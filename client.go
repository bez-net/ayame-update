package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	host     string
	roomId   string
	clientId string
	send     chan []byte
	sync.Mutex
}

// send json
func (c *Client) SendJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.conn.WriteJSON(v)
}

func (c *Client) Setup(roomId string, clientId string) *Client {
	c.Lock()
	defer c.Unlock()
	c.roomId = roomId
	c.clientId = clientId
	return c
}
