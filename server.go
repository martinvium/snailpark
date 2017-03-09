package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

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
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				log.Println("Error:", err.Error())
			}
		}()

		game := NewGameServer(ws)
		game.Listen()
	}

	http.Handle(s.path, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		case <-s.doneCh:
			return
		}
	}
}
