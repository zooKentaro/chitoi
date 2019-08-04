package model

import (
	"io"
	"log"

	"github.com/uenoryo/chitoi/model"
	"golang.org/x/net/websocket"
)

type Client struct {
	ws     *websocket.Conn
	room   *Room
	ID     uint64
	ch     chan string
	doneCh chan bool
}

func NewClient(ws *websocket.Conn, room *Room, user *model.User) *Client {
	var (
		doneCh = make(chan bool)
		ch     = make(chan string)
	)
	return &Client{
		ID:     user.Row.ID,
		ws:     ws,
		room:   room,
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
				// c.server.Err(err)
			default:
				// c.server.SendAll(string(data))
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
		// c.server.Err(fmt.Errorf("client %d is disconnected.", c.ID))
	}
}
