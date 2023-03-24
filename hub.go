package main

import (
	"collabdraw/messages"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"strconv"
)

var upgrade = websocket.Upgrader{
	// Allow all origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
}

// newHub method creates a new Hub struct.
func newHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// run Hub method
func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

// handleWebSocket method process socket message.
func (hub *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// check client is supported for websocket.
	socket, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}
	client := newClient(hub, socket)
	hub.clients = append(hub.clients, client)
	hub.register <- client
	client.run()
	go client.read()
	go client.write()
}

// broadcast method broadcasts a message to all clients, except one(sender).
func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

// Send client data.
func (hub *Hub) send(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	client.outbound <- data
}

// onConnect method of Hub.
func (hub *Hub) onConnect(client *Client) {
	log.Println("client connected: ", client.socket.RemoteAddr())

	// Make list of all users
	var users []messages.User
	for _, c := range hub.clients {
		users = append(users, messages.User{ID: c.id, Color: c.color})
	}

	// Send exists clients to client
	hub.send(messages.NewConnected(client.color, users), client)
	// Convert Client ID from string to Int
	id, _ := strconv.Atoi(client.id)
	// Notify user joined
	hub.broadcast(messages.NewUserJoined(id, client.color), client)
}

// Handle client disconnect
func (hub *Hub) onDisconnect(client *Client) {
	log.Println("client disconnected: ", client.socket.RemoteAddr())
	client.close()

	// Find index of client
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
			break
		}
	}

	// Delete client from list
	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]

	// Convert Client ID from string to Int
	id, _ := strconv.Atoi(client.id)
	// Notify user left
	hub.broadcast(messages.NewUserLeft(id), nil)
}

// onMessage method of Hub.
func (hub *Hub) onMessage(data []byte, client *Client) {
	// Convert Client ID from string to Int
	id, _ := strconv.Atoi(client.id)
	kind := gjson.GetBytes(data, "kind").Int()
	if kind == messages.KindStroke {
		var msg messages.Stroke
		if json.Unmarshal(data, &msg) != nil {
			return
		}
		msg.UserID = id
		hub.broadcast(msg, client)
	} else if kind == messages.KindClear {
		var msg messages.Clear
		if json.Unmarshal(data, &msg) != nil {
			return
		}
		msg.UserID = id
		hub.broadcast(msg, client)
	}
}
