package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"router/internal/models"
)

const (
	// Multicast address that must match the one used by the servers.
	multicastAddr = "224.0.0.1:9999"
	serverTimeout = 15 * time.Second // Timeout for server expiration.
)

// serversMu protects the map of discovered servers.
var (
	serversMu sync.Mutex
	servers   = make(map[string]*models.DiscoveredServer)
)

func main() {
	// Create a context that cancels on SIGINT or SIGTERM.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Use a WaitGroup to wait for goroutines to finish.
	var wg sync.WaitGroup

	// Start the multicast listener in a separate goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		listenMulticast(ctx)
	}()

	log.Println("Router started. Waiting for multicast messages...")

	// Block until a termination signal is received.
	<-ctx.Done()
	log.Println("Termination signal received. Shutting down...")

	// Wait for all goroutines to complete.
	wg.Wait()
	log.Println("Graceful shutdown complete.")
}

// listenMulticast listens for UDP multicast messages and updates the list of discovered servers.
func listenMulticast(ctx context.Context) {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddr)
	if err != nil {
		log.Fatalf("Error resolving multicast address: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("Error listening on multicast: %v", err)
	}
	defer conn.Close()

	// Set a read deadline to periodically check the context.
	setReadDeadline(conn)

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			log.Println("Multicast listener: shutting down.")
			return
		default:
			n, src, err := conn.ReadFromUDP(buf)
			if err != nil {
				// If timeout, reset the deadline and continue.
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					setReadDeadline(conn)
					continue
				}
				log.Printf("Error reading multicast: %v", err)
				continue
			}

			message := string(buf[:n])
			// log.Printf("Received multicast message from %s: %s", src, message)
			const prefix = "SERVER:"
			if len(message) > len(prefix) && message[:len(prefix)] == prefix {
				srvAddr := message[len(prefix):]
				serversMu.Lock()
				servers[srvAddr] = &models.DiscoveredServer{
					Address:  srvAddr,
					LastSeen: time.Now(),
				}
				serversMu.Unlock()
				log.Printf("Discovered server from %s: %s", src, srvAddr)
			}
		}
	}
}

// setReadDeadline sets a deadline for reading from the UDP connection.
func setReadDeadline(conn *net.UDPConn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
}

func chooseServer() (*models.DiscoveredServer, error) {
	serversMu.Lock()
	defer serversMu.Unlock()

	// Remove expired servers.
	now := time.Now()
	for key, srv := range servers {
		if now.Sub(srv.LastSeen) > serverTimeout {
			log.Printf("Router: removing expired server: %s", key)
			delete(servers, key)
		}
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("no available servers")
	}
	// Select the first server from the map.
	for _, srv := range servers {
		fmt.Println("Server: ", srv)
		return srv, nil
	}
	return nil, fmt.Errorf("no available servers")
}
