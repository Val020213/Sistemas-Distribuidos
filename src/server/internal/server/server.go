package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"server/internal/chord"
)

type Server struct {
	node *chord.RingNode
}

func NewServer(node *chord.RingNode) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(node.CaCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{node.Cert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // ⬅️ Exige certificado del router
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	NewServer := &Server{
		node: node,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewServer.RegisterRoutes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
