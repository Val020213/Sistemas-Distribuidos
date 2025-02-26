package chord

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"time"

	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "server/internal/chord/chordpb"
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

	Data    map[string]string // Simple key-value storage
	m       int               // Number of bits in the hash space
	idSpace uint64            // Number of nodes in the hash space

	mu sync.Mutex // Protects access to mutable fields

	pb.UnimplementedChordServiceServer
}

var (
	grpcAddr      = os.Getenv("IP_ADDRESS")
	grpcPort      = os.Getenv("RPC_PORT")
	tolerance     = 3 // update this to the environment variables
	mBits         = utils.GetEnvAsInt("CHORD_BITS", 3)
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
		Data:    make(map[string]string),
		m:       mBits,
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
	n.initNode(n.m)
	n.joinNetwork()

}

func (n *RingNode) initNode(mBits int) {
	for i := 0; i < mBits; i++ {
		n.Finger[i] = n.MakeNode()
	}
}

// gRPC Chord Protocol

func (n *RingNode) Notify(ctx context.Context, node *pb.Node) (*pb.StoreDataRequest, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Println("Notified by ", node.Address)

	data := []*pb.Data{}

	if n.Predecessor == nil || utils.Between(node.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = &pb.Node{Id: node.Id, Address: node.Address}
		fmt.Println("Updated predecessor to ", node.Address)

		for key := range n.Data {
			if !utils.BetweenRightInclusive(utils.ChordHash(key, n.m), n.Predecessor.Id, n.Id) {
				data = append(data, &pb.Data{Key: key, Value: n.Data[key]})
			}
		}
	}

	return &pb.StoreDataRequest{Data: data}, nil
}

func (n *RingNode) Health(ctx context.Context, empty *pb.Empty) (*pb.HealthResponse, error) {
	fmt.Println(" * * * * * *")
	fmt.Println("Health check")
	fmt.Println("Node ", n.Address, " is alive")
	fmt.Println("Successors: ", n.Successors)
	fmt.Println("Predecessor: ", n.Predecessor)
	fmt.Println("Finger: ", n.Finger)
	fmt.Println("Data: ", n.Data)
	fmt.Println(" * * * * * *")
	return &pb.HealthResponse{
		Id:      n.Id,
		Address: n.Address,
	}, nil
}

func (n *RingNode) GetSuccessors(ctx context.Context, empty *pb.Empty) (*pb.GetSuccessorsResponse, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	return &pb.GetSuccessorsResponse{Successors: n.Successors}, nil
}

func (n *RingNode) GetPredecessor(ctx context.Context, empty *pb.Empty) (*pb.Node, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.Predecessor, nil
}

