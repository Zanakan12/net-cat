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
	// LinuxLogo is sent to clients upon connection
	LinuxLogo = `
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

// Message struct to hold message details
// A message consists of a timestamp, the client who sent it, and the content of the message.
type Message struct {
	Timestamp time.Time
	Client    string
	Content   string
}

// Client struct represents connected clients
// A client has a connection (Conn), a username, and a channel for outgoing messages (Out).
type Client struct {
	Conn     net.Conn
	Username string
	Out      chan string
}

// Server struct holds the server state
// This struct contains information about the protocol (TCP/UDP), the port it's listening on,
// the connected clients, chat messages, and mutexes to handle concurrency.
type Server struct {
	Protocol    Protocol
	Port        string
	Clients     map[string]*Client
	Messages    []Message
	ClientsLock sync.Mutex
	MsgLock     sync.Mutex
	LogFile     *os.File
}

// NewServer creates a new server instance
// It initializes the log file and sets up the server with the chosen protocol and port.
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

// Start initiates the server based on the protocol (TCP or UDP)
func (s *Server) Start() {
	if s.Protocol == UDP {
		s.startUDP()
	} else {
		s.startTCP()
	}
}

// startTCP starts a TCP server, accepts incoming connections and handles each client in a new goroutine
func (s *Server) startTCP() {
	listener, err := net.Listen(string(TCP), ":"+s.Port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()
	log.Printf("Listening on port %s with TCP", s.Port)

	for {
		// If the maximum number of clients is reached, reject new connections
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

		// Handle each client in a new goroutine
		go s.handleClient(conn)
	}
}

// startUDP starts a UDP server, listens for incoming messages, and prints the message along with the sender's address
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
		// Read incoming UDP messages and print them along with the sender's address
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading UDP data: %v", err)
			continue
		}

		message := string(buf[:n])
		fmt.Printf("[%s]: %s\n", addr, message)
	}
}

// handleClient manages the interaction with a newly connected TCP client
func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	// Send Linux logo to the client
	conn.Write([]byte(LinuxLogo))
	conn.Write([]byte("Enter your name: "))

	// Read the username from the client
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

	// Create a new client object
	client := &Client{
		Conn:     conn,
		Username: username,
		Out:      make(chan string),
	}

	// Add client to the server's client map
	s.ClientsLock.Lock()
	if _, exists := s.Clients[username]; exists {
		s.ClientsLock.Unlock()
		conn.Write([]byte("Username already taken. Disconnecting...\n"))
		return
	}
	s.Clients[username] = client
	s.ClientsLock.Unlock()

	// Log the new client connection and broadcast a message to other clients
	s.logActivity(fmt.Sprintf("Client %s joined.", username))
	s.broadcast(fmt.Sprintf("[INFO]: %s joined the chat\n", username), "INFO")

	// Send previous chat messages to the new client
	s.MsgLock.Lock()
	for _, msg := range s.Messages {
		conn.Write([]byte(fmt.Sprintf("[%s][%s]: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Client, msg.Content)))
	}
	s.MsgLock.Unlock()

	// Start goroutine to send messages to the client
	go s.sendMessagesToClient(client)

	// Receive messages from the client
	s.receiveMessagesFromClient(client)

	// Once the client disconnects, remove them from the client list and notify others
	s.ClientsLock.Lock()
	delete(s.Clients, username)
	s.ClientsLock.Unlock()

	s.broadcast(fmt.Sprintf("[INFO]: %s left the chat\n", username), "INFO")
	s.logActivity(fmt.Sprintf("Client %s left.", username))
}

// sendMessagesToClient sends messages to a specific client
func (s *Server) sendMessagesToClient(client *Client) {
	for msg := range client.Out {
		_, err := client.Conn.Write([]byte(msg))
		if err != nil {
			return
		}
	}
}

// receiveMessagesFromClient listens for incoming messages from a client and broadcasts them to others
func (s *Server) receiveMessagesFromClient(client *Client) {
	buf := make([]byte, 1024)
	for {
		// Read message from client
		n, err := client.Conn.Read(buf)
		if err != nil {
			return
		}

		message := string(buf[:n])
		if message == "/exit" {
			return
		}

		// Log the message and broadcast it to other clients
		timestamp := time.Now()
		msg := Message{Timestamp: timestamp, Client: client.Username, Content: message}
		s.MsgLock.Lock()
		s.Messages = append(s.Messages, msg)
		s.MsgLock.Unlock()

		formattedMsg := fmt.Sprintf("[%s][%s]: %s\n", timestamp.Format("2006-01-02 15:04:05"), client.Username, message)
		s.broadcast(formattedMsg, client.Username)
	}
}

// broadcast sends a message to all clients except the sender
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

// logActivity logs activities to the server's log file
func (s *Server) logActivity(activity string) {
	log.Println(activity)
	s.LogFile.WriteString(activity + "\n")
}

func main() {
	listen := flag.Bool("l", false, "Listen for incoming connections")

	// Parse the flags
	flag.Parse()

	// Set default port
	port := DefaultPort

	// Check if any arguments (port) are provided after flags
	args := flag.Args()

	protocol := flag.String("u", string(TCP), "Choose between tcp or udp")
	

	flag.Parse()

	// Validate the protocol flag
	if *protocol != string(TCP) && *protocol != string(UDP) {
		log.Fatalf("Invalid protocol: %s. Use 'tcp' or 'udp'.", *protocol)
	}
		if len(args) > 0{
			port = args[0]
		}
		
	// Start the server if the -l flag is set
	if *listen || len(flag.Args())==0{
		 // Use the first argument as the port if provided
		server := NewServer(Protocol(*protocol), port)
		server.Start()
	} else {
		// If the -l flag is not set, display the usage message
		fmt.Println("[USAGE 1]: ./TCPChat -l -p <port> -u <tcp|udp>\n[USAGE 2]: ./TCPChat $port\n[USAGE 3]: ./TCPChat")
	}
}
