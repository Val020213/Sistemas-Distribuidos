package chord

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "server/internal/chord/chordpb"
	"server/internal/scraper"
	"server/internal/utils"
)

type RingNode struct {
	Id          uint64            // Node's Id (computed from its address)
	Address     string            // Host
	Port        string            // Port
	Successor   *RemoteNode       // Immediate successor in the ring
	Predecessor *RemoteNode       // Immediate predecessor in the ring
	Finger      []*RemoteNode     // Finger table entries
	Data        map[string]string // Simple key-value storage
	m           int               // Number of bits in the hash space
	mu          sync.Mutex        // Protects access to mutable fields
	pb.UnimplementedChordServiceServer

	Scraper *scraper.Scraper // = Scrapper Application
}

type RemoteNode struct {
	Id      uint64
	Address string
}

var (
	grpcAddr      = os.Getenv("IP_ADDRESS")
	grpcPort      = os.Getenv("RPC_PORT")
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
		m:       mBits,
		Data:    make(map[string]string),
		Finger:  make([]*RemoteNode, mBits),
		Scraper: scraper,
	}
}

func (n *RingNode) StartRPCServer(grpcServer *grpc.Server) {
	pb.RegisterChordServiceServer(grpcServer, n)
	fmt.Println("Starting gRPC Server on ", n.Address, ":", n.Port)
	n.initNode(n.m)
	n.joinNetwork()

}

func (n *RingNode) initNode(mBits int) {
	remoteNode := &RemoteNode{Id: n.Id, Address: n.Address}
	for i := 0; i < mBits; i++ {
		n.Finger[i] = remoteNode
	}
	n.Predecessor = remoteNode
	n.Successor = remoteNode
}

func (n *RingNode) Notify(ctx context.Context, node *pb.Node) (*pb.Successful, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	fmt.Println("Notified by ", node.Address)

	if n.Predecessor == nil || utils.BetweenRightInclusive(node.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = &RemoteNode{Id: node.Id, Address: node.Address}
		fmt.Println("Updated predecessor to ", node.Address)
		return &pb.Successful{Successful: true}, nil
	}

	return &pb.Successful{Successful: false}, nil
}

func (n *RingNode) Health(ctx context.Context, empty *pb.Empty) (*pb.HealthResponse, error) {
	fmt.Printf("Debugging node:\naddr: %s\nid: %d\nsucc: %s\npred: %s\nfinger table: %v\n", n.Address, n.Id, n.Successor.Address, n.Predecessor.Address, n.Finger)
	return &pb.HealthResponse{
		Id:      n.Id,
		Address: n.Address,
	}, nil
}

func (n *RingNode) FindSuccessor(ctx context.Context, req *pb.KeyRequest) (*pb.Node, error) {
	fmt.Println("Finding successor for ", req.Key)

	key := req.Key

	if utils.BetweenRightInclusive(key, n.Id, n.Successor.Id) {
		fmt.Println("Successor found", "Id:", n.Successor.Id, " Address:", n.Successor.Address)
		return &pb.Node{Id: n.Successor.Id, Address: n.Successor.Address}, nil
	}

	fmt.Println("Finding closest preceding node")
	nextNode, _ := n.closestPrecedingNode(key)

	conn, err := grpc.NewClient(utils.ChangePort(nextNode.Address, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Error connecting to next node: ", err)
		return nil, err
	}

	defer conn.Close()

	client := pb.NewChordServiceClient(conn)
	fmt.Println("Asking next node for successor", "Id:", nextNode.Id, " Address:", nextNode.Address)
	return client.FindSuccessor(ctx, req)
}

func (n *RingNode) closestPrecedingNode(key uint64) (*pb.Node, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Println("Finding closest preceding node for key:", key)

	for i := len(n.Finger) - 1; i >= 0; i-- {
		if n.Finger[i] == nil {
			continue
		}

		if utils.BetweenRightInclusive(n.Finger[i].Id, n.Id, key) {
			fmt.Println("Closest preceding node found ", "Id:", n.Finger[i].Id, " Address:", n.Finger[i].Address)
			return &pb.Node{Id: n.Finger[i].Id, Address: n.Finger[i].Address}, nil
		}
	}
	fmt.Println("Closest preceding node found ", "Id:", n.Successor.Id, " Address:", n.Successor.Address)
	return &pb.Node{Id: n.Successor.Id, Address: n.Successor.Address}, nil
}

func (n *RingNode) joinNetwork() (string, error) {

	fmt.Println("Joining network...")

	addr, err := n.Discover()

	if err != nil {
		return "", fmt.Errorf("error discovering network: %v", err)
	}

	if addr == n.Address {
		fmt.Println("I'm the first node in the network")
		return n.Address, nil
	}

	conn, err := grpc.NewClient(utils.ChangePort(addr, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Error connecting to node: ", err)
		return n.joinNetwork()
	}
	defer conn.Close()

	client := pb.NewChordServiceClient(conn)
	succ, err := client.FindSuccessor(context.Background(), &pb.KeyRequest{Key: n.Id})

	if err != nil {
		fmt.Println("Error en la conexión: ", err)
		return "", errors.New("no se encontraron nodos para enlazar")
	}

	fmt.Println("Successor encontrado en ", succ.Address)
	n.Successor = &RemoteNode{Id: succ.Id, Address: succ.Address}
	fmt.Println("Notifying successor")
	client.Notify(context.Background(), &pb.Node{Id: n.Id, Address: n.Address})
	return n.Successor.Address, nil
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
