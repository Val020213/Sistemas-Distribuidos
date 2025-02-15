package chord

import (
	pb "server/internal/chord/chordpb"

	"server/internal/utils"
	"sync"
)

type RingNode struct {
	ID          uint64            // Node's ID (computed from its address)
	Address     string            // Host:port (used for RPC)
	Successor   *RemoteNode       // Immediate successor in the ring
	Predecessor *RemoteNode       // Immediate predecessor in the ring
	Finger      []*RemoteNode     // Finger table entries
	Data        map[string]string // Simple key-value storage
	m           int               // Number of bits in the hash space

	mu sync.Mutex // Protects access to mutable fields
}

type RemoteNode struct {
	ID      uint64
	Address string
}

func NewNode(grpcAddr, restAddr string, mBits int) *RingNode {
	id := utils.ChordHash(grpcAddr, mBits)

	return &RingNode{

		ID:      id,
		Address: restAddr,
		m:       mBits,
		Data:    make(map[string]string),
		Finger:  make([]*RemoteNode, mBits),
	}
}

func (n *RingNode) Notify(newPredecessor *RemoteNode) (*pb.NotifyResponse, error) {

	n.mu.Lock()
	defer n.mu.Unlock()

	if n.Predecessor == nil || utils.Between(newPredecessor.ID, n.Predecessor.ID, n.ID) {
		n.Predecessor = newPredecessor
		return &pb.NotifyResponse{Updated: true}, nil
	}

	return &pb.NotifyResponse{Updated: false}, nil
}

func (n *RingNode) Health() (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Id:      n.ID,
		Address: n.Address,
	}, nil
}
