package model

import (
	"fmt"
	"io"
	"log"

	"github.com/uenoryo/chitoi/model"
	"golang.org/x/net/websocket"
)

type Client struct {
	ws          *websocket.Conn
	room        *Room
	ID          uint64
	IsListening bool
	ch          chan *Packet
	doneCh      chan bool
}

// Packet は送受信1回分のデータ
type Packet struct {
	SenderID   uint64 `json:"sender_id,string"`
	ActionType uint32 `json:"action_type"`
	RoomCode   uint32
}

func NewClient(ws *websocket.Conn, room *Room, user *model.User) *Client {
	var (
		doneCh = make(chan bool)
		ch     = make(chan *Packet)
	)
	return &Client{
		ID:          user.Row.ID,
		IsListening: false,
		ws:          ws,
		room:        room,
		ch:          ch,
		doneCh:      doneCh,
	}
}

func (c *Client) Listen() {
	if c.IsListening {
		return
	}

	c.IsListening = true
	go c.listenWrite()

	log.Println("Listening read from client")
	for {
		select {

		case <-c.doneCh:
			// c.server.Del(c)
			c.doneCh <- true
			return

		default:
			packet := &Packet{}
			err := websocket.JSON.Receive(c.ws, &packet)
			switch {
			case err == io.EOF:
				fmt.Println("close listenning for reading, id:", c.ID)
				c.doneCh <- true
				return
			case err != nil:
				server, err2 := c.room.Server()
				if err2 != nil {
					fmt.Println("[ERROR] error ocurred in room, but failed to get server", err2.Error())
				} else {
					server.Err(err)
				}
			default:
				if err := c.room.Submit(packet); err != nil {
					fmt.Println("[ERROR] invalid packet", err.Error())
				}
			}
		}
	}
}

func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case packet := <-c.ch:
			log.Println("Send:", packet)
			websocket.JSON.Send(c.ws, packet)

		// receive done request
		case <-c.doneCh:
			fmt.Println("close listenning for writing, id:", c.ID)
			return
		}
	}
}

func (c *Client) Receive(packet *Packet) {
	select {
	case c.ch <- packet:
		log.Println(packet, "を受け取ったよ！")
	default:
		// c.server.Err(fmt.Errorf("client %d is disconnected.", c.ID))
	}
}
