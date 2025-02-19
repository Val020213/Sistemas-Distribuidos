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
		Data:    make(map[string]string),
		m:       mBits,
		idSpace: 1 << mBits,
		Scraper: scraper,

		Successors:     make([]*pb.Node, 1),
		Predecessor:    nil,
		SuccessorCache: make([]*pb.Node, 1),
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

func (n *RingNode) Notify(ctx context.Context, node *pb.Node) (*pb.Successful, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	fmt.Println("Notified by ", node.Address)

	if n.Predecessor == nil || utils.Between(node.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = &pb.Node{Id: node.Id, Address: node.Address}
		fmt.Println("Updated predecessor to ", node.Address)
		return &pb.Successful{Successful: true}, nil
	}

	return &pb.Successful{Successful: false}, nil
}

func (n *RingNode) Health(ctx context.Context, empty *pb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Id:      n.Id,
		Address: n.Address,
	}, nil
}

func (n *RingNode) joinNetwork() (string, error) {

	fmt.Println("Joining network...")

	bootstrapNode, err := n.GetBootstrapNode()

	if err != nil {
		fmt.Println("Error getting bootstrap node: ", err)
		return n.joinNetwork()
	}

	if bootstrapNode == nil {
		n.updateSuccessors([]*pb.Node{n.MakeNode()})
	} else {
		clientNode, conn, err := n.GetClient(bootstrapNode.Address)

		if err != nil {
			fmt.Println("Error connecting to node: ", err)
			return n.joinNetwork()
		}
		defer conn.Close()

		succ, err := clientNode.FindSuccessor(context.Background(), &pb.FindSuccessorRequest{Key: n.Id, Hops: 0, Visited: nil})

		if err != nil {
			fmt.Println("Error en la conexión: ", err)
			return "", errors.New("no se encontraron nodos para enlazar")
		}

		predecessors := []*pb.Node{}
		if succ.Predecessor != nil {
			predecessors = append(predecessors, succ.Predecessor.Successors...) // use get predecessor service if available
		}

		n.updateSuccessors(append([]*pb.Node{succ}, predecessors...))
	}

	// start background tasks
	n.RunPeriodicTasks()
	return n.Address, nil
}

func (n *RingNode) RunPeriodicTasks() {
	go func() {
		for {
			n.CheckPredecessor()
			n.Stabilize()
			n.FixFingersTable()
			// TODO : REPLICATE DATA
			time.Sleep(1 * time.Second)
		}
	}()
}

func (n *RingNode) CheckPredecessor() {
	if n.Predecessor != nil && !n.IsAlive(n.Predecessor) {
		n.Predecessor = nil
	}

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
	if succ != nil {
		x := succ.Predecessor
		if x != nil && n.IsAlive(x) && utils.Between(x.Id, n.Id, succ.Id) {
			n.updateSuccessors(append([]*pb.Node{x}, x.Successors...))
		}
		succClient, conn, err := n.GetClient(succ.Address)

		if err != nil {
			fmt.Println("Error connecting to successor: ", err)
			n.Stabilize()
			return
		}
		defer conn.Close()
		succClient.Notify(context.Background(), n.MakeNode())
		// TODO: TRANSFER DATA
	}
}

func (n *RingNode) FixFingersTable() {
	for i := 0; i < n.m; i++ {
		fingerKey := (n.Id + (1 << i)) % n.idSpace
		node, err := n.FindSuccessor(context.Background(), &pb.FindSuccessorRequest{Key: fingerKey, Hops: 0, Visited: nil})
		if node != nil && err == nil {
			n.Finger[i] = node
		} else {
			n.Finger[i] = n.GetFirstAliveSuccessor()
		}
	}
}

// Chord Utils

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
}

func (n *RingNode) GetFirstAliveSuccessor() *pb.Node {
	candidates := append(n.Successors, n.MakeNode())

	for _, node := range candidates {

		if node.Id == n.Id {
			return node
		}

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
		return true
	}

	return false
}

// gRPC Client
func (n *RingNode) GetClient(addr string) (pb.ChordServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(utils.ChangePort(addr, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to node: ", err)
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
		fmt.Println("I'm the first node in the network")
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
