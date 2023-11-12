package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub *Hub
	id string
	socket *websocket.Conn
	outBound chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub: hub,
		socket: socket,
		outBound: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		select {
		case msg,ok := <-c.outBound:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (c *Client) Close() {
	// c.hub.unregister <- c
	c.socket.Close()
	close(c.outBound)
}

// func (c *Client) Read() {
// 	defer func() {
// 		c.hub.unregister <- c
// 		c.socket.Close()
// 	}()

// 	for {
// 		_, msg, err := c.socket.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		c.hub.broadcast <- msg
// 	}
// }