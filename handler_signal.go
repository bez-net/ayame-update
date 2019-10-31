package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ryanuber/go-glob"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// allow cross orgin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SignalMessage struct {
	Type     string       `json:"type"`
	RoomId   string       `json:"roomId"`
	ClientId string       `json:"clientId"`
	Metadata *interface{} `json:"authn_metadata,omitempty"`
	Key      *string      `json:"key,omitempty"`
}

type PingMessage struct {
	Type string `json:"type"`
}

func (c *Client) listen(cancel context.CancelFunc) {
	defer func() {
		cancel()
		c.hub.unregister <- &RegisterInfo{
			client: c,
			roomId: c.roomId,
		}
		c.conn.Close()
	}()

	upgrader.CheckOrigin = func(r *http.Request) bool {
		if Options.AllowOrigin == "" {
			return true
		}
		origin := r.Header.Get("Origin")
		// trim origin
		host, err := TrimOriginToHost(origin)
		if err != nil {
			log.Println("Invalid Origin Header, header=", origin)
		}
		// check the origin is same with one of Allow Origin in config.yaml
		log.Printf("[WS] Request Origin=%s, AllowOrigin=%s", origin, Options.AllowOrigin)
		if &Options.AllowOrigin == host {
			return true
		}
		if glob.Glob(Options.AllowOrigin, *host) {
			return true
		}
		return false
	}

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	log.Printf("ws error: %v", err)
			// }
			log.Printf("%v", err)
			return
		}

		msg := &SignalMessage{}
		json.Unmarshal(message, &msg)
		log.Printf("message=%s", message)

		switch msg.Type {
		case "":
			log.Printf("ignore null message")
			continue
		case "pong":
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
		case "register":
			if msg.ClientId == "" || msg.RoomId == "" {
				log.Printf("%s error: clientId=%s, roomId=%s", msg.Type, msg.ClientId, msg.RoomId)
				return
			}
			c.hub.register <- &RegisterInfo{
				clientId: msg.ClientId,
				client:   c,
				roomId:   msg.RoomId,
				key:      msg.Key,
				metadata: msg.Metadata,
			}
		// case "onmessage":
		default:
			if c.clientId == "" || c.roomId == "" {
				log.Printf("%s error: client not registered: %v", msg.Type, c)
				return
			}
		}

		// Broadcast the signaling message received
		broadcast := &Broadcast{
			client:   c,
			roomId:   c.roomId,
			messages: message,
		}
		c.hub.broadcast <- broadcast
	}
}

func (c *Client) broadcast(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			// exit the loop if the channel already close
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("%v", err)
				return
			}
			w.Write(message)
			w.Close()
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// if over_ws_ping_pong option is set
			if Options.OverWsPingPong {
				log.Println("send ping over WS")
				pingMsg := &PingMessage{Type: "ping"}
				if err := c.SendJSON(pingMsg); err != nil {
					log.Printf("%v", err)
					return
				}
			} else {
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("%v", err)
					return
				}
			}
		}
	}
}

func signalHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf("%s, %s", r.URL.Path, r.RemoteAddr)
	defer log.Printf("signalHandler exit")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	origin := r.Header.Get("Origin")
	host, err := TrimOriginToHost(origin)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: c, host: *host, send: make(chan []byte, 256)}
	log.Printf("[WS] connected")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go client.listen(cancel)
	go client.broadcast(ctx)
}
