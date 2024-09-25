
# TCPChat - NetCat-like Chat Application in Go

## Introduction

**TCPChat** is a Go-based network chat server that replicates basic functionality of `NetCat` for group chat communication. This application uses a server-client architecture and supports multiple clients connected over TCP. Clients can send and receive messages in real-time, change their usernames, and exit the chat.

## Features

- **Multiple Clients Support**: Supports up to 10 concurrent clients.
- **TCP & UDP Support**: The server can be started in either TCP or UDP mode.
- **Client Naming**: Clients must provide a unique username when joining the server.
- **Message Broadcasting**: Messages sent by clients are broadcast to all connected clients.
- **Chat History**: New clients receive all previous messages when they join the chat.
- **Name Change**: Clients can change their username using `/name <newname>`.
- **Join/Leave Notifications**: All clients are notified when someone joins or leaves.
- **Concurrency**: Utilizes Go’s goroutines and synchronization mechanisms to handle multiple clients concurrently.
- **Graceful Shutdown**: Server resources are cleaned up upon shutdown.

## Requirements

- Go 1.16 or later.
- Internet connection for distributed clients.
- Supported OS: Windows, Linux, macOS.

## Project Structure

```
.
├── main.go          # Main server code
├── server_test.go   # Test code for TCP and UDP servers
├── README.md        # This README file
```

## Installation and Setup

### 1. Build the Project

To build the project, run the following command:

```bash
go build -o TCPchat main.go
```

### 2. Starting the Server

To run the TCP or UDP chat server, use one of the following commands:

#### TCP Mode
```bash
./TCPchat -l -u tcp
```

#### UDP Mode
```bash
./TCPchat -l -u udp
```

You can specify a different port by passing it as an argument:
```bash
./TCPchat -l -u tcp 9000
```

### 3. Connecting Clients

Clients can connect using `telnet` or `netcat`:

#### Using Telnet
```bash
telnet localhost 8989
```

#### Using NetCat
```bash
nc localhost 8989
```

## Commands

### Changing Username

A client can change their username by sending the following command:
```
/name <newname>
```

### Exiting the Chat

A client can exit the chat by sending:
```
/exit
```

## Testing

Run tests to ensure the server functions correctly by executing:

```bash
go test
```

This will run the test cases defined in the `server_test.go` file.

## Example Session

### Starting the Server on Port 8989:
```bash
$ ./TCPchat -l -u tcp
Listening on port :8989
```

### Starting a Client:
```bash
$ telnet localhost 8989
          .--.
         |o_o |
         |:_/ |
        //   \ \
       (|     | )
      /'\_   _/`\
      \___)=(___/
Enter your name: Alice
[INFO]: Alice has joined the chat.
```

### Another Client Joins:
```bash
$ telnet localhost 8989
          .--.
         |o_o |
         |:_/ |
        //   \ \
       (|     | )
      /'\_   _/`\
      \___)=(___/
Enter your name: Bob
[INFO]: Bob has joined the chat.
[Alice]: Hello everyone!
```

## License

This project is open-sourced under the MIT License.

## Authors

- Developed by [DJIHADI RAFTANDJANI AND YANIS Bellahouel..]
- Open for contributions and improvements.
