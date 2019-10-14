package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// CAUTION: don't use small case in fields of structure
type EventData struct {
	UserId  string `json:"user_id,omitempty"`
	OccurAt string `json:"occur_at,omitempty"`
}

type EventMessage struct {
	Event string
	Id    string
	Retry string
	Data  EventData
}

func eventHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf(r.URL.Path)
	// check if SSE is supported
	f, ok := w.(http.Flusher)
	if !ok {
		log.Printf("SSE Streaming unsuported")
		return
	}

	// flusher required headers for SSE streams
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // cojam.tv

	sendEventStream(hub, w, r, f)
}

func sendEventStream(hub *Hub, w http.ResponseWriter, r *http.Request, f http.Flusher) {
	emsg := EventMessage{Event: "notify", Id: "ayame", Retry: "2"}
	emsg.Data.UserId = "sikang99@gmail.com"
	// fmt.Println(emsg)

	for i := 0; i < 100; i++ {
		str := genStringEventMessage(emsg)
		fmt.Fprintf(w, str)
		f.Flush()
		time.Sleep(1 * time.Second)
	}
	log.Printf("closed")
}

func genStringEventMessage(emsg EventMessage) (str string) {
	emsg.Data.OccurAt = time.Now().Format("2006/01/02 15:04:05")
	jdata, err := json.Marshal(emsg.Data)
	// jdata, err := json.MarshalIndent(edata, "", " ")
	if err != nil {
		log.Printf("json error: ", err)
		return
	}
	str = fmt.Sprintf("event:%s\nretry:%s\nid:%s\ndata:%s\n\n", emsg.Event, emsg.Retry, emsg.Id, string(jdata))
	return
}
