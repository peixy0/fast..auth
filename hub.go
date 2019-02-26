package main

import (
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

func newId(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type Hub struct {
	register   chan *websocket.Conn
	unregister chan string
	api        chan map[string]string
	nodes      map[string]*Node
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *websocket.Conn),
		unregister: make(chan string),
		api:        make(chan map[string]string),
		nodes:      make(map[string]*Node),
	}
}

func (hub *Hub) Run() {
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		select {
		case conn := <-hub.register:
			id := newId(16)
			for hub.nodes[id] != nil {
				id = newId(16)
			}
			node := newNode(hub, id, conn)
			hub.nodes[id] = node
			go node.Run()
		case id := <-hub.unregister:
			node := hub.nodes[id]
			if node != nil {
				node.conn.Close()
				delete(hub.nodes, id)
			}
		case event := <-hub.api:
			if event["event"] == "Auth.Id" {
				hub.idHandler(event["source"])
			}

			if event["event"] == "Auth.Token" {
				hub.tokenHandler(event["source"], event["target"], event["token"])
			}
		}
	}
}

func (hub *Hub) idHandler(id string) {
	node := hub.nodes[id]
	if node == nil {
		return
	}
	result := make(map[string]string)
	result["event"] = "Auth.Id"
	result["id"] = id
	node.conn.WriteJSON(result)
}

func (hub *Hub) tokenHandler(sourceId string, targetId string, token string) {
	source := hub.nodes[sourceId]
	if source == nil {
		return
	}
	target := hub.nodes[targetId]
	if target == nil {
		return
	}
	result := make(map[string]string)
	result["event"] = "Auth.Token"
	result["token"] = token
	target.conn.WriteJSON(result)
	go func() {
		hub.unregister <- sourceId
		hub.unregister <- targetId
	}()
}
