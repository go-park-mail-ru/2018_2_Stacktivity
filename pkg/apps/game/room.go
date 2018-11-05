package game

import (
	"log"

	"github.com/google/uuid"
)

type Room struct {
	ID      uint32
	player1 *Player
	player2 *Player
	rm      *RoomManager
}

func NewRoom(player1 *Player, player2 *Player, rm *RoomManager) *Room {
	log.Println("creating room...")
	return &Room{uuid.New().ID(), player1, player2, rm}
}

func (r *Room) Start() {
	log.Println("starting room...")
}
