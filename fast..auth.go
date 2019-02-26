package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hub = newHub()

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub.register <- conn
}

func main() {
	go hub.Run()
	http.HandleFunc("/api/fast..auth", handler)
	log.Println("fast..auth working on :8003")
	log.Fatal(http.ListenAndServe("127.0.0.1:8003", nil))
}
