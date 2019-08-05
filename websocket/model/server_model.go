package model

import (
	"log"

	"github.com/uenoryo/chitoi/core"
)

// Server (､´･ω･)▄︻┻┳═一
type Server struct {
	core        *core.Core
	doneCh      chan bool
	errCh       chan error
	broadCastCh chan *Packet
	clients     map[uint64]*Client
	rooms       map[uint32]*Room
}

// NewServer (､´･ω･)▄︻┻┳═一
func NewServer(core *core.Core) *Server {
	var (
		doneCh      = make(chan bool)
		errCh       = make(chan error)
		broadCastCh = make(chan *Packet)
		clients     = make(map[uint64]*Client)
		rooms       = make(map[uint32]*Room)
	)
	return &Server{
		core,
		doneCh,
		errCh,
		broadCastCh,
		clients,
		rooms,
	}
}

// Listen はイベントを監視する
func (s *Server) Listen() {
	for {
		select {
		case err := <-s.errCh:
			log.Println("[ERROR]", err.Error())

		case <-s.doneCh:
			return

		case packet := <-s.broadCastCh:
			if room, ok := s.rooms[packet.RoomCode]; ok {
				room.SendToMembers(packet)
			}
		}
	}
}

// Launch はserverにroomを立てます
func (s *Server) Launch(room *Room) *Room {
	if exists, ok := s.rooms[room.Row.Code]; ok {
		return exists
	}
	room.RegisterServer(s)
	s.rooms[room.Row.Code] = room
	return room
}

// LaunchedRoom は roomCode の room を返す
// 起動していなければ nil を返す
func (s *Server) LaunchedRoom(roomCode uint32) *Room {
	return s.rooms[roomCode]
}

// Add は client を server に追加する
func (s *Server) Add(client *Client) {
	log.Println("ADDED")
	s.clients[client.ID] = client
}

// Receive はroom codeの部屋のメンバーにpacketを送信する
func (s *Server) Receive(packet *Packet) {
	s.broadCastCh <- packet
}

// Err はserverにerrを通知する
func (s *Server) Err(err error) {
	s.errCh <- err
}
