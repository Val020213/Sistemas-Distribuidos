package chord

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sort"
	"time"

	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "server/internal/chord/chordpb"
	"server/internal/models"
	"server/internal/scraper"
	"server/internal/utils"
)

type RingNode struct {
	Id      uint64 // Node's Id (computed from its address)
	Address string // Host
	Port    string // Port

	Scraper *scraper.Scraper // = Scraper Application

	Successors     []*pb.Node // Immediate successor in the ring
	Predecessor    *pb.Node   // Immediate predecessor in the ring
	Finger         []*pb.Node // Finger table entries
	SuccessorCache []*pb.Node // Cache of verified successors

	// Data    map[uint64]models.TaskType // Simple key-value storage
	M       int    // Number of bits in the hash space
	idSpace uint64 // Number of nodes in the hash space

	mu sync.Mutex // Protects access to mutable fields

	pb.UnimplementedChordServiceServer
}

var (
	grpcAddr      = os.Getenv("IP_ADDRESS")
	grpcPort      = os.Getenv("RPC_PORT")
	tolerance     = 3 // update this to the environment variables
	mBits         = utils.GetEnvAsInt("CHORD_BITS", 8)
	multicastAddr = "224.0.0.1:9999"
)

func NewNode() *RingNode {

	Id := utils.ChordHash(grpcAddr, mBits)
	scraper := scraper.NewScraper()
	return &RingNode{
		Id:      Id,
		Address: grpcAddr,
		Port:    grpcPort,
		Finger:  make([]*pb.Node, mBits),
		M:       mBits,
		idSpace: 1 << mBits,
		Scraper: scraper,

		Successors:     []*pb.Node{},
		Predecessor:    nil,
		SuccessorCache: []*pb.Node{},
		mu:             sync.Mutex{},
	}
}

func (n *RingNode) StartRPCServer(grpcServer *grpc.Server) {
	pb.RegisterChordServiceServer(grpcServer, n)
	fmt.Println("Starting gRPC Server on ", n.Address, ":", n.Port)
	n.initNode()
	n.joinNetwork()

}

func (n *RingNode) initNode() {
	for i := 0; i < n.M; i++ {
		n.Finger[i] = n.MakeNode()
	}
}

// http server

func (n *RingNode) CallCreateData(data models.TaskType) error {

	node, err := n.FindSuccessor(context.Background(), &pb.FindSuccessorRequest{Key: data.Key, Hops: 0, Visited: nil})

	if err != nil {
		fmt.Println("ERROR STORING ", data.URL)
		return err
	}

	client, conn, err := n.GetClient(node.Address)

	if err != nil {
		return n.CallCreateData(data)
	}
	defer conn.Close()

	pbData := *ToPbData(&data)

	_, err = client.CreateData(context.Background(), &pb.CreateDataRequest{Data: &pbData})

	return err
}

func (n *RingNode) CallGetData(url string) (string, error) {

	key := uint64(utils.ChordHash(url, n.M))
	node, err := n.FindSuccessor(context.Background(), &pb.FindSuccessorRequest{Key: key, Hops: 0, Visited: nil})

	if err != nil {
		fmt.Println("ERROR RETRIEVING ", url)
		return "", err
	}

	client, conn, err := n.GetClient(node.Address)

	if err != nil {
		return n.CallGetData(url)
	}
	defer conn.Close()

	data, err := client.RetrieveData(context.Background(), &pb.Id{Id: key})

	if err != nil {
		return "", err
	}

	return data.Content, nil
}

func (n *RingNode) CallGetStatus() ([]models.TaskType, error) {

	node := n.GetFirstAliveSuccessor()

	client, conn, err := n.GetClient(node.Address)

	if err != nil {
		return n.CallGetStatus()
	}
	defer conn.Close()

	client.PrintState(context.Background(), &pb.Empty{})

	return nil, nil
}

