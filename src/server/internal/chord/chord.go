package chord

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

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
	grpcAddr = os.Getenv("IP_ADDRESS")
	grpcPort = os.Getenv("RPC_PORT")
	mBits    = utils.GetEnvAsInt("CHORD_BITS", 8)
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

func (n *RingNode) Discover(ctx context.Context) (*pb.DiscoveryResponse, error) {
	log.Printf("Nodo %d: Recibido discover request", n.ID)

	return &pb.DiscoveryResponse{
		Id:      n.ID,
		Address: n.Address,
	}, nil
}

func (n *RingNode) StartRPCServer(grpcServer *grpc.Server) {

	fmt.Println("RPC Server started")

	pb.RegisterChordServiceServer(grpcServer, n)

	_, err := n.Lookup()
	if err != nil {
		fmt.Println(err)
	}
}

func (n *RingNode) Notify(ctx context.Context, req *pb.NotifyRequest) (*pb.NotifyResponse, error) {
	newPredecessor := &RemoteNode{
		ID:      req.GetId(),
		Address: req.GetAddress(),
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if n.Predecessor == nil || utils.Between(newPredecessor.ID, n.Predecessor.ID, n.ID) {
		n.Predecessor = newPredecessor
		return &pb.NotifyResponse{Updated: true}, nil
	}

	return &pb.NotifyResponse{Updated: false}, nil
}

func (n *RingNode) Health(ctx context.Context, empty *pb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Id:      n.ID,
		Address: n.Address,
	}, nil
}

func (n *RingNode) Lookup() (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChan := make(chan string, 1)
	errChan := make(chan error)

	var wg sync.WaitGroup

	for i := 1; i < 54; i++ {
		addr := fmt.Sprintf("10.0.11.2%d:50051", i)

		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			conn, err := grpc.NewClient(a, grpc.WithTransportCredentials(insecure.NewCredentials()))

			if err != nil {
				errChan <- err
				return
			}

			defer conn.Close()

			// Health check	service from gRPC
			healthClient := healthpb.NewHealthClient(conn)
			resp, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{Service: "chord"})

			if err == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING {
				n.mu.Lock()
				defer n.mu.Unlock()

				select {
				case resultChan <- a:
					cancel()
				default:
				}
			}

		}(addr)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case addr := <-resultChan:
		return addr, nil
	case <-ctx.Done():
		return "", fmt.Errorf("no hay nodos disponibles")
	}

}
