package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
)

const channelBufSize = 100

type GameServer struct {
	ws       *websocket.Conn
	server   *Server
	messages []*Message
	msgCh    chan *Message
	doneCh   chan bool
}

func NewGameServer(ws *websocket.Conn, server *Server) *GameServer {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	messages := []*Message{}
	msgCh := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)

	return &GameServer{
		ws,
		server,
		messages,
		msgCh,
		doneCh,
	}
}

// Listen Write and Read request via chanel
func (g *GameServer) Listen() {
	go g.listenWrite()
	g.listenRead()
}

// Send stuff to the client over socket
func (g *GameServer) listenWrite() {
	log.Println("Listening write to client")

	for {
		select {

		// send message to the client
		case msg := <-g.msgCh:
			log.Println("Send:", msg)
			websocket.JSON.Send(g.ws, msg)

		// receive done request
		case <-g.doneCh:
			g.doneCh <- true // for listenRead method
			return
		}
	}
}

// Receive stuff from the client over socket
func (g *GameServer) listenRead() {
	log.Println("Listening read from client")

	for {
		select {

		// receive done request
		case <-g.doneCh:
			g.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(g.ws, &msg)
			if err == io.EOF {
				g.doneCh <- true
			} else if err != nil {
				g.server.Err(err)
			} else {
				g.handleAction(&msg)
			}
		}
	}
}

func (g *GameServer) handleAction(msg *Message) {
	log.Println("Handleling action from client!")
}
