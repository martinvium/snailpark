package main

import (
	"github.com/gorilla/websocket"
	"io"
	"log"
)

type Client interface {
	Listen(chan *Message)
	SendResponse(msg *ResponseMessage)
	PlayerId() string
	Done()
}

// # BaseClient definition

type BaseClient struct {
	playerId string
	msgCh    chan *ResponseMessage
	doneCh   chan bool
}

func (c *BaseClient) SendResponse(msg *ResponseMessage) {
	c.msgCh <- msg
}

func (c *BaseClient) PlayerId() string {
	return c.playerId
}

func (c *BaseClient) String() string {
	return "Client(" + c.PlayerId() + ")"
}

func (c *BaseClient) Done() {
	log.Println("Client done")
	c.doneCh <- true
}

// # AIClient definition

type AIClient struct {
	BaseClient
	ai *AI
}

func NewAIClient(ai *AI) *AIClient {
	return &AIClient{
		BaseClient{
			"ai",
			make(chan *ResponseMessage, channelBufSize),
			make(chan bool),
		},
		ai,
	}
}

func (c *AIClient) Listen(requestCh chan *Message) {
	go c.listenWrite()
	go c.listenRead(requestCh)
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
			log.Println("Done received")
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Receive stuff from the AI over channel
func (c *AIClient) listenRead(requestCh chan *Message) {
	log.Println("Listening read from AI")

	for {
		select {

		// read data from websocket connection
		case msg := <-c.ai.outCh:
			if c.PlayerId() == msg.PlayerId {
				requestCh <- msg
			} else {
				log.Println("Error: Wrong client id: " + c.PlayerId() + " != " + msg.PlayerId)
			}

		// receive done request
		case <-c.doneCh:
			log.Println("Done received")
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

func NewSocketClient(ws *websocket.Conn) *SocketClient {
	return &SocketClient{
		BaseClient{
			"player",
			make(chan *ResponseMessage, channelBufSize),
			make(chan bool),
		},
		ws,
	}
}

// Listen Write and Read request via chanel
func (c *SocketClient) Listen(requestCh chan *Message) {
	go c.listenWrite()
	c.listenRead(requestCh)
}

// Send stuff to the client over socket
func (c *SocketClient) listenWrite() {
	log.Println("Listening write to client")

	for {
		select {

		// send message to the client
		case msg := <-c.msgCh:
			log.Println("Send:", msg)
			if err := c.ws.WriteJSON(msg); err != nil {
				log.Println("ERROR: WriteJSON failed:", err)
			}

		// receive done request
		case <-c.doneCh:
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Receive stuff from the client over socket
func (c *SocketClient) listenRead(requestCh chan *Message) {
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
			err := c.ws.ReadJSON(&msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				log.Println("Error:", err.Error())
				c.doneCh <- true
			} else {
				if c.PlayerId() == msg.PlayerId {
					requestCh <- &msg
				} else {
					log.Println("Error: Wrong client id: " + c.PlayerId() + " != " + msg.PlayerId)
				}
			}
		}
	}
}
