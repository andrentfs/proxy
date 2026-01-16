package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/hashicorp/yamux"
)

func main() {
	relayAddr := flag.String("relay", "", "Endereço do Relay no K8s (ex: 1.2.3.4:8080)")
	port := flag.Int("port", 1080, "Porta para o servidor SOCKS5 local")
	flag.Parse()

	if *relayAddr != "" {
		runTunnelClient(*relayAddr)
	} else {
		runLocalServer(*port)
	}
}

func runLocalServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Falha ao iniciar listener: %v", err)
	}
	defer listener.Close()
	log.Printf("Modo Local: Servidor Proxy SOCKS5 rodando na porta %d...", port)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro accept: %v", err)
			continue
		}
		go handleSOCKS5(clientConn)
	}
}

func runTunnelClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Erro ao conectar no Relay %s: %v", addr, err)
	}
	log.Printf("Modo Túnel: Conectado ao Relay em %s. Aguardando streams...", addr)

	session, err := yamux.Client(conn, nil)
	if err != nil {
		log.Fatalf("Erro yamux client: %v", err)
	}

	for {
		stream, err := session.Accept()
		if err != nil {
			log.Printf("Erro accept stream: %v", err)
			break
		}
		log.Println("Túnel: Novo stream recebido do K8s, tratando como SOCKS5.")
		go handleSOCKS5(stream)
	}
}

func handleSOCKS5(conn net.Conn) {
	defer conn.Close()

	// 1. Handshake
	buf := make([]byte, 256)
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		return
	}
	if buf[0] != 0x05 {
		return
	}
	nMethods := int(buf[1])
	if _, err := io.ReadFull(conn, buf[:nMethods]); err != nil {
		return
	}
	conn.Write([]byte{0x05, 0x00})

	// 2. Request
	if _, err := io.ReadFull(conn, buf[:4]); err != nil {
		return
	}
	if buf[1] != 0x01 { // CONNECT
		return
	}

	var target string
	switch buf[3] {
	case 0x01: // IPv4
		if _, err := io.ReadFull(conn, buf[:4]); err != nil {
			return
		}
		target = net.IP(buf[:4]).String()
	case 0x03: // Domain
		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			return
		}
		length := int(buf[0])
		if _, err := io.ReadFull(conn, buf[:length]); err != nil {
			return
		}
		target = string(buf[:length])
	case 0x04: // IPv6
		if _, err := io.ReadFull(conn, buf[:16]); err != nil {
			return
		}
		target = net.IP(buf[:16]).String()
	}

	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		return
	}
	targetPort := (uint16(buf[0]) << 8) | uint16(buf[1])
	targetAddr := net.JoinHostPort(target, fmt.Sprintf("%d", targetPort))

	// 3. Dial Target
	log.Printf("SOCKS5: Conectando a %s", targetAddr)
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		conn.Write([]byte{0x05, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}
	defer targetConn.Close()
	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})

	// 4. Pipe
	done := make(chan bool, 2)
	go func() { io.Copy(targetConn, conn); done <- true }()
	go func() { io.Copy(conn, targetConn); done <- true }()
	<-done
}
