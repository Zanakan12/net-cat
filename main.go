package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	DefaultPort = "8989"
	MaxClients  = 10
	LogFile     = "server.log"
	LinuxLogo   = `
          .--.
         |o_o |
         |:_/ |
        //   \ \
       (|     | )
      /'\_   _/` + "`\\ " + `
      \___)=(___/
`
)

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

// Message struct holds message details.
type Message struct {
	Timestamp time.Time
	Client    string
	Content   string
}

// Client struct represents connected clients.
type Client struct {
	Conn     net.Conn
	Username string
	Out      chan string
}

// Server struct holds the server state.
type Server struct {
	Protocol    Protocol
	Port        string
	Clients     map[string]*Client
	Messages    []Message
	ClientsLock sync.Mutex
	MsgLock     sync.Mutex
	LogFile     *os.File
}

// NewServer creates a new server instance.
func NewServer(protocol Protocol, port string) *Server {
	file, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}

	return &Server{
		Protocol: protocol,
		Port:     port,
		Clients:  make(map[string]*Client),
		Messages: []Message{},
		LogFile:  file,
	}
}

// Start initiates the server based on the protocol (TCP or UDP).
func (s *Server) Start() {
	if s.Protocol == UDP {
		s.startUDP()
	} else {
		s.startTCP()
	}
}

// startTCP starts a TCP server and handles connections.
func (s *Server) startTCP() {
	listener, err := net.Listen(string(TCP), ":"+s.Port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()
	log.Printf("Listening on port %s with TCP", s.Port)

	for {
		if len(s.Clients) >= MaxClients {
			log.Println("Max clients connected. Rejecting new connection.")
			conn, err := listener.Accept()
			if err == nil {
				conn.Write([]byte("Server is full. Try again later.\n"))
				conn.Close()
			}
			continue
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleClient(conn)
	}
}

// startUDP starts a UDP server and handles incoming messages.
func (s *Server) startUDP() {
	udpAddr, err := net.ResolveUDPAddr(string(UDP), ":"+s.Port)
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}

	conn, err := net.ListenUDP(string(UDP), udpAddr)
	if err != nil {
		log.Fatalf("Error starting UDP server: %v", err)
	}
	defer conn.Close()

	log.Printf("Listening on port %s with UDP", s.Port)

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading UDP data: %v", err)
			continue
		}
		message := string(buf[:n])
		fmt.Printf("[%s]: %s\n", addr, message)
	}
}

// handleClient manages the interaction with a TCP client.
func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte(LinuxLogo))
	conn.Write([]byte("Enter your name: "))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	username := strings.TrimSpace(string(buf[:n]))
	if username == "" {
		conn.Write([]byte("Invalid username. Disconnecting...\n"))
		return
	}

	client := &Client{
		Conn:     conn,
		Username: username,
		Out:      make(chan string, 100), // Increased buffer size even further
	}

	s.ClientsLock.Lock()
	if _, exists := s.Clients[username]; exists {
		s.ClientsLock.Unlock()
		conn.Write([]byte("Username already taken. Disconnecting...\n"))
		return
	}
	s.Clients[username] = client
	s.ClientsLock.Unlock()

	s.logActivity(fmt.Sprintf("Client %s joined.", username))
	s.broadcast(fmt.Sprintf("[INFO]: %s joined the chat\n", username), "INFO")

	s.MsgLock.Lock()
	for _, msg := range s.Messages {
		conn.Write([]byte(fmt.Sprintf("[%s][%s]: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Client, msg.Content)))
	}
	s.MsgLock.Unlock()

	go s.sendMessagesToClient(client)
	s.receiveMessagesFromClient(client)

	s.ClientsLock.Lock()
	delete(s.Clients, username)
	s.ClientsLock.Unlock()

	s.broadcast(fmt.Sprintf("[INFO]: %s left the chat\n", username), "INFO")
	s.logActivity(fmt.Sprintf("Client %s left.", username))
}

// sendMessagesToClient sends messages to a specific client.
func (s *Server) sendMessagesToClient(client *Client) {
	for msg := range client.Out {
		_, err := client.Conn.Write([]byte(msg))
		if err != nil {
			return
		}
	}
}

// receiveMessagesFromClient listens for incoming messages from a client, including the name change command.
func (s *Server) receiveMessagesFromClient(client *Client) {
	buf := make([]byte, 1024)
	for {
		n, err := client.Conn.Read(buf)
		if err != nil {
			return
		}

		message := strings.TrimSpace(string(buf[:n]))

		// Handle name change command
		if strings.HasPrefix(message, "/name ") {
			newName := strings.TrimSpace(strings.TrimPrefix(message, "/name "))
			if newName == "" {
				client.Conn.Write([]byte("Invalid new name.\n"))
				continue
			}

			// Ensure the new name isn't already taken
			s.ClientsLock.Lock()
			if _, exists := s.Clients[newName]; exists {
				client.Conn.Write([]byte("This name is already taken.\n"))
				s.ClientsLock.Unlock()
				continue
			}

			// Broadcast the name change
			oldName := client.Username
			delete(s.Clients, client.Username) // Remove the old name
			client.Username = newName          // Update the name
			s.Clients[newName] = client        // Add the new name

			s.ClientsLock.Unlock()

			// Notify others of the name change
			s.broadcast(fmt.Sprintf("[INFO]: %s changed their name to %s\n", oldName, newName), "INFO")
			s.logActivity(fmt.Sprintf("Client %s changed their name to %s", oldName, newName))

			continue
		}

		if message == "/exit" {
			return
		}

		timestamp := time.Now()
		msg := Message{Timestamp: timestamp, Client: client.Username, Content: message}
		s.MsgLock.Lock()
		s.Messages = append(s.Messages, msg)
		s.MsgLock.Unlock()

		formattedMsg := fmt.Sprintf("[%s][%s]: %s\n", timestamp.Format("2006-01-02 15:04:05"), client.Username, message)
		s.broadcast(formattedMsg, client.Username)
	}
}

// broadcast sends a message to all clients except the sender.
func (s *Server) broadcast(message, sender string) {
	s.ClientsLock.Lock()
	defer s.ClientsLock.Unlock()

	for _, client := range s.Clients {
		if client.Username == sender {
			continue
		}
		select {
		case client.Out <- message:
		default:
			log.Printf("Client %s is slow. Dropping message.", client.Username)
		}
	}
}

// logActivity logs activities to the server's log file.
func (s *Server) logActivity(activity string) {
	log.Println(activity)
	s.LogFile.WriteString(activity + "\n")
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown() {
	s.ClientsLock.Lock()
	for _, client := range s.Clients {
		client.Conn.Close()
	}
	s.ClientsLock.Unlock()
	s.LogFile.Close()
}

func main() {
	listen := flag.Bool("l", false, "Listen for incoming connections")
	protocol := flag.String("u", string(TCP), "Choose between tcp or udp")
	flag.Parse()

	port := DefaultPort
	args := flag.Args()

	if len(args) == 1 {
		port = args[0]
	}

	if *listen || len(flag.Args()) == 0 || port != DefaultPort {
		server := NewServer(Protocol(*protocol), port)
		server.Start()
	} else {
		fmt.Println("[USAGE 1]: ./TCPChat -l -p <port> -u <tcp|udp>\n[USAGE 2]: ./TCPChat $port\n[USAGE 3]: ./TCPChat")
	}
}
