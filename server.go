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
	errCh  chan error
}

// Create new chat server.
func NewServer(path string) *Server {
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		path,
		doneCh,
		errCh,
	}
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

// Listen and serve.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		game := NewGameServer(ws, s)
		game.Listen()
	}

	http.Handle(s.path, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Probably really not worth it to have a channel dedicated to this, just
		// dry up the error handling below?
		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
