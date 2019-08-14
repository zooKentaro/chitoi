package model

import (
	"log"

	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/database/row"
	"github.com/uenoryo/chitoi/websocket/packet"
)

// Server (､´･ω･)▄︻┻┳═一
type Server struct {
	core        *core.Core
	doneCh      chan bool
	errCh       chan error
	broadCastCh chan *packet.BloadcastPacket
	rooms       map[uint32]*Room
}

// NewServer (､´･ω･)▄︻┻┳═一
func NewServer(core *core.Core) *Server {
	var (
		doneCh      = make(chan bool)
		errCh       = make(chan error)
		broadCastCh = make(chan *packet.BloadcastPacket)
		rooms       = make(map[uint32]*Room)
	)
	return &Server{
		core,
		doneCh,
		errCh,
		broadCastCh,
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

// Receive はroom codeの部屋のメンバーにpacketを送信する
func (s *Server) Receive(pkt *packet.Packet) {
	room, ok := s.rooms[pkt.RoomCode]
	if !ok {
		return
	}

	var (
		player1 *row.User
		player2 *row.User
	)
	if room.Player1 != nil {
		if client, ok := room.Clients[room.Player1.Row.ID]; ok {
			player1 = client.Player.Row
		}
	}
	if room.Player2 != nil {
		if client, ok := room.Clients[room.Player2.Row.ID]; ok {
			player2 = client.Player.Row
		}
	}
	bloadcastPacket := &packet.BloadcastPacket{
		pkt,
		player1,
		player2,
	}
	s.broadCastCh <- bloadcastPacket
}

func (s *Server) Validate(pkt *packet.Packet) error {
	if pkt.RoomCode == 0 {
		return errors.New("error room code is empty")
	}
	if pkt.SenderID == 0 {
		return errors.New("error sender id (user id) is empty")
	}
	room, ok := s.rooms[pkt.RoomCode]
	if !ok {
		return errors.Errorf("invalid room code:%d", pkt.RoomCode)
	}
	if _, ok := room.Clients[pkt.SenderID]; !ok {
		return errors.Errorf("invalid sender id:%d", pkt.SenderID)
	}
	return nil
}

// Err はserverにerrを通知する
func (s *Server) Err(err error) {
	s.errCh <- err
}
