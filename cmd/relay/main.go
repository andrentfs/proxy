package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/yamux"
)

var (
	currentSession *yamux.Session
	sessionMutex   sync.Mutex
)

func main() {
	// 1. Listen para conexões de Pods (SOCKS5/1080) - UMA ÚNICA VEZ
	proxyListener, err := net.Listen("tcp", ":1080")
	if err != nil {
		log.Fatalf("Erro ao abrir porta 1080: %v", err)
	}
	log.Println("Relay: Porta 1080 aberta para tráfego SOCKS5.")

	go func() {
		for {
			podConn, err := proxyListener.Accept()
			if err != nil {
				log.Printf("Erro accept pod: %v", err)
				continue
			}
			go handlePodConnection(podConn)
		}
	}()

	// 2. Listen para o Agente Local (Túnel/8080)
	tunnelListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Erro tunnel listener: %v", err)
	}
	log.Println("Relay: Aguardando Agente Local na 8080...")

	for {
		conn, err := tunnelListener.Accept()
		if err != nil {
			log.Printf("Erro accept tunnel: %v", err)
			continue
		}
		log.Println("Relay: Agente Local conectado! Iniciando nova sessão yamux...")

		setupNewSession(conn)
	}
}

func setupNewSession(conn net.Conn) {
	yamuxCfg := yamux.DefaultConfig()
	yamuxCfg.KeepAliveInterval = 25 * time.Second
	yamuxCfg.ConnectionWriteTimeout = 60 * time.Second

	session, err := yamux.Server(conn, yamuxCfg)
	if err != nil {
		log.Printf("Erro ao criar sessão yamux: %v", err)
		return
	}

	sessionMutex.Lock()
	if currentSession != nil {
		currentSession.Close() // Fecha a anterior se existir
	}
	currentSession = session
	sessionMutex.Unlock()

	// Mantém a conexão viva até ela cair
	<-session.CloseChan()
	log.Println("Relay: Sessão yamux encerrada.")
}

func handlePodConnection(podConn net.Conn) {
	sessionMutex.Lock()
	session := currentSession
	sessionMutex.Unlock()

	if session == nil || session.IsClosed() {
		log.Println("Relay: Erro - Nenhuma sessão de túnel ativa para o Pod.")
		podConn.Write([]byte("Proxy Error: No tunnel connected"))
		podConn.Close()
		return
	}

	log.Println("Relay: Abrindo stream no túnel para nova conexão de Pod...")
	stream, err := session.Open()
	if err != nil {
		log.Printf("Erro ao abrir stream no túnel: %v", err)
		podConn.Close()
		return
	}

	defer podConn.Close()
	defer stream.Close()

	done := make(chan struct{}, 2)
	go func() { io.Copy(stream, podConn); done <- struct{}{} }()
	go func() { io.Copy(podConn, stream); done <- struct{}{} }()
	<-done
}
