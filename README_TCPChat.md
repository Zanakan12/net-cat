
# TCPChat - NetCat-like Chat Application in Go

## Introduction

**TCPChat** is a network chat application written in Go that replicates basic functionality of the famous `NetCat` utility in a server-client architecture. The server listens for multiple client connections over TCP and facilitates group chat communication among clients. It supports features such as message broadcasting, client join/leave notifications, and message history sharing.

## Features

- **TCP connection**: Server listens for multiple client connections and handles communications over TCP.
- **Client naming**: Clients are required to provide a unique name upon joining.
- **Message broadcasting**: Messages sent by clients are broadcast to all other clients.
- **Chat history**: When a new client joins, they receive all previous messages in the chat.
- **Join/Leave notifications**: All clients are informed when a new client joins or an existing client leaves.
- **Timestamped messages**: Messages are timestamped and identified by the client's username.
- **Concurrency**: Utilizes Go's goroutines and synchronization mechanisms for handling multiple clients.
- **Maximum connections**: Server supports up to 10 concurrent clients.

## Requirements

- Go 1.16 or later.
- Internet connection (for distributed clients).
- Compatible OS: Windows, Linux, macOS.

## Allowed Packages

The following Go packages are used in this project:

- `io`
- `log`
- `os`
- `fmt`
- `net`
- `sync`
- `time`
- `bufio`
- `errors`
- `strings`
- `reflect`

## Project Structure

```
TCPChat/
├── client/
│   └── client.go
├── server/
│   └── server.go
├── README.md
```

- **client/client.go**: Client-side implementation to connect, send, and receive messages.
- **server/server.go**: Server-side implementation to manage multiple clients, broadcast messages, and handle chat history.

## Usage Instructions

1. **Starting the Server**:

   Navigate to the server directory and run the following command:

   ```bash
   go run server.go
   ```

   **Optional**: To start the server on a different port, specify the port as an argument:

   ```bash
   go run server.go <port>
   ```

   For example, to run the server on port `2525`, use:

   ```bash
   go run server.go 2525
   ```

2. **Connecting Clients**:

   Open a new terminal window for each client and run the client program from the client directory:

   ```bash
   go run client.go
   ```

   **Optional**: Specify a different port if the server is running on a non-default port:

   ```bash
   go run client.go <port>
   ```

3. **Using `nc`**:

   You can use `nc` (NetCat) to connect to the TCPChat server:

   ```bash
   nc <server_ip> <port>
   ```

## Features Walkthrough

- Upon connection, the server sends a Linux logo and prompts the client for a name.
- When a client sends a message, it is timestamped and broadcast to all connected clients.
- When a new client joins, they receive the chat history, and the other clients are notified of the new joiner.
- When a client leaves, the remaining clients are notified.
- Clients can send messages in real-time, and all clients receive messages sent by other participants.

## Testing

To run tests, navigate to the server directory and use the `go test` command:

```bash
cd server
go test
```

## Example Session

**Starting the Server on Port 8989**:

```bash
$ go run server.go
Listening on the port :8989
```

**Starting a Client**:

```bash
$ go run client.go
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

**Another Client Joins**:

```bash
$ go run client.go
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

- Developed by [DJIHADI RAFTANDJANI AND YANIS B..]
- Open for contributions and improvements.
