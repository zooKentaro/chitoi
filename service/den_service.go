package service

import (
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
	)
	return &Listener{
		doneCh,
		errCh,
		sendAllCh,
	}
}

type Listener struct {
	doneCh    chan bool
	errCh     chan error
	sendAllCh chan string
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
		}
	}
}

func (l *Listener) Add(c *Client) {
	log.Println("ADDED")
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
	doneCh   chan bool
}

func NewClient(ws *websocket.Conn, listener *Listener) *Client {
	maxID++
	var (
		id     = maxID
		doneCh = make(chan bool)
	)

	return &Client{
		ID:       id,
		ws:       ws,
		listener: listener,
		doneCh:   doneCh,
	}
}

func (c *Client) Listen() {
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
