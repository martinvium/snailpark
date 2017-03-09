package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
)

type Client interface {
	Listen(gs *GameServer)
	SendMessage(msg *Message)
}

type AIClient struct {
	msgCh chan *Message
}

type SocketClient struct {
	ws     *websocket.Conn
	msgCh  chan *Message
	doneCh chan bool
}

func (c *AIClient) Listen(gs *GameServer) {
	// TODO implement
}

func (c *AIClient) SendMessage(msg *Message) {
	c.msgCh <- msg
}

func (c *SocketClient) SendMessage(msg *Message) {
	c.msgCh <- msg
}

// Listen Write and Read request via chanel
func (c *SocketClient) Listen(gs *GameServer) {
	go c.listenWrite()
	c.listenRead(gs)
}

// Send stuff to the client over socket
func (c *SocketClient) listenWrite() {
	log.Println("Listening write to client")

	for {
		select {

		// send message to the client
		case msg := <-c.msgCh:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Receive stuff from the client over socket
func (c *SocketClient) listenRead(gs *GameServer) {
	log.Println("Listening read from client")

	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				log.Println("Error:", err.Error())
			} else {
				gs.handleAction(&msg)
			}
		}
	}
}
