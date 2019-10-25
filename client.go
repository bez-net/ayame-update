package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Common part of structs such as client, room, hub
type Common struct {
	uuid    string
	name    string
	class   string
	created string
}

type Client struct {
	Common
	device   string // device unique string
	hub      *Hub
	conn     *websocket.Conn
	host     string
	roomId   string
	clientId string
	send     chan []byte
	sync.Mutex
}

// send json data
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

type User struct {
	Common
	nick string
	mail string
	SNSs map[string]string
}
type Group struct {
	Common
	users map[string]*User
}
