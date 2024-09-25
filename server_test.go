package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

// TestTCPServer tests the TCP chat server's basic functionality.
func TestTCPServer(t *testing.T) {
	// Start the server in a separate goroutine
	server := NewServer(TCP, "9000")
	go server.Start()

	// Allow time for the server to start
	time.Sleep(1 * time.Second)

	// Connect as a client
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Read the initial Linux logo and prompt
	scanner := bufio.NewScanner(conn)
	scanner.Scan() // Linux logo
	scanner.Scan() // Enter your name:

	// Send the username
	conn.Write([]byte("TestUser\n"))

	// Test receiving broadcast message
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay before broadcasting
		server.broadcast("[INFO]: Test broadcast message\n", "INFO")
	}()

	// Timeout after 5 seconds if no broadcast received
	timeout := time.NewTimer(5 * time.Second)

	// Channel to signal received message
	done := make(chan bool)

	// Reading message from the server
	go func() {
		for scanner.Scan() {
			msg := scanner.Text()
			if msg == "[INFO]: Test broadcast message" {
				t.Log("Received broadcast message")
				done <- true
				return
			}
		}
	}()

	// Select for either a message or timeout
	select {
	case <-done:
		// Successfully received the message
	case <-timeout.C:
		t.Fatal("Test timed out waiting for broadcast message")
	}

	// Shutdown the server gracefully
	server.Shutdown()
}

// TestUDPServer tests the UDP message receipt functionality.
func TestUDPServer(t *testing.T) {
	// Start the server in a separate goroutine
	server := NewServer(UDP, "9001")
	go server.startUDP()

	// Allow time for the server to start
	time.Sleep(1 * time.Second)

	// Send a UDP message
	conn, err := net.Dial("udp", "localhost:9001")
	if err != nil {
		t.Fatalf("Failed to connect to UDP server: %v", err)
	}
	defer conn.Close()

	message := "Hello from UDP client"
	conn.Write([]byte(message))

	// Allow time for the server to process the message
	time.Sleep(1 * time.Second)

	// Shutdown the server after test
	server.Shutdown()
}