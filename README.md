# Chat App

This is a simple chat application built with Golang and WebSocket.

## Setup

1. Install [Golang](https://golang.org/dl/).
2. Clone the repository.
3. Navigate to the project directory.
4. Run `go mod tidy` to download dependencies.
5. Run `go run main.go` to start the server.
6. Open `client/index.html` in your browser.

## Usage

- Open multiple browser tabs pointing to `client/index.html`.
- Type a message in the input field and click "Send".
- Messages will be broadcasted to all connected clients.
