package main

import (
	"log"
	"net/http"
)

func adminHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf("admin")
}
