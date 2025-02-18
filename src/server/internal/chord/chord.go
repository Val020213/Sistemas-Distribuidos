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

	fmt.Println("RPC Server started")

	pb.RegisterChordServiceServer(grpcServer, n)

	address, err := n.Discover()

	fmt.Println("Address: ", address)
	fmt.Println("Error: ", err)

}

func (n *RingNode) Notify(ctx context.Context, node *pb.Node) (*pb.Succesfull, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.Predecessor == nil || utils.Between(node.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = &RemoteNode{Id: node.Id, Address: node.Address}
		return &pb.Succesfull{Succesfull: true}, nil
	}

	return &pb.Succesfull{Succesfull: false}, nil
}

func (n *RingNode) Health(ctx context.Context, empty *pb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Id:      n.Id,
		Address: n.Address,
	}, nil
}

func (n *RingNode) FindSuccessor(ctx context.Context, req *pb.KeyRequest) (*pb.Node, error) {

	key := req.Key

	if utils.Between(n.Id, key, n.Successor.Id) {
		return &pb.Node{Id: n.Successor.Id, Address: n.Successor.Address}, nil
	}

	nextNode, err := n.closestPrecedingNode(key)

	conn, err := grpc.NewClient(nextNode.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewChordServiceClient(conn)
	return client.FindSuccessor(ctx, req)

}

func (n *RingNode) closestPrecedingNode(key uint64) (*pb.Node, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	for i := len(n.Finger) - 1; i >= 0; i-- {
		if n.Finger[i] == nil {
			continue
		}

		if utils.Between(n.Finger[i].Id, n.Id, key) {
			return &pb.Node{Id: n.Finger[i].Id, Address: n.Finger[i].Address}, nil
		}
	}

	return &pb.Node{Id: n.Successor.Id, Address: n.Successor.Address}, nil
}

func (n *RingNode) joinNetwork() (string, error) {

	addr, _ := n.Discover()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return n.joinNetwork()
	}
	defer conn.Close()

	client := pb.NewChordServiceClient(conn)
	succ, err := client.FindSuccessor(context.Background(), &pb.KeyRequest{Key: n.Id})

	if err != nil {
		return "", errors.New("no se encontraron nodos para enlazar")
	}

	n.Successor = &RemoteNode{Id: succ.Id, Address: succ.Address}
	client.Notify(context.Background(), &pb.Node{Id: n.Id, Address: n.Address})
	return n.Successor.Address, nil
}

func (n *RingNode) Discover() (string, error) {

	addr, err := net.ResolveUDPAddr("udp4", multicastAddr)
	if err != nil {
		log.Fatalf("Error resolving multicast address: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("Error listening on multicast: %v", err)
	}
	defer conn.Close()

	// Buffer para recibir mensajes.
	buf := make([]byte, 1024)
	// Establece un tiempo límite de lectura de 6 segundos.
	deadline := time.Now().Add(6 * time.Second)
	conn.SetReadDeadline(deadline)

	// Se queda esperando hasta recibir un mensaje o agotar el tiempo.
	for {
		nbytes, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			// Si se alcanza el tiempo límite, se asigna a si mismo como responsable.
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return n.Address, nil
			}
			return "", fmt.Errorf("error leyendo desde multicast: %v", err)
		}
		// Si el mensaje recibIdo proviene de un nodo distinto (comparando direcciones)
		if src.String() != n.Address {
			// Se asume que el mensaje contiene el Identificador del nodo responsable.
			responsable := string(buf[:nbytes])
			return responsable, nil
		}
		// Si el mensaje es de si mismo, continúa esperando.
	}
}