func (n *RingNode) CallList() ([]models.TaskType, error) {

	predecessor := n.CheckPredecessor()
	successors := n.Successors
	visited := make(map[uint64]bool)

	tasks, err := n.Scraper.DB.GetTasksWithFilter(utils.GetFilterBetweenRightInclusive(predecessor.Id, n.Id))

	if err != nil {
		utils.RedPrint("ERROR CALLING LIST IN FIRST NODE ", n.Address)
		return nil, err
	}

	for ind := range tasks {
		tasks[ind].Content = ""
	}

	for _, node := range successors {

		if _, ok := visited[node.Id]; node.Address != "" && ok && node.Address != n.Address {

			client, conn, err := n.GetClient(node.Address)

			if err != nil {
				utils.RedPrint("ERROR CALLING LIST IN NODE :", node.Address)
				continue
			}
			defer conn.Close()

			response, err := client.List(context.Background(), &pb.Empty{})

			if err != nil {
				utils.RedPrint("ERROR WHILE LISTING ON NODE ", node.Address, ": ", err)
				return nil, err
			}

			for _, task := range response.Data {
				tasks = append(tasks, *FromPbData(task))
			}

			visited[node.Id] = true
			successors = append(successors, response.Successors...)
		}

	}

	return tasks, nil

}

// gRPC Chord Protocol

func (n *RingNode) Notify(ctx context.Context, node *pb.Node) (*pb.Successful, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Println("Notified by ", node.Address)

	if n.Predecessor == nil || utils.Between(node.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = &pb.Node{Id: node.Id, Address: node.Address}
	}

	return &pb.Successful{Successful: true}, nil
}

func (n *RingNode) PrintState(ctx context.Context, empty *pb.Empty) (*pb.State, error) {

	stateData := []*pb.Data{}

	tasks, err := n.Scraper.DB.GetTasks()

	if err != nil {
		utils.RedPrint("Error communicating with database in Print State")
		return nil, err
	}

	for _, cData := range tasks {
		stateData = append(stateData, ToPbData(&cData))
	}

	return &pb.State{
		Id:          n.Id,
		Addr:        n.Address,
		Data:        stateData,
		Finger:      n.Finger,
		Successors:  n.Successors,
		Predecessor: n.Predecessor,
	}, nil
}

func (n *RingNode) Health(ctx context.Context, empty *pb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Id:      n.Id,
		Address: n.Address,
	}, nil
}

func (n *RingNode) GetSuccessors(ctx context.Context, empty *pb.Empty) (*pb.GetSuccessorsResponse, error) {
	return &pb.GetSuccessorsResponse{Successors: n.SuccessorCache}, nil
}

func (n *RingNode) GetPredecessor(ctx context.Context, empty *pb.Empty) (*pb.Node, error) {
	return n.Predecessor, nil
}

func (n *RingNode) FindSuccessor(ctx context.Context, request *pb.FindSuccessorRequest) (*pb.Node, error) {

	key, hops, visited := request.Key, request.Hops, request.Visited

	if visited == nil {
		visited = make(map[uint64]bool)
	}

	if int(hops) > n.M {
		return nil, errors.New("too many hops")
	}

	if key == n.Id {
		return n.MakeNode(), nil
	}

	successor := n.GetFirstAliveSuccessor()
	if utils.BetweenRightInclusive(key, n.Id, successor.Id) {
		return successor, nil
	}

	closest, _ := n.ClosestPrecedingFinger(key)
	if closest.Id == n.Id || visited[closest.Id] {
		return successor, nil
	}

	closestClient, conn, err := n.GetClient(closest.Address)

	if err != nil {
		return n.FindSuccessor(ctx, request)
	}
	defer conn.Close()

	visited[n.Id] = true
	return closestClient.FindSuccessor(ctx, &pb.FindSuccessorRequest{Key: key, Hops: hops + 1, Visited: visited})
}

func (n *RingNode) CreateData(ctx context.Context, data *pb.CreateDataRequest) (*pb.Successful, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.createData(data.Data)
	return &pb.Successful{Successful: true}, nil
}

func (n *RingNode) StoreData(ctx context.Context, data *pb.StoreDataRequest) (*pb.Successful, error) {

	n.updateData(data.Data)
	return &pb.Successful{Successful: true}, nil
}

func (n *RingNode) RetrieveData(ctx context.Context, key *pb.Id) (*pb.Data, error) {

	task, err := n.Scraper.DB.GetTask(key.Id)

	if err != nil {
		utils.RedPrint("Error retrieving data: ", err)
		return nil, err
	}

	return ToPbData(&task), nil
}

