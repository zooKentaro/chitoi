package model

import (
	"fmt"
	"log"

	"github.com/uenoryo/chitoi/core"
)

// Server (､´･ω･)▄︻┻┳═一
type Server struct {
	core      *core.Core
	doneCh    chan bool
	errCh     chan error
	sendAllCh chan string
	clients   map[uint64]*Client
	rooms     map[uint32]*Room
}

// NewServer (､´･ω･)▄︻┻┳═一
func NewServer(core *core.Core) *Server {
	var (
		doneCh    = make(chan bool)
		errCh     = make(chan error)
		sendAllCh = make(chan string)
		clients   = make(map[uint64]*Client)
		rooms     = make(map[uint32]*Room)
	)
	return &Server{
		core,
		doneCh,
		errCh,
		sendAllCh,
		clients,
		rooms,
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

// Launch はserverにroomを立てます
func (s *Server) Launch(room *Room) {
	if s.IsLaunched(room) {
		return
	}
	s.rooms[room.Row.Code] = room
}

// IsLaunched はserverにroomが立っているかどうかを返す
func (s *Server) IsLaunched(room *Room) bool {
	_, isLaunched := s.rooms[room.Row.Code]
	return isLaunched
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
