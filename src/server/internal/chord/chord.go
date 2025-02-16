package chord

import (
	"context"
	"os"
	pb "server/internal/chord/chordpb"
	"server/internal/scraper"

	"server/internal/utils"
	"sync"

	"google.golang.org/grpc"
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

func (n *RingNode) StartRPCServer(grpcServer *grpc.Server) {
	pb.RegisterChordServiceServer(grpcServer, n)
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
