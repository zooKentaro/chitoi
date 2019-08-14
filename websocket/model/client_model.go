package model

import (
	"fmt"
	"io"
	"log"

	"github.com/uenoryo/chitoi/websocket/packet"
	"golang.org/x/net/websocket"
)

type Client struct {
	ws          *websocket.Conn
	room        *Room
	Player      *User
	IsListening bool
	ch          chan *packet.BloadcastPacket
	doneCh      chan bool
}

func NewClient(ws *websocket.Conn, room *Room, player *User) *Client {
	var (
		doneCh = make(chan bool)
		ch     = make(chan *packet.BloadcastPacket)
	)
	return &Client{
		Player:      player,
		IsListening: false,
		ws:          ws,
		room:        room,
		ch:          ch,
		doneCh:      doneCh,
	}
}

func (c *Client) ID() uint64 {
	return c.Player.Row.ID
}

func (c *Client) Listen() {
	if c.IsListening {
		return
	}

	c.IsListening = true
	go c.listenWrite()

	log.Println("Listening read from client")

	if err := c.room.SubmitOnEnterPlayer(); err != nil {
		fmt.Println("[ERROR] submit on exit player failed", err.Error())
	}
	for {
		select {

		case <-c.doneCh:
			// c.server.Del(c)
			c.doneCh <- true
			return

		default:
			request := &packet.RequestPacket{}
			err := websocket.JSON.Receive(c.ws, &request)
			switch {
			case err == io.EOF:
				fmt.Println("close listenning for reading, player id:", c.Player.Row.ID)
				c.room.PushOut(c)
				if err := c.room.SubmitOnExitPlayer(request); err != nil {
					fmt.Println("[ERROR] submit on exit player failed", err.Error())
				}
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
				userID, err := c.room.Authenticate(request.SessionID)
				if err != nil {
					server, err2 := c.room.Server()
					if err2 != nil {
						fmt.Println("[ERROR] error ocurred in room, but failed to get server", err2.Error())
					} else {
						server.Err(err)
					}
				} else {
					request.SenderID = userID
					if err := c.room.Submit(request); err != nil {
						fmt.Println("[ERROR] invalid packet", err.Error())
					}
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
		case pkt := <-c.ch:
			log.Println("Send:", pkt)
			websocket.JSON.Send(c.ws, pkt)

		// receive done request
		case <-c.doneCh:
			fmt.Println("close listenning for writing, player id:", c.Player.Row.ID)
			return
		}
	}
}

func (c *Client) Receive(pkt *packet.BloadcastPacket) {
	select {
	case c.ch <- pkt:
		log.Println(pkt, "を受け取ったよ！")
	default:
		// c.server.Err(fmt.Errorf("client %d is disconnected.", c.ID))
	}
}
