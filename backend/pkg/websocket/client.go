package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"

	"backend/pkg/message"
)

type Client struct {
	ID   int32
	Name string
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message from client - %v", err)
			return
		}
		c.Pool.Broadcast <- message.Message{
			ClientId:   c.ID,
			ClientName: c.Name,
			Time:       time.Now(),
			Content:    string(p),
		}
	}
}