func (n *RingNode) DeleteData(ctx context.Context, deleteId *pb.Id) (*pb.Successful, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	err := n.Scraper.DB.DeleteData(utils.GetFilterBetweenRightInclusive(n.Id, deleteId.Id))

	if err != nil {
		utils.RedPrint("ERROR DELETING ", deleteId.Id)
		return nil, err
	}

	return &pb.Successful{Successful: true}, nil
}

func (n *RingNode) GetRangeNodeData(ctx context.Context, request *pb.GetNodeDataRequest) (*pb.StoreDataRequest, error) {
	data, err := n.Scraper.DB.GetTasksWithFilter(utils.GetFilterBetweenRightInclusive(request.PredecesorId, request.Id))
	if err != nil {
		utils.RedPrint("GetRangeNodeData Error", err)
		return &pb.StoreDataRequest{}, err
	}
	var pbData []*pb.Data
	for _, task := range data {
		pbData = append(pbData, ToPbData(&task))
	}
	return &pb.StoreDataRequest{Data: pbData}, nil
}

func (n *RingNode) List(ctx context.Context, empty *pb.Empty) (*pb.ListResponse, error) {

	predecessor := n.CheckPredecessor()

	data, err := n.Scraper.DB.GetTasksWithFilter(utils.GetFilterBetweenRightInclusive(predecessor.Id, n.Id))

	if err != nil {
		utils.RedPrint("DATABASE ERROR IN LIST ", err)
		return nil, err
	}

	var pbData []*pb.Data

	for _, task := range data {
		task.Content = ""
		pbData = append(pbData, ToPbData(&task))
	}

	return &pb.ListResponse{Successors: n.Successors, Data: pbData}, nil
}

// Chord Utils

func (n *RingNode) joinNetwork() (string, error) {

	fmt.Println("Joining network...")

	bootstrapNode, err := n.GetBootstrapNode()

	if err != nil {
		fmt.Println("JoinNetwork: Error getting bootstrap node: ", err)
		return n.joinNetwork()
	}

	if bootstrapNode == nil || bootstrapNode.Address == n.Address {
		fmt.Println("JoinNetwork: I am the bootstrap node")
		n.updateSuccessors([]*pb.Node{n.MakeNode()})
	} else {
		fmt.Println("Estableciendo conexión con ", bootstrapNode.Address)
		clientNode, conn, err := n.GetClient(bootstrapNode.Address)

		if err != nil {
			fmt.Println("JoinNetwork: Error connecting to node: ", err)
			return n.joinNetwork()
		}

		defer conn.Close()
		fmt.Println("Conexión establecida: ", n.Id, " -> ", bootstrapNode.Id)
		fmt.Println("Buscando sucesor...")
		succ, err := clientNode.FindSuccessor(context.Background(), &pb.FindSuccessorRequest{Key: n.Id, Hops: 0, Visited: nil})
		fmt.Println("Sucesor encontrado: ", succ.Id)
		if err != nil {
			fmt.Println("JoinNetwork: Error en la conexión: ", err)
			return n.joinNetwork()
		}

		predecessors := []*pb.Node{}

		succClient, conn, err := n.GetClient(succ.Address)
		fmt.Println("Conexión establecida con sucesor: ", succ.Address)
		if err != nil {
			fmt.Println("JoinNetwork: Error connecting to successor: ", err)
			return n.joinNetwork()
		}

		defer conn.Close()

		fmt.Println("Buscando predecesor...")
		succPredecessor, err := succClient.GetPredecessor(context.Background(), &pb.Empty{})
		fmt.Println("Predecesor encontrado: ", succPredecessor.Id)

		if err != nil {
			fmt.Println("JoinNetwork: Error getting predecessor: ", err)
			return n.joinNetwork()
		}

		candidatesSuccessors := []*pb.Node{}

		if succPredecessor.Address != "" {
			fmt.Println("Estableciendo conexión con predecesor: ", succPredecessor.Address)
			succPredecessorClient, conn, err := n.GetClient(succPredecessor.Address)
			fmt.Println("Conexión establecida con predecesor: ", succPredecessor.Address)

			if err != nil {
				fmt.Println("JoinNetwork: Error connecting to predecessor: ", err)
				return n.joinNetwork()
			}

			defer conn.Close()
			fmt.Println("Buscando sucesores del predecesor...")
			candidatesSucc, err := succPredecessorClient.GetSuccessors(context.Background(), &pb.Empty{})

			if err != nil {
				fmt.Println("JoinNetwork: Error getting successors: ", err)
				return n.joinNetwork()
			}

			fmt.Println("Sucesores del predecesor encontrados: ", candidatesSucc.Successors)
			candidatesSuccessors = append(candidatesSuccessors, candidatesSucc.Successors...)
		}

		fmt.Println("Actualizando predecesor...")
		predecessors = append(predecessors, candidatesSuccessors...)
		fmt.Println("Predecesores: ", predecessors)
		n.updateSuccessors(append([]*pb.Node{succ}, predecessors...))
	}
	n.RunPeriodicTasks()
	return n.Address, nil
}

