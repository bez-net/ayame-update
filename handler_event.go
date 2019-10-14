package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// CAUTION: don't use small case in fields of structure
type Event struct {
	Kind    string `json:"kind,omitempty"`
	OccurAt string `json:"occur_at,omitempty"`
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

	sendEventStream(hub, w, r, f)
}

func sendEventStream(hub *Hub, w http.ResponseWriter, r *http.Request, f http.Flusher) {
	edata := Event{Kind: "event", OccurAt: "2019/10/14"}
	// fmt.Println(edata)

	for i := 0; i < 100; i++ {
		edata.OccurAt = time.Now().Format("2006/01/02 15:04:05")
		jdata, err := json.Marshal(edata)
		// jdata, err := json.MarshalIndent(edata, "", "\t")
		if err != nil {
			log.Printf("json error", err)
			return
		}
		// fmt.Println(string(jdata))

		fmt.Fprintf(w, "[%2d] stream data %s\n", i+1, string(jdata))
		f.Flush()
		time.Sleep(1 * time.Second)
	}
}
