package model

import (
	"fmt"
	"io"
	"log"

	"github.com/uenoryo/chitoi/database/row"
	"golang.org/x/net/websocket"
)

type Client struct {
	ws          *websocket.Conn
	room        *Room
	Player      *User
	IsListening bool
	ch          chan *BloadcastPacket
	doneCh      chan bool
}

// Packet は各クライアントから送信される1回分のデータ
type Packet struct {
	SessionID  string `json:"session_id"`
	ActionType uint32 `json:"action_type"`
	SenderID   uint64
	RoomCode   uint32
}

// BloadcastPacket は全体に送信するデータ
type BloadcastPacket struct {
	*Packet
	Player1 *row.User `json:"player1"`
	Player2 *row.User `json:"player2"`
}

func NewClient(ws *websocket.Conn, room *Room, player *User) *Client {
	var (
		doneCh = make(chan bool)
		ch     = make(chan *BloadcastPacket)
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
				fmt.Println("close listenning for reading, player id:", c.Player.Row.ID)
				c.room.PushOut(c)
				if err := c.room.SubmitOnExitPlayer(packet); err != nil {
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
				userID, err := c.room.Authenticate(packet.SessionID)
				if err != nil {
					server, err2 := c.room.Server()
					if err2 != nil {
						fmt.Println("[ERROR] error ocurred in room, but failed to get server", err2.Error())
					} else {
						server.Err(err)
					}
				} else {
					packet.SenderID = userID
					if err := c.room.Submit(packet); err != nil {
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
		case packet := <-c.ch:
			log.Println("Send:", packet)
			websocket.JSON.Send(c.ws, packet)

		// receive done request
		case <-c.doneCh:
			fmt.Println("close listenning for writing, player id:", c.Player.Row.ID)
			return
		}
	}
}

func (c *Client) Receive(packet *BloadcastPacket) {
	select {
	case c.ch <- packet:
		log.Println(packet, "を受け取ったよ！")
	default:
		// c.server.Err(fmt.Errorf("client %d is disconnected.", c.ID))
	}
}
