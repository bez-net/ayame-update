package main

type Broadcast struct {
	client   *Client
	roomId   string
	messages []byte
}

type RegisterInfo struct {
	roomId   string
	clientId string
	client   *Client
	metadata *interface{}
	key      *string
}

type EventInfo struct {
	class   string
	content string
}

// TODO: Think more for the need
const (
	ROOM_TYPE_PRE_DEFINED  = iota
	ROOM_TYPE_USER_DEFINED = iota
)

type Room struct {
	Common
	clients map[*Client]bool
	roomId  string
}

type Hub struct {
	Common
	rooms      map[string]*Room
	broadcast  chan *Broadcast
	register   chan *RegisterInfo
	unregister chan *RegisterInfo
	event      chan *EventInfo
}

func newHub(name string) (hub *Hub) {
	hub = &Hub{
		broadcast:  make(chan *Broadcast),
		register:   make(chan *RegisterInfo),
		unregister: make(chan *RegisterInfo),
		event:      make(chan *EventInfo),
		rooms:      make(map[string]*Room),
	}
	hub.name = name
	hub.uuid = getStringUUID()
	return
}

func (h *Hub) run() {
	for {
		select {
		case registerInfo := <-h.register:
			client := registerInfo.client
			clientId := registerInfo.clientId
			roomId := registerInfo.roomId
			client = client.Setup(roomId, clientId)
			room := h.rooms[roomId]
			if _, ok := h.rooms[roomId]; !ok {
				room = &Room{
					clients: make(map[*Client]bool),
					roomId:  roomId,
				}
				h.rooms[roomId] = room
			}
			ok := len(room.clients) < Options.MaxSessions
			if !ok {
				msg := &RejectMessage{
					Type:   "reject",
					Reason: "TOO-MANY-USERS",
				}
				client.SendJSON(msg)
				client.conn.Close()
				break
			}
			isExistUser := len(room.clients) > 0
			// use auth webhook
			if Options.AuthWebhookURL != "" {
				resp, err := AuthWebhookRequest(registerInfo.key, roomId, registerInfo.metadata, client.host)
				if err != nil {
					msg := &RejectMessage{
						Type:   "reject",
						Reason: "AUTH-WEBHOOK-ERROR",
					}
					if resp != nil {
						msg.Reason = resp.Reason
					}
					client.SendJSON(msg)
					client.conn.Close()
					break
				}

				msg := &AcceptMetadataMessage{
					Type:        "accept",
					IceServers:  resp.IceServers,
					IsExistUser: isExistUser,
				}
				if resp.AuthzMetadata != nil {
					msg.Metadata = resp.AuthzMetadata
				}

				room.clients[client] = true
				client.SendJSON(msg)
			} else {
				msg := &AcceptMessage{
					Type:        "accept",
					IsExistUser: isExistUser,
				}
				room.clients[client] = true
				client.SendJSON(msg)
			}
		case registerInfo := <-h.unregister:
			roomId := registerInfo.roomId
			client := registerInfo.client
			if room, ok := h.rooms[roomId]; ok {
				if _, ok := room.clients[client]; ok {
					delete(room.clients, client)
					close(client.send)
				}
			}
		case broadcast := <-h.broadcast:
			if room, ok := h.rooms[broadcast.roomId]; ok {
				for client := range room.clients {
					if client.clientId != broadcast.client.clientId {
						select {
						case client.send <- broadcast.messages:
						default:
							close(client.send)
							delete(room.clients, client)
						}
					}
				}
			}
		}
	}
}
