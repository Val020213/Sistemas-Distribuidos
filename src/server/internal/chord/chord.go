package chord

import (
	"context"
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
	ID          uint64            // Node's ID (computed from its address)
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
	ID      uint64
	Address string
}

var (
	grpcAddr      = os.Getenv("IP_ADDRESS")
	grpcPort      = os.Getenv("RPC_PORT")
	mBits         = utils.GetEnvAsInt("CHORD_BITS", 8)
	multicastAddr = "224.0.0.1:9999"
)

func NewNode() *RingNode {

	id := utils.ChordHash(grpcAddr, mBits)
	scraper := scraper.NewScraper()
	return &RingNode{
		ID:      id,
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

func (n *RingNode) Notify(ctx context.Context, req *pb.Node) (*pb.Succesfull, error) {
	newPredecessor := &RemoteNode{
		ID:      req.GetId(),
		Address: req.GetAddress(),
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if n.Predecessor == nil || utils.Between(newPredecessor.ID, n.Predecessor.ID, n.ID) {
		n.Predecessor = newPredecessor
		return &pb.Succesfull{succesfull: true}, nil
	}

	return &pb.Succesfull{succesfull: false}, nil
}

func (n *RingNode) Health(ctx context.Context, empty *pb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Id:      n.ID,
		Address: n.Address,
	}, nil
}

func (n *RingNode) FindSuccessor(ctx context.Context, key *pb.KeyRequest) (*pb.Node, error) {
	
}

func (n *RingNode) joinNetwork() (string, error) {

	addr, _ := n.Discover()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return n.joinNetwork()
	}
	defer conn.Close()

	client := pb.NewChordServiceClient(conn)
	succ, err := client.

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
		// Si el mensaje recibido proviene de un nodo distinto (comparando direcciones)
		if src.String() != n.Address {
			// Se asume que el mensaje contiene el identificador del nodo responsable.
			responsable := string(buf[:nbytes])
			return responsable, nil
		}
		// Si el mensaje es de si mismo, continúa esperando.
	}
}
