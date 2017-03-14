package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Chat server.
type Server struct {
	path   string
	doneCh chan bool
}

// Create new chat server.
func NewServer(path string) *Server {
	doneCh := make(chan bool)

	return &Server{
		path,
		doneCh,
	}
}

func (s *Server) Done() {
	s.doneCh <- true
}

// Listen and serve.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		game := NewGameServer(ws)
		game.Listen()
	}

	http.HandleFunc(s.path, onConnected)

	for {
		select {

		case <-s.doneCh:
			return
		}
	}
}
