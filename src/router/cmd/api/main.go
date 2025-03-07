package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
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
	serversMu       sync.Mutex
	servers         = make(map[string]*models.DiscoveredServer)
	routerCert      tls.Certificate
	caCertPool      *x509.CertPool
	tlsConfigServer *tls.Config
	tlsConfigClient *tls.Config
)

func init() {
	var err error
	routerCert, err = tls.LoadX509KeyPair("./certs/router.crt", "./certs/router.key")
	if err != nil {
		log.Fatalf("Error cargando certificado del router: %v", err)
	}

	caCert, err := os.ReadFile("./certs/ca.crt")
	if err != nil {
		log.Fatalf("Error leyendo CA: %v", err)
	}
	caCertPool = x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfigServer = &tls.Config{
		Certificates: []tls.Certificate{routerCert},
		ClientAuth:   tls.NoClientCert, // ⬅️ Exige certificado del cliente
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12, // Usa TLS 1.2 o superior
	}

	tlsConfigClient = &tls.Config{
		Certificates:       []tls.Certificate{routerCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // Disable certificate verification
		// No se usa ClientAuth aquí porque el router actúa como cliente
	}
}

func main() {

	// init() is automatically called by the Go runtime

	// Create a context that cancels on SIGINT or SIGTERM.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Use a WaitGroup to wait for goroutines to finish.
	var wg sync.WaitGroup

	srv := &http.Server{
		Addr:      ":8080", // Puerto para recibir requests del front
		Handler:   http.HandlerFunc(ServeHTTP),
		TLSConfig: tlsConfigServer,
	}

	// Start the HTTP server in a separate goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("HTTPS server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTPS server error: %v", err)
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

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfigClient,
		},
	}

	candidates, _ := filterServer()
	fmt.Print("Handle request from client addr", r.RemoteAddr, " to ", r.URL.Path, "\n")

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

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	for _, srv := range targetServers {
		go func(server string) {
			req := r.Clone(ctx)
			req.RequestURI = ""
			req.URL = &url.URL{
				Scheme: "https",
				Host:   fmt.Sprintf("%s:%d", server, 8080),
				Path:   req.URL.Path,
			}

			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			req.ContentLength = int64(len(bodyBytes))

			resp, err := httpClient.Do(req)

			if err != nil {
				errChan <- err
				return
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err != nil {
				errChan <- err
				return
			}
			// fmt.Println("RESPONSE", resp.Body)
			fmt.Println("RESPONSE", resp.Body)

			if resp.StatusCode < 500 {
				// Clonar la respuesta
				clonedResp := &http.Response{
					StatusCode: resp.StatusCode,
					Header:     resp.Header,
					Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
				}

				select {
				case respChan <- clonedResp:
				default:
				}
			} else {
				errChan <- fmt.Errorf("servidor %s respondió con error: %d", server, resp.StatusCode)
			}
		}(srv)
	}

	var successfulResp *http.Response
	errors := make([]error, 0, len(targetServers))

	// Esperar por la primera respuesta exitosa o todos los errores
outerLoop:
	for i := 0; i < len(targetServers); i++ {
		select {
		case resp := <-respChan:
			successfulResp = resp
			cancel() // Cancelar otras peticiones
			break outerLoop
		case err := <-errChan:
			errors = append(errors, err)
		}
	}

	if successfulResp != nil {
		defer successfulResp.Body.Close()

		for k, v := range successfulResp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(successfulResp.StatusCode)

		if _, err := io.Copy(w, successfulResp.Body); err != nil {
			log.Printf("Error copying response body: %v", err)
		}
		return
	}

	// Todos fallaron
	log.Printf("Todos los servidores fallaron: %v", errors)
	http.Error(w, "Todos los servidores fallaron", http.StatusBadGateway)
}

/// REFACTOR PARA LUEGO

// func ServeHTTP(w http.ResponseWriter, r *http.Request) {
//     candidates, _ := filterServer()
//     log.Printf("Handle request from %s to %s", r.RemoteAddr, r.URL.Path)

//     if len(candidates) == 0 {
//         http.Error(w, "No servers available", http.StatusServiceUnavailable)
//         return
//     }

