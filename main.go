package main

import (
	"bufio"
	"fmt"
	"net"

	"strings"
	"sync"
	"time"
)

const (
	port           = ":8989"
	maxConnections = 10
)

var (
	clients    = make(map[net.Conn]string)
	clientLock = sync.Mutex{}
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		if len(clients) >= maxConnections {
			conn.Close()
			continue
		}

		clientLock.Lock()
		clients[conn] = ""
		clientLock.Unlock()

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	conn.Write([]byte("Welcome to TCP-Chat!\n[ENTER YOUR NAME]: "))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	clientLock.Lock()
	clients[conn] = name
	clientLock.Unlock()

	broadcast(fmt.Sprintf("%s has joined the chat...", name), conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}
		broadcast(fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), name, message), conn)
	}

	clientLock.Lock()
	delete(clients, conn)
	clientLock.Unlock()

	broadcast(fmt.Sprintf("%s has left the chat...", name), conn)
}

func broadcast(message string, sender net.Conn) {
	clientLock.Lock()
	defer clientLock.Unlock()

	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message + "\n"))
		}
	}
}