func (n *RingNode) FindSuccessor(ctx context.Context, request *pb.FindSuccessorRequest) (*pb.Node, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	key, hops, visited := request.Key, request.Hops, request.Visited

	if visited == nil {
		visited = make(map[uint64]bool)
	}

	if int(hops) > n.m {
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

func (n *RingNode) StoreData(ctx context.Context, data *pb.StoreDataRequest) (*pb.Successful, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.updateData(data.Data)
	return &pb.Successful{Successful: true}, nil
}

func (n *RingNode) DeleteData(ctx context.Context, deleteId *pb.Id) (*pb.Successful, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	deleteData := []*pb.Data{}
	for dataKey := range n.Data {
		if utils.BetweenRightInclusive(utils.ChordHash(dataKey, n.m), n.Id, deleteId.Id) {
			deleteData = append(deleteData, &pb.Data{Key: dataKey, Value: n.Data[dataKey]})
		}
	}

	fmt.Println("Deleting data: ", deleteData)
	for _, data := range deleteData {
		delete(n.Data, data.Key)
	}
	return &pb.Successful{Successful: true}, nil
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
		clientNode, conn, err := n.GetClient(bootstrapNode.Address)

		if err != nil {
			fmt.Println("JoinNetwork: Error connecting to node: ", err)
			return n.joinNetwork()
		}

		defer conn.Close()

		succ, err := clientNode.FindSuccessor(context.Background(), &pb.FindSuccessorRequest{Key: n.Id, Hops: 0, Visited: nil})

		if err != nil {
			fmt.Println("JoinNetwork: Error en la conexión: ", err)
			return n.joinNetwork()
		}

		predecessors := []*pb.Node{}

		succClient, conn, err := n.GetClient(succ.Address)

		if err != nil {
			fmt.Println("JoinNetwork: Error connecting to successor: ", err)
			return n.joinNetwork()
		}

		defer conn.Close()

		succPredecessor, err := succClient.GetPredecessor(context.Background(), &pb.Empty{})

		if err != nil {
			fmt.Println("JoinNetwork: Error getting predecessor: ", err)
			return n.joinNetwork()
		}

		candidatesSuccessors := []*pb.Node{}

		if succPredecessor.Address != "" {

			succPredecessorClient, conn, err := n.GetClient(succPredecessor.Address)

			if err != nil {
				fmt.Println("JoinNetwork: Error connecting to predecessor: ", err)
				return n.joinNetwork()
			}

			defer conn.Close()

			candidatesSucc, err := succPredecessorClient.GetSuccessors(context.Background(), &pb.Empty{})

			if err != nil {
				fmt.Println("JoinNetwork: Error getting successors: ", err)
				return n.joinNetwork()
			}

			candidatesSuccessors = append(candidatesSuccessors, candidatesSucc.Successors...)
		}

		predecessors = append(predecessors, candidatesSuccessors...)

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

	for i := n.m - 1; i >= 0; i-- {
		if utils.BetweenRightInclusive(n.Finger[i].Id, n.Id, key) && n.IsAlive(n.Finger[i]) { // contemplando quitar el isAlive
			return n.Finger[i], nil
		}
	}

	return n.MakeNode(), nil
}

func (n *RingNode) Stabilize() {

	succ := n.GetFirstAliveSuccessor()

	if succ.Id != n.Id {

		succClient, conn, err := n.GetClient(succ.Address)

		if err != nil {
			fmt.Println("STABILIZE: Error connecting to successor: ", err)
			return
		}

		defer conn.Close()

		succPredecessor, err := succClient.GetPredecessor(context.Background(), &pb.Empty{})

		if err != nil {
			fmt.Println("STABILIZE: Error getting predecessor: ", err)
			return
		}

		if succPredecessor.Address != "" {

			succPredecessorClient, conn, err := n.GetClient(succPredecessor.Address)

			if err != nil {
				fmt.Println("STABILIZE: Error connecting to predecessor: ", err)
				return
			}

			defer conn.Close()

			if utils.Between(succPredecessor.Id, n.Id, succ.Id) {
				succPredecessorSuccessors, err := succPredecessorClient.GetSuccessors(context.Background(), &pb.Empty{})

				if err != nil {
					fmt.Println("STABILIZE: Error getting successors: Successor Address: ", succPredecessor.Address, " Error:", err)
					return
				}

				n.updateSuccessors(append([]*pb.Node{succPredecessor}, succPredecessorSuccessors.Successors...))
			}
		}

		fmt.Println("STABILIZE: Notify to successor")
		retrievedData, err := succClient.Notify(context.Background(), n.MakeNode())

		if err != nil {
			fmt.Println("STABILIZE: Error notifying successor: ", err)
			return
		}

		newSuccessors, err := succClient.GetSuccessors(context.Background(), &pb.Empty{})

		if err != nil {
			fmt.Println("STABILIZE: Error getting successors: ", err)
			return
		}

		n.updateSuccessors(append([]*pb.Node{succ}, newSuccessors.Successors...))

		n.updateData(retrievedData.Data)
	}

	n.FixFingersTable()

	if succ.Id != n.Id {
		n.replicateData()
	}

}

func (n *RingNode) FixFingersTable() {
	fmt.Println("FixFingerTable...")
	for i := 0; i < n.m; i++ {
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
	n.mu.Lock()
	defer n.mu.Unlock()

	merged := []*pb.Node{}
	seen := make(map[uint64]bool)

	for _, node := range append(n.Successors, newSuccessors...) {
		if node.Id != n.Id {
			if _, ok := seen[node.Id]; !ok { // Handle alive nodes in a gorutine or channel to avoid blocking
				merged = append(merged, node)
				seen[node.Id] = true

				if len(merged) > tolerance {
					break
				}
			}
		}
	}

	sort.Slice(merged, func(i, j int) bool { // Warning with None
		diffI := (merged[i].Id - n.Id) % n.idSpace
		diffJ := (merged[j].Id - n.Id) % n.idSpace
		return diffI < diffJ
	})

	newSuccessors = make([]*pb.Node, 0, tolerance+1)
	for i := 0; i < len(merged) && i < tolerance+1; i++ {
		newSuccessors = append(newSuccessors, merged[i])
	}
	n.Successors = newSuccessors
	n.SuccessorCache = newSuccessors // here we can add a goroutine to update the cache

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
		return false
	}

	defer conn.Close()

	resp, err := client.Health(context.Background(), &pb.Empty{})

	if err != nil || resp == nil {
		return false
	}

	return true
}

// data management

func (n *RingNode) replicateData() {
	fmt.Println("Replicating data...")
	predecessor := n.CheckPredecessor()
	fmt.Println("REPLICATE DATA: predecessor ", predecessor.Id)
	predecessorId := predecessor.Id

	for _, successor := range n.Successors {
		replicated := []*pb.Data{}
		for key, value := range n.Data {
			if n.Id != successor.Id && utils.Between(utils.ChordHash(key, n.m), predecessorId, successor.Id) {
				replicated = append(replicated, &pb.Data{Key: key, Value: value})
			}
		}
		successorClient, conn, err := n.GetClient(successor.Address)

		if err != nil { // handle error
			fmt.Println("Error connecting to successor: ", err)
			n.Stabilize()
			continue
		}

		defer conn.Close()
		successorClient.StoreData(context.Background(), &pb.StoreDataRequest{Data: replicated})
	}

	if len(n.Successors) >= tolerance {
		lastSuccessorClient, conn, err := n.GetClient(n.Successors[tolerance-1].Address)
		if err != nil {
			fmt.Println("Error connecting to last successor: ", err)
			return
		}
		defer conn.Close()
		lastSuccessorClient.DeleteData(context.Background(), &pb.Id{Id: predecessorId})
	}

}

func (n *RingNode) updateData(data []*pb.Data) {
	for _, d := range data {
		n.Data[d.Key] = d.Value
	}
}

// gRPC Client
func (n *RingNode) GetClient(addr string) (pb.ChordServiceClient, *grpc.ClientConn, error) {
	if addr == "" {
		fmt.Printf("\033[31mGetClient: Addrs is EMPTY!!!\033[0m\n")
	}
	conn, err := grpc.NewClient(utils.ChangePort(addr, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("GetClient: Error connecting to node: ", err)
		return nil, nil, err
	}
	return pb.NewChordServiceClient(conn), conn, nil
}

// Discover multicast

func (n *RingNode) GetBootstrapNode() (*pb.Node, error) {
	fmt.Println("Getting bootstrap node...")
	addr, err := n.Discover()

	if err != nil {
		fmt.Println("Error discovering network: ", err)
		return nil, err
	}

	if addr == n.Address {
		fmt.Println("I am the bootstrap node")
		return nil, nil
	}

	return &pb.Node{Id: utils.ChordHash(addr, n.m), Address: addr}, nil
}

func (n *RingNode) Discover() (string, error) {
	fmt.Println("Discovering network...")

	addr, err := net.ResolveUDPAddr("udp4", multicastAddr)
	if err != nil {
		log.Fatalf("Error resolving multicast address: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("Error listening on multicast: %v", err)
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
