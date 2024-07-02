package server

import (
	"bytes"
	"log"
	"time"
	"github.com/gorilla/websocket"
)

const (
	wirteWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var newline = []byte{'\n'}
var space = []byte{' '}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(
		func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod) //tạo 1 Ticker để gửi tin nhắn ping định kì
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select { // select la goroutine, cho tin nhan tu 'c.send' hoac 'ticker.C'
		case message, ok := <-c.send: // cho tin nhan tu channel c.send
			c.conn.SetWriteDeadline(time.Now().Add(wirteWait)) // dat thoi han ghi (write deadline)
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage) // tao 1 writer de ghi tin nhan van ban (text message)
			if err != nil {
				return
			}
			w.Write(message)

			// Ghi các tin nhắn còn lại trong c.send
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(wirteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
