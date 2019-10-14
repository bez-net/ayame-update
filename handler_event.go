package main

import (
	"log"
	"net/http"
)

func eventHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf(r.URL.Path)
}
