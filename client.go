package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
)

type Client interface {
	Listen(g *GameServer)
	SendResponse(msg *Message)
	ClientId() string
}

// # BaseClient definition

type BaseClient struct {
	clientId string
	msgCh    chan *Message
	doneCh   chan bool
}

func (c *BaseClient) SendResponse(msg *Message) {
	c.msgCh <- msg
}

func (c *BaseClient) ClientId() string {
	return c.clientId
}

func (c *BaseClient) String() string {
	return "Client(" + c.ClientId() + ")"
}

// # AIClient definition

type AIClient struct {
	BaseClient
	ai *AI
}

func (c *AIClient) Listen(g *GameServer) {
	go c.listenWrite()
	go c.listenRead(g)
}

// Send stuff to the AI over channel
func (c *AIClient) listenWrite() {
	log.Println("Listening write to AI")

	for {
		select {

		// send message to the client
		case msg := <-c.msgCh:
			log.Println("Send:", msg)
			c.ai.Send(msg)

		// receive done request
		case <-c.doneCh:
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Receive stuff from the AI over channel
func (c *AIClient) listenRead(g *GameServer) {
	log.Println("Listening read from AI")

	for {
		select {

		// read data from websocket connection
		case msg := <-c.ai.msgCh:
			if c.ClientId() == msg.ClientId {
				g.SendRequest(msg)
			} else {
				log.Println("Error: Wrong client id: " + c.ClientId() + " != " + msg.ClientId)
			}

		// receive done request
		case <-c.doneCh:
			c.doneCh <- true // for listenWrite method
			return

		}
	}
}

// # SocketClient definition

type SocketClient struct {
	BaseClient
	ws *websocket.Conn
}

// Listen Write and Read request via chanel
func (c *SocketClient) Listen(g *GameServer) {
	go c.listenWrite()
	c.listenRead(g)
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
func (c *SocketClient) listenRead(g *GameServer) {
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
				if c.ClientId() == msg.ClientId {
					g.SendRequest(&msg)
				} else {
					log.Println("Error: Wrong client id: " + c.ClientId() + " != " + msg.ClientId)
				}
			}
		}
	}
}
