package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"server/internal/chord"
	"server/internal/server"

	"google.golang.org/grpc"
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

	node := chord.NewNode()
	grpcServer := grpc.NewServer()
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

	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