//     // 1. Configuración de timeouts
//     ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
//     defer cancel()

//     // 2. Selección de servidores con lógica de preferencia
//     targetServers := selectServers(candidates, 3)

//     // 3. Canal único para resultados
//     resultChan := make(chan *serverResult, len(targetServers))

//     // 4. Grupo de espera para todas las goroutines
//     var wg sync.WaitGroup
//     wg.Add(len(targetServers))

//     for _, srv := range targetServers {
//         go func(server string) {
//             defer wg.Done()
//             result := &serverResult{server: server}

//             // 5. Validar dirección del servidor
//             if !isValidServer(server) {
//                 result.err = fmt.Errorf("dirección inválida: %s", server)
//                 resultChan <- result
//                 return
//             }

//             // 6. Crear cliente con timeout personalizado
//             client := &http.Client{
//                 Timeout: 10 * time.Second,
//                 Transport: &http.Transport{
//                     DisableCompression:  false,
//                     MaxIdleConnsPerHost: 10,
//                 },
//             }

//             // 7. Construir request de manera segura
//             req, err := http.NewRequestWithContext(ctx, r.Method, fmt.Sprintf("http://%s%s", server, r.URL.Path), r.Body)
//             if err != nil {
//                 result.err = err
//                 resultChan <- result
//                 return
//             }

//             // 8. Copiar headers importantes
//             copyHeaders(req.Header, r.Header)

//             resp, err := client.Do(req)
//             if err != nil {
//                 result.err = fmt.Errorf("error de conexión: %w", err)
//                 resultChan <- result
//                 return
//             }
//             defer resp.Body.Close()

//             // 9. Leer solo si es respuesta válida
//             if resp.StatusCode < 500 {
//                 result.resp = resp
//                 var buf bytes.Buffer
//                 tee := io.TeeReader(resp.Body, &buf)

//                 // 10. Leer y mantener buffer para posibles reintentos
//                 body, _ := io.ReadAll(tee)
//                 result.body = body
//                 result.headers = resp.Header
//                 result.statusCode = resp.StatusCode

//                 // 11. Reconstruir body para posibles lecturas
//                 resp.Body = io.NopCloser(&buf)
//             } else {
//                 result.err = fmt.Errorf("código de error: %d", resp.StatusCode)
//             }

//             resultChan <- result
//         }(srv)
//     }

//     // 12. Esperar resultados en goroutine separada
//     go func() {
//         wg.Wait()
//         close(resultChan)
//     }()

//     // 13. Seleccionar primera respuesta exitosa
//     var bestResponse *serverResult
//     for result := range resultChan {
//         if result.resp != nil && bestResponse == nil {
//             bestResponse = result
//             cancel() // Cancelar otras solicitudes
//         } else if result.err != nil {
//             log.Printf("Error del servidor %s: %v", result.server, result.err)
//         }
//     }

//     // 14. Escribir respuesta
//     if bestResponse != nil {
//         copyHeaders(w.Header(), bestResponse.headers)
//         w.WriteHeader(bestResponse.statusCode)
//         w.Write(bestResponse.body)
//         return
//     }

//     // 15. Manejo de errores detallado
//     log.Printf("All servers failed for %s", r.URL.Path)
//     http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
// }

// // Estructura para manejar resultados
// type serverResult struct {
//     server     string
//     resp       *http.Response
//     body       []byte
//     headers    http.Header
//     statusCode int
//     err        error
// }

// // Función para copiar headers
// func copyHeaders(dst, src http.Header) {
//     for k, vv := range src {
//         for _, v := range vv {
//             dst.Add(k, v)
//         }
//     }
// }

// // Validación de servidores
// func isValidServer(addr string) bool {
//     host, port, err := net.SplitHostPort(addr)
//     if err != nil || port == "" {
//         return false
//     }
//     return net.ParseIP(host) != nil || isValidHostname(host)
// }

// // Lógica mejorada de selección de servidores
// func selectServers(servers []string, max int) []string {
//     if len(servers) <= max {
//         return servers
//     }

//     // Implementar lógica de selección inteligente
//     return servers[:max]
// }
