package main

import (
	"fmt"
	"log"
	"net/http"
)

func uploadHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf(r.URL.Path)
	fmt.Fprintf(w, "%s not implemented", r.URL.Path)
}
