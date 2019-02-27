package main

import (
	"github.com/gorilla/websocket"
	"time"
)

type Node struct {
	hub  *Hub
	id   string
	conn *websocket.Conn
}

func newNode(hub *Hub, id string, conn *websocket.Conn) *Node {
	return &Node{
		hub:  hub,
		id:   id,
		conn: conn,
	}
}

func (node *Node) Run() {
	expire := time.NewTimer(2 * time.Minute)
	go func() {
		<-expire.C
		node.hub.unregister <- node.id
	}()
	for {
		event := make(map[string]string)
		err := node.conn.ReadJSON(&event)
		if err != nil {
			expire.Stop()
			node.hub.unregister <- node.id
			return
		}
		event["source"] = node.id
		node.hub.api <- event
	}
}
