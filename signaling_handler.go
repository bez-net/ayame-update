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

type Message struct {
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
		// origin を trim
		host, err := TrimOriginToHost(origin)
		if err != nil {
			log.Println("Invalid Origin Header, header=", origin)
		}
		// config.yaml で指定した Allow Origin と一致するかで検査する
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
		msg := &Message{}
		json.Unmarshal(message, &msg)
		log.Printf("signaling: %s %s", msg.Type, message)

		switch msg.Type {
		case "":
			// log.Println("invalid signaling type: ", msg.Type)
			break
		case "pong":
			// log.Println("recv ping over WS")
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
			break
		case "register":
			if msg.RoomId == "" {
				log.Printf(msg.Type, "invalid room id=", msg.RoomId)
				return
			}
			log.Printf("%s: %v", msg.Type, msg)
			c.hub.register <- &RegisterInfo{
				clientId: msg.ClientId,
				client:   c,
				roomId:   msg.RoomId,
				key:      msg.Key,
				metadata: msg.Metadata,
			}
			break
		case "onmessage":
			if c.roomId == "" {
				log.Printf("client does not registered: %v", c)
				return
			}
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				return
			}
			break
		default:
			// log.Println("pass signaling type: ", msg.Type)
			break
		}

		log.Printf("clientId=%s, roomId=%s", c.clientId, c.roomId)

		// Broadcast the signaling message
		broadcast := &Broadcast{
			client:   c,
			roomId:   c.roomId,
			messages: message,
		}
		c.hub.broadcast <- broadcast

		/*
			if msg.Type == "" {
				log.Println("Invalid Signaling Type")
				break
			}
			if msg.Type == "pong" {
				log.Println("recv ping over WS")
				c.conn.SetReadDeadline(time.Now().Add(pongWait))
			} else {
				if msg.Type == "register" && msg.RoomId != "" {
					log.Printf("register: %v", msg)
					c.hub.register <- &RegisterInfo{
						clientId: msg.ClientId,
						client:   c,
						roomId:   msg.RoomId,
						key:      msg.Key,
						metadata: msg.Metadata,
					}
				} else {
					log.Printf("onmessage: %s", message)
					log.Printf("client roomId: %s", c.roomId)
					if c.roomId == "" {
						log.Printf("client does not registered: %v", c)
						return
					}
					if err != nil {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
							log.Printf("error: %v", err)
						}
						break
					}
					broadcast := &Broadcast{
						client:   c,
						roomId:   c.roomId,
						messages: message,
					}
					c.hub.broadcast <- broadcast
				}
			}
		*/
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
			// exit loop if channel already close
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
				return
			}
			w.Write(message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// if over_ws_ping_pong is set
			if Options.OverWsPingPong {
				log.Println("send ping over WS")
				pingMsg := &PingMessage{Type: "ping"}
				if err := c.SendJSON(pingMsg); err != nil {
					return
				}
			} else {
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

func signalingHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
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
