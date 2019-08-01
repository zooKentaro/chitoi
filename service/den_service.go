package service

import (
	"fmt"
	"io"
	"log"

	"github.com/uenoryo/chitoi/core"
	"golang.org/x/net/websocket"
)

type DenService interface {
	Listener() *Listener
	Entry(*websocket.Conn)
}

var maxID uint64 = 1001

type denService struct {
	core   *core.Core
	server *Listener
}

// NewDenService (､´･ω･)▄︻┻┳═一
func NewDenService(core *core.Core) DenService {
	listener := newListener()

	return &denService{
		core:   core,
		server: listener,
	}
}

func (srv *denService) Listener() *Listener {
	return srv.server
}

func (srv *denService) Entry(ws *websocket.Conn) {
	client := NewClient(ws, srv.server)
	srv.server.Add(client)
	client.Listen()
}

func newListener() *Listener {
	var (
		doneCh    = make(chan bool)
		errCh     = make(chan error)
		sendAllCh = make(chan string)
		clients   = make(map[uint64]*Client)
	)
	return &Listener{
		doneCh,
		errCh,
		sendAllCh,
		clients,
	}
}

type Listener struct {
	doneCh    chan bool
	errCh     chan error
	sendAllCh chan string
	clients   map[uint64]*Client
}

func (l *Listener) Listen() {
	for {
		select {
		case err := <-l.errCh:
			log.Println("Error:", err.Error())

		case <-l.doneCh:
			return

		case msg := <-l.sendAllCh:
			log.Println("Send all:", msg)
			for id, client := range l.clients {
				client.Write(fmt.Sprintf("message: %s, id:%d", msg, id))
			}
		}
	}
}

func (l *Listener) Add(c *Client) {
	log.Println("ADDED")
	l.clients[c.ID] = c
}

func (l *Listener) SendAll(msg string) {
	l.sendAllCh <- msg
}

func (l *Listener) Err(err error) {
	l.errCh <- err
}

type Client struct {
	ID       uint64
	ws       *websocket.Conn
	listener *Listener
	ch       chan string
	doneCh   chan bool
}

func NewClient(ws *websocket.Conn, listener *Listener) *Client {
	maxID++
	var (
		id     = maxID
		doneCh = make(chan bool)
		ch     = make(chan string)
	)

	return &Client{
		ID:       id,
		ws:       ws,
		listener: listener,
		ch:       ch,
		doneCh:   doneCh,
	}
}

func (c *Client) Listen() {
	go c.listenWrite()

	log.Println("Listening read from client")
	for {
		select {

		case <-c.doneCh:
			// c.listener.Del(c)
			c.doneCh <- true
			return

		default:
			var data []byte
			err := websocket.Message.Receive(c.ws, &data)
			switch {
			case err == io.EOF:
				c.doneCh <- true
			case err != nil:
				c.listener.Err(err)
			default:
				c.listener.SendAll(string(data))
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
		c.listener.Err(fmt.Errorf("client %d is disconnected.", c.ID))
	}
}
