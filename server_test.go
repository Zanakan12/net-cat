package main

import (
	"net"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {
	server := NewServer("tcp","8989")
	go server.Start()
	time.Sleep(time.Second) // Allow server to start

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
}

func TestBroadcast(t *testing.T) {
	server := NewServer("tcp","8989")
	client1 := &Client{
		Username: "Alice",
		Out:      make(chan string, 10),
	}
	client2 := &Client{
		Username: "Bob",
		Out:      make(chan string, 10),
	}
	server.Clients["Alice"] = client1
	server.Clients["Bob"] = client2

	message := "[2024-04-01 12:00:00][Alice]: Hello Bob!"
	server.broadcast(message, "Alice")

	select {
	case msg := <-client2.Out:
		if msg != message+"\n" {
			t.Errorf("Expected message '%s', got '%s'", message, msg)
		}
	default:
		t.Errorf("No message received by Bob")
	}

	select {
	case msg := <-client1.Out:
		t.Errorf("Alice should not receive her own message, but got '%s'", msg)
	default:
		// Expected, no message
	}
}

func TestClientJoinLeave(t *testing.T) {
	server := NewServer("tcp","8989")

	// Simulate client join
	client := &Client{
		Username: "Charlie",
		Out:      make(chan string, 10),
	}
	server.Clients["Charlie"] = client
	server.broadcast("[INFO]: Charlie has joined the chat.\n", "INFO")

	select {
	case msg := <-client.Out:
		// New client should not receive join message
		t.Errorf("New client should not receive join message, but got '%s'", msg)
	default:
		// Expected, no message
	}

	// Add another client to receive the join message
	client2 := &Client{
		Username: "Dave",
		Out:      make(chan string, 10),
	}
	server.Clients["Dave"] = client2
	server.broadcast("[INFO]: Dave has joined the chat.\n", "INFO")

	select {
	case msg := <-client.Out:
		if msg != "[INFO]: Dave has joined the chat.\n" {
			t.Errorf("Expected join message, got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Errorf("Dave's join message not received by Charlie")
	}

	// Simulate client leave
	delete(server.Clients, "Charlie")
	server.broadcast("[INFO]: Charlie has left the chat.\n", "INFO")

	select {
	case msg := <-client2.Out:
		if msg != "[INFO]: Charlie has left the chat.\n" {
			t.Errorf("Expected leave message, got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Errorf("Leave message not received by Dave")
	}
}
