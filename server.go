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
	path    string
	gamesCh chan *GameServer
	games   map[string]*GameServer
}

// Create new chat server.
func NewServer(path string) *Server {
	return &Server{
		path,
		make(chan *GameServer),
		make(map[string]*GameServer),
	}
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

		defer ws.Close()

		id := r.URL.Query().Get("gameId")
		game := s.FindOrCreateGameServer(id, ws)
		game.Listen()
		log.Println("Websocket exit")
	}

	http.HandleFunc(s.path, onConnected)
}

func (s *Server) FindOrCreateGameServer(id string, ws *websocket.Conn) *GameServer {
	client := NewSocketClient(ws)

	if _, ok := s.games[id]; ok {
		log.Println("Reconnecting to game:", id)
	} else {
		log.Println("Creating game:", id)
		s.games[id] = NewGameServer()
	}

	s.games[id].SetClient(client)

	return s.games[id]
}
