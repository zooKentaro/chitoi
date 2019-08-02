package model

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

var maxID uint64 = 1001

type Client struct {
	ID     uint64
	ws     *websocket.Conn
	server *Server
	ch     chan string
	doneCh chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	maxID++
	var (
		id     = maxID
		doneCh = make(chan bool)
		ch     = make(chan string)
	)

	return &Client{
		ID:     id,
		ws:     ws,
		server: server,
		ch:     ch,
		doneCh: doneCh,
	}
}

func (c *Client) Listen() {
	go c.listenWrite()

	log.Println("Listening read from client")
	for {
		select {

		case <-c.doneCh:
			// c.server.Del(c)
			c.doneCh <- true
			return

		default:
			var data []byte
			err := websocket.Message.Receive(c.ws, &data)
			switch {
			case err == io.EOF:
				c.doneCh <- true
			case err != nil:
				c.server.Err(err)
			default:
				c.server.SendAll(string(data))
			}
		}
	}
}

func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.Message.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

func (c *Client) Write(msg string) {
	select {
	case c.ch <- msg:
	default:
		c.server.Err(fmt.Errorf("client %d is disconnected.", c.ID))
	}
}
