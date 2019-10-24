package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ambelovsky/gosf"
)

func adminHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf("%s, %s", r.URL.Path, r.RemoteAddr)

	code := http.StatusBadRequest
	var body string
	if r.Method == "GET" {
		code = http.StatusOK
		body = getStringHub(hub)
		listStringHub(hub)
	} else {
		body = http.StatusText((code))
	}
	w.WriteHeader(code)
	w.Write([]byte(body))
}

func listStringHub(hub *Hub) (str string) {
	for hk, hv := range hub.rooms {
		for rk, rv := range hv.clients {
			log.Printf("room=%s,%s client=%s,%t\n", hk, hv.roomId, rk.clientId, rv)
		}
	}
	return
}

func getStringHub(hub *Hub) (str string) {
	str += fmt.Sprintf("HUB=%s, %s<br>", hub.uuid, hub.name)
	for hk, hv := range hub.rooms {
		for rk, rv := range hv.clients {
			str += fmt.Sprintf("ROOM=%s,%s CLIENT=%s,%t<br>", hk, hv.roomId, rk.clientId, rv)
		}
	}
	return
}

func handleSignalMessage(client *gosf.Client, request *gosf.Request) *gosf.Message {
	log.Printf("message: %v", request.Message)
	return gosf.NewSuccessMessage(request.Message.Text)
}
