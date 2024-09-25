
# TCP Chat Server in Go

This project implements a simple TCP chat server in Go that supports multiple clients. It allows users to connect via TCP or UDP, send messages, and change their username using commands.

## Features
- **Multiple Clients Support**: Handles up to 10 clients concurrently.
- **TCP & UDP Support**: The server can be started in either TCP or UDP mode.
- **Message Broadcasting**: Broadcasts messages to all connected clients.
- **Name Change Command**: Clients can change their username using the `/name <newname>` command.
- **Graceful Shutdown**: Server shuts down gracefully, ensuring all resources are cleaned up.
- **Linux Logo**: Displays a cool Linux logo to clients upon connection.

## Usage

### Build the Project
You can build the project using the following command:
```bash
go build -o TCPchat main.go
```

### Run the Server
To run the TCP or UDP chat server, use the following commands:

#### TCP Mode
```bash
./TCPchat -l -u tcp
```

#### UDP Mode
```bash
./TCPchat -l -u udp
```

By default, the server listens on port 8989. You can specify a different port by passing it as an argument:
```bash
./TCPchat -l -u tcp 9000
```

### Connecting as a Client
You can use `telnet` or `netcat` to connect to the server:
```bash
telnet localhost 8989
```

Or using `netcat`:
```bash
nc localhost 8989
```

## Commands

### Changing Username
A user can change their username by sending the following command:
```
/name <newname>
```

### Exiting the Chat
A user can exit the chat by sending the following command:
```
/exit
```

## Project Structure

```
.
├── main.go          # The main server code
├── server_test.go   # The test code for TCP and UDP servers
└── README.md        # This README file
```

## Testing

The project contains unit tests for the TCP and UDP functionality. You can run the tests using:
```bash
go test
```

This will run the test cases defined in the `server_test.go` file to ensure the server functions correctly.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