func (n *RingNode) RunPeriodicTasks() {
	go func() {
		for {
			n.CheckPredecessor()
			n.Stabilize()
			time.Sleep(1 * time.Second)
		}
	}()
}

func (n *RingNode) CheckPredecessor() *pb.Node {
	fmt.Println("CHECK PREDECESSOR ")
	if n.Predecessor != nil && !n.IsAlive(n.Predecessor) {
		n.Predecessor = nil
		return n.MakeNode()
	}

	if n.Predecessor == nil {
		return n.MakeNode()
	}

	return n.Predecessor
}

func (n *RingNode) ClosestPrecedingFinger(key uint64) (*pb.Node, error) {

	for i := n.M - 1; i >= 0; i-- {
		if utils.BetweenRightInclusive(n.Finger[i].Id, n.Id, key) && n.IsAlive(n.Finger[i]) { // contemplando quitar el isAlive
			return n.Finger[i], nil
		}
	}

	return n.MakeNode(), nil
}

func (n *RingNode) FindNewSuccessor(candidate, oldSuccessor *pb.Node) (*pb.Node, error) {

	if candidate.Address != "" && n.Id != candidate.Id && utils.Between(candidate.Id, n.Id, oldSuccessor.Id) {
		succClient, conn, err := n.GetClient(candidate.Address)

		if err != nil {
			utils.RedPrint("FindNewSuccessor: Error connecting to candidate: ", err)
			return oldSuccessor, nil
		}

		defer conn.Close()
		newCandidate, err := succClient.GetPredecessor(context.Background(), &pb.Empty{})

		if err != nil {
			utils.RedPrint("FindNewSuccessor: Error getting predecessor: ", err)
			return oldSuccessor, nil
		}

		return n.FindNewSuccessor(newCandidate, candidate)

	}
	return oldSuccessor, nil
}

func (n *RingNode) GetNodeDataFromOldSuccessor() {
	utils.YellowPrint("GET NODE DATA FROM OLD SUCCESSORS")
	if n.Predecessor != nil {
		if len(n.Successors) > 0 {
			for _, succ := range n.Successors[:len(n.Successors)-1] {
				utils.YellowPrint("FROM SUCCESSOR: ", succ.Address)
				succClient, conn, err := n.GetClient(succ.Address)

				if err != nil {
					utils.RedPrint("Error Getting Node data from old successor ", n.Id, " error:", err)
					continue
				}

				defer conn.Close()

				data, err := succClient.GetNodeData(context.Background(), &pb.GetNodeDataRequest{Id: n.Id, PredecesorId: n.Predecessor.Id})

				if err != nil {
					utils.RedPrint("Error Getting Node data from old successor ", n.Id, " error:", err)
					continue
				}

				n.updateData(data.Data)
			}
		}
	}
}

func (n *RingNode) Stabilize() {
	n.GetNodeDataFromOldSuccessor()

	// UPDATE SUCCESSORS
	succ := n.GetFirstAliveSuccessor()

	succClient, conn, err := n.GetClient(succ.Address)

	if err != nil {
		fmt.Println("STABILIZE: Error connecting to successor: ", err)
		return
	}

	defer conn.Close()

	newSuccessor, err := succClient.GetPredecessor(context.Background(), &pb.Empty{})

	if err != nil {
		fmt.Println("STABILIZE: Error getting predecessor: ", err)
		return
	}

	newestSuccessor, err := n.FindNewSuccessor(newSuccessor, succ)

	if err != nil {
		utils.RedPrint("STABILIZE: Error Finding new Successor ", err)
	}

	newestSuccessorClient, conn, err := n.GetClient(newestSuccessor.Address)

	if err != nil {
		utils.RedPrint("STABILIZE: Error connecting to predecessor: ", err)
		return
	}

	defer conn.Close()

	newestSuccessorResponse, err := newestSuccessorClient.GetSuccessors(context.Background(), &pb.Empty{})

	if err != nil {
		utils.RedPrint("STABILIZE: Error getting successors: Successor Address: ", newestSuccessor.Address, " Error:", err)
		return
	}

	n.updateSuccessors(append([]*pb.Node{newestSuccessor}, newestSuccessorResponse.Successors...))

	// UPDATE SUCCESSORS ENDS

	n.FixFingersTable()

	if succ.Id != n.Id {
		n.replicateData()
	}

}

