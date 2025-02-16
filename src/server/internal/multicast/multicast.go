package multicast

import (
	"fmt"
	"net"
	"os"

	"log"
	"time"
)

const (
	multicastAddr     = "224.0.0.1:9999"
	multicastInterval = 5 * time.Second
)

// Periodically sends a multicast message with the server's IP address
func MulticastAnnouncer() {
	log.Printf("Iniciando anunciador multicast en %s", multicastAddr)
	addr, err := net.ResolveUDPAddr("udp4", multicastAddr)
	if err != nil {
		log.Fatalf("No se pudo resolver la dirección multicast: %v", err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("No se pudo conectar a la dirección multicast: %v", err)
	}
	defer conn.Close()

	for {
		message := fmt.Sprintf("SERVER:%s", os.Getenv("IP_ADDRESS"))
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Printf("Error al enviar anuncio multicast: %v", err)
		} else {
			log.Printf("Anunciado servidor: %s", message)
		}
		time.Sleep(multicastInterval)
	}
}
