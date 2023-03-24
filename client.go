package main

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Client struct {
	hub      *Hub
	id       string
	color    string
	socket   *websocket.Conn
	outbound chan []byte
	send     chan []byte
}

// Constructor function for the client that assigns UUID and a random color
func newClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
		color:    generateColor(),
		id:       uuid.NewV4().String(),
	}
}

// Reads messages sent from the client and forwards them to the hub
func (client *Client) read() {
	defer func() {
		client.hub.unregister <- client
	}()
	for {
		_, data, err := client.socket.ReadMessage()
		if err != nil {
			break
		}
		client.hub.onMessage(data, client)
	}
}

// Write takes messages from the outbound channel and sends them to the client
func (client *Client) write() {
	for {
		select {
		case data, ok := <-client.outbound:
			if !ok {
				err := client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			err := client.socket.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}
		}
	}
}

// Close channel connection to client
func (client *Client) close() {
	err := client.socket.Close()
	if err != nil {
		return
	}
	close(client.outbound)
}

// Client run
func (client Client) run() {
	go client.read()
	go client.write()
}
