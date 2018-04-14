package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func checkOrigin(r *http.Request) bool {
	// TODO actually check origin...
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// Chat server.
type Server struct {
	gamesCh chan *GameServer
	games   map[string]*GameServer
	ticker  *time.Ticker
}

// Create new chat server.
func NewServer() *Server {
	return &Server{
		make(chan *GameServer),
		make(map[string]*GameServer),
		time.NewTicker(60 * time.Second),
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

	http.HandleFunc("/game/connect", onConnected)

	http.HandleFunc("/game/stats", s.handleStats)

	go s.gameCleanup()
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(s.games))
	for k := range s.games {
		keys = append(keys, k)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(keys); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func (s *Server) gameCleanup() {
	killEm := make(chan string, channelBufSize)

	for {
		select {
		case id := <-killEm:
			log.Println("Killing:", id)
			go s.games[id].Done()
			delete(s.games, id)
		case <-s.ticker.C:
			for id, game := range s.games {
				if time.Since(game.StartTime).Minutes() > 15 {
					killEm <- id
				}
			}
		}
	}
}
