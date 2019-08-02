package model

import (
	"fmt"
	"log"
)

// Server (､´･ω･)▄︻┻┳═一
type Server struct {
	doneCh    chan bool
	errCh     chan error
	sendAllCh chan string
	clients   map[uint64]*Client
}

// NewServer (､´･ω･)▄︻┻┳═一
func NewServer() *Server {
	var (
		doneCh    = make(chan bool)
		errCh     = make(chan error)
		sendAllCh = make(chan string)
		clients   = make(map[uint64]*Client)
	)
	return &Server{
		doneCh,
		errCh,
		sendAllCh,
		clients,
	}
}

// Listen はイベントを監視する
func (s *Server) Listen() {
	for {
		select {
		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return

		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			for id, client := range s.clients {
				client.Write(fmt.Sprintf("message: %s, id:%d", msg, id))
			}
		}
	}
}

// Add は client を server に追加する
func (s *Server) Add(client *Client) {
	log.Println("ADDED")
	s.clients[client.ID] = client
}

// SendAll は全ての client にメッセージを送信する
func (s *Server) SendAll(msg string) {
	s.sendAllCh <- msg
}

// Err はserverにerrを通知する
func (s *Server) Err(err error) {
	s.errCh <- err
}
