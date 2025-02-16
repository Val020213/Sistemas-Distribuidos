package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
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

	srv := &http.Server{
		Addr:    ":8080", // Puerto para recibir requests del front
		Handler: http.HandlerFunc(ServeHTTP),
	}

	// Start the HTTP server in a separate goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("HTTP server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start the multicast listener in a separate goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Router started. Waiting for multicast messages...")
		listenMulticast(ctx)
	}()

	// Block until a termination signal is received.
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP shutdown error: %v", err)
		}
	}()

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

func filterServer() ([]string, error) {
	filtered := []string{}
	for key, srv := range servers {
		if time.Since(srv.LastSeen) > serverTimeout {
			log.Printf("Router: removing expired server: %s", key)
			continue
		}
		filtered = append(filtered, key)
	}
	return filtered, nil
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	candidates, _ := filterServer()

	if len(candidates) == 0 {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		return
	}

	targetServers := candidates
	if len(targetServers) > 3 {
		targetServers = targetServers[:3]
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	respChan := make(chan *http.Response, 1)
	errChan := make(chan error, len(targetServers))

	for _, srv := range targetServers {
		go func(server string) {
			req := r.Clone(ctx)
			req.URL = &url.URL{
				Scheme: "http",
				Host:   server,
				Path:   req.URL.Path,
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errChan <- err
				return
			}

			if resp.StatusCode < 500 {
				select {
				case respChan <- resp:
				default:
					resp.Body.Close()
				}
			} else {
				resp.Body.Close()
				errChan <- fmt.Errorf("servidor %s respondió con error: %d", server, resp.StatusCode)
			}
		}(srv)
	}

	var successfulResp *http.Response
	errors := make([]error, 0, len(targetServers))

	// Esperar por la primera respuesta exitosa o todos los errores
	for i := 0; i < len(targetServers); i++ {
		select {
		case resp := <-respChan:
			successfulResp = resp
			cancel() // Cancelar otras peticiones
			break
		case err := <-errChan:
			errors = append(errors, err)
		}
	}

	if successfulResp != nil {
		defer successfulResp.Body.Close()
		// Copiar headers
		for k, v := range successfulResp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(successfulResp.StatusCode)
		io.Copy(w, successfulResp.Body)
		return
	}

	// Todos fallaron
	log.Printf("Todos los servidores fallaron: %v", errors)
	http.Error(w, "Todos los servidores fallaron", http.StatusBadGateway)
}