func (n *RingNode) FixFingersTable() {
	fmt.Println("FixFingerTable...")
	for i := 0; i < n.M; i++ {
		fingerKey := (n.Id + (1 << i)) % n.idSpace
		node, err := n.FindSuccessor(context.Background(), &pb.FindSuccessorRequest{Key: fingerKey, Hops: 0, Visited: nil})
		if node != nil && err == nil {
			n.Finger[i] = node
		} else {
			n.Finger[i] = n.GetFirstAliveSuccessor()
		}
	}
	fmt.Println("FingerTable: ", n.Finger)
}

func (n *RingNode) updateSuccessors(newSuccessors []*pb.Node) {

	merged := []*pb.Node{}
	seen := make(map[uint64]bool)

	for _, node := range append(n.Successors, newSuccessors...) {
		if node.Address != "" && node.Id != n.Id && !seen[node.Id] && n.IsAlive(node) {
			merged = append(merged, node)
			seen[node.Id] = true
			if len(merged) > tolerance {
				break
			}
		}
	}

	sort.Slice(merged, func(i, j int) bool {
		diffI := (merged[i].Id - n.Id) % n.idSpace
		diffJ := (merged[j].Id - n.Id) % n.idSpace
		return diffI < diffJ
	})

	newSuccessors = make([]*pb.Node, 0, tolerance)
	for i := 0; i < len(merged) && i < tolerance; i++ {
		newSuccessors = append(newSuccessors, merged[i])
	}

	if len(newSuccessors) > 0 && (len(n.SuccessorCache) == 0 || newSuccessors[0].Id != n.SuccessorCache[0].Id) {
		newSuccessorClient, conn, err := n.GetClient(newSuccessors[0].Address)

		if err != nil {
			utils.RedPrint("UPDATE SUCCESSORS ERROR: Notify connection error to", newSuccessors[0].Address, err)
			return
		}

		defer conn.Close()

		newSuccessorClient.Notify(context.Background(), n.MakeNode())
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.Successors = newSuccessors
	n.SuccessorCache = newSuccessors

	fmt.Println("UpdateSuccessors: ", n.Successors)
}

func (n *RingNode) GetFirstAliveSuccessor() *pb.Node {
	for _, node := range n.Successors {
		if n.IsAlive(node) {
			return node
		}
	}
	return n.MakeNode()
}

func (n *RingNode) IsAlive(remoteNode *pb.Node) bool {
	client, conn, err := n.GetClient(remoteNode.Address)

	if err != nil {
		utils.RedPrint("IsAlive: Error connecting to node ", remoteNode.Address)
		return false
	}

	defer conn.Close()

	resp, err := client.Health(context.Background(), &pb.Empty{})

	if err != nil || resp == nil {
		utils.RedPrint("IsAlive: Node is not alive ", remoteNode.Address)
		return false
	}

	return true
}

func (n *RingNode) replicateData() {
	fmt.Println("Replicating data...")
	predecessor := n.CheckPredecessor()
	fmt.Println("REPLICATE DATA: predecessor ", predecessor.Id)
	predecessorId := predecessor.Id

	if predecessorId == n.Id {
		return
	}

	tasks, err := n.Scraper.DB.GetTasksWithFilter(utils.GetFilterBetweenRightInclusive(predecessorId, n.Id))

	if err != nil {
		utils.RedPrint("DATABASE ERROR IN REPLICATE DATA")
		return
	}

	replicated := []*pb.Data{}

	for _, task := range tasks {
		replicated = append(replicated, ToPbData(&task))
	}

	for _, successor := range n.Successors {

		if n.Id != successor.Id {

			successorClient, conn, err := n.GetClient(successor.Address)

			if err != nil { // handle error
				utils.RedPrint("Error connecting to successor: ", err)
				return
			}

			defer conn.Close()
			successorClient.StoreData(context.Background(), &pb.StoreDataRequest{Data: replicated})
		}

	}

	if len(n.Successors) >= tolerance {
		lastSuccessorClient, conn, err := n.GetClient(n.Successors[tolerance-1].Address)
		if err != nil {
			utils.RedPrint("Error connecting to last successor: ", err)
			return
		}
		defer conn.Close()
		lastSuccessorClient.DeleteData(context.Background(), &pb.Id{Id: predecessorId})
	}

}

func (n *RingNode) createData(d *pb.Data) error {
	task := models.TaskType{
		URL:       d.Url,
		Key:       d.Key,
		Status:    models.TaskStatusType(d.Status),
		Content:   d.Content,
		CreatedAt: d.CreatedAt.AsTime(),
		UpdatedAt: d.UpdatedAt.AsTime(),
	}

	_, err := n.Scraper.DB.CreateTask(task)

	if err != nil {
		utils.RedPrint("ERROR ADDING ", d.Url, " TO DATABASE")
		return err
	}

	return nil

}

func (n *RingNode) updateData(data []*pb.Data) {
	n.mu.Lock()
	defer n.mu.Unlock()

	modelData := []models.TaskType{}
	for _, d := range data {
		task := models.TaskType{
			URL:       d.Url,
			Key:       d.Key,
			Status:    models.TaskStatusType(d.Status),
			Content:   d.Content,
			CreatedAt: d.CreatedAt.AsTime(),
			UpdatedAt: d.UpdatedAt.AsTime(),
		}

		modelData = append(modelData, task)
	}

	err := n.Scraper.DB.UpdateTasks(modelData)
	if err != nil {
		utils.RedPrint("ERROR UPDATING TASKS IN DATABASE: ", err)
	}
}

// gRPC Client

func (n *RingNode) GetClient(addr string) (pb.ChordServiceClient, *grpc.ClientConn, error) {
	if addr == "" {
		return nil, nil, errors.New("empty address")
	}
	conn, err := grpc.NewClient(utils.ChangePort(addr, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.RedPrint("GetClient: Error connecting to node: ", err)
		return nil, nil, err
	}
	return pb.NewChordServiceClient(conn), conn, nil
}

// Discover multicast

func (n *RingNode) GetBootstrapNode() (*pb.Node, error) {
	fmt.Println("Getting bootstrap node...")
	addr, err := n.Discover()

	if err != nil {
		utils.RedPrint("Error discovering network: ", err)
		return nil, err
	}

	if addr == n.Address {
		fmt.Println("I am the bootstrap node")
		return nil, nil
	}

	return &pb.Node{Id: utils.ChordHash(addr, n.M), Address: addr}, nil
}

func (n *RingNode) Discover() (string, error) {
	fmt.Println("Discovering network...")

	addr, err := net.ResolveUDPAddr("udp4", multicastAddr)
	if err != nil {
		utils.RedPrint("Error resolving multicast address: %v", err)
		return "", err
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		utils.RedPrint("Error listening on multicast: %v", err)
		return "", err
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	deadline := time.Now().Add(6 * time.Second)
	conn.SetReadDeadline(deadline)

	for {
		_, src, err := conn.ReadFromUDP(buf)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Timeout, no se encontraron nodos")
				return n.Address, nil
			}
			return "", fmt.Errorf("error leyendo desde multicast: %v", err)
		}

		if src.String() != n.Address {
			fmt.Printf("Encontrado nodo en %s\n", src.String())
			return src.String(), nil
		}
	}
}

// utils

func (n *RingNode) MakeNode() *pb.Node {
	return &pb.Node{Id: n.Id, Address: n.Address}
}

func ToPbData(data *models.TaskType) *pb.Data {
	return &pb.Data{
		Url:       data.URL,
		Key:       data.Key,
		Status:    string(data.Status),
		Content:   data.Content,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}
}

func FromPbData(data *pb.Data) *models.TaskType {
	return &models.TaskType{
		URL:       data.Url,
		Key:       data.Key,
		Status:    models.TaskStatusType(data.Status),
		Content:   data.Content,
		CreatedAt: data.CreatedAt.AsTime(),
		UpdatedAt: data.UpdatedAt.AsTime(),
	}
}
