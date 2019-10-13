package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ambelovsky/gosf"
)

func adminHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf(r.URL.Path)
	code := http.StatusBadRequest
	var body string
	if r.Method == "GET" {
		code = http.StatusOK
		body = strHub(hub)
	} else {
		body = http.StatusText((code))
	}
	w.WriteHeader(code)
	w.Write([]byte(body))

	listHub(hub)
}

func listHub(hub *Hub) (str string) {
	for hk, hv := range hub.rooms {
		// log.Printf("%s %s", hk, hv.roomId)
		for rk, rv := range hv.clients {
			log.Printf("room=%s,%s client=%s,%t", hk, hv.roomId, rk.clientId, rv)
		}
	}
	return
}

func strHub(hub *Hub) (str string) {
	str += fmt.Sprintf("hub=%s\n", hub.ID)
	for hk, hv := range hub.rooms {
		for rk, rv := range hv.clients {
			str += fmt.Sprintf("room=%s,%s client=%s,%t\n", hk, hv.roomId, rk.clientId, rv)
		}
	}
	return
}

func handleSignalMessage(client *gosf.Client, request *gosf.Request) *gosf.Message {
	log.Printf("message: %v", request.Message)
	return gosf.NewSuccessMessage(request.Message.Text)
}
