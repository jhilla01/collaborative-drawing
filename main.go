package main

import (
	"log"
	"net/http"
)

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", hub.handleWebSocket)
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
	}
}
