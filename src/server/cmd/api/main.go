package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/internal/chord"
	"server/internal/multicast"
	"server/internal/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func gracefulShutdown(apiServer *http.Server, grpcServer *grpc.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down gracefully")

	// Shutdown HTTP
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer httpCancel()
	if err := apiServer.Shutdown(httpCtx); err != nil {
		log.Printf("Error HTTP shutdown: %v", err)
	}

	// Apagar gRPC
	grpcServer.GracefulStop()

	done <- true
}

func main() {

	cert, err := tls.LoadX509KeyPair("./certs/server.crt", "./certs/server.key")
	if err != nil {
		fmt.Println("failed to load key pair: ", err)
	}

	caCert, err := os.ReadFile("./certs/ca.crt")

	if err != nil {
		fmt.Println("failed to load CA certificate: ", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // ⬅️ Exige mTLS
		ClientCAs:    caCertPool,
	})

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	node := chord.NewNode(cert, caCert)
	node.StartRPCServer(grpcServer)

	httpServer := server.NewServer(node)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(httpServer, grpcServer, done)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", node.Address, node.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("gRPC server initialized on %s:%s", node.Address, node.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC failed: %v", err)
		}
	}()

	go multicast.MulticastAnnouncer()

	err = httpServer.ListenAndServeTLS("", "")
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
