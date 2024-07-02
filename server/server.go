package server

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type Hub struct {
	// Registered connectios
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }

		case message := <-h.broadcast:
			for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
		}
	}
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,//xác định kích thước bộ đệm đọc (read buffer size) 
						//là 1024 byte. Bộ đệm này được sử dụng để đọc dữ liệu từ kết nối WebSocket. 
						//Kích thước này xác định lượng dữ liệu tối đa có thể được đọc từ kết nối WebSocket trong một lần đọc.
    WriteBufferSize: 1024, //xác định kích thước bộ đệm ghi (write buffer size) là 1024 byte. 
						//Bộ đệm này được sử dụng để ghi dữ liệu vào kết nối WebSocket. 
						//Kích thước này xác định lượng dữ liệu tối đa có thể được ghi vào kết nối WebSocket trong một lầ
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}