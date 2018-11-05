package game

import (
	log "github.com/sirupsen/logrus"
)

type RoomManager struct {
	queue    chan *Player
	stopchan chan struct{}
	log      *log.Logger
}

func NewRoomManager(logger *log.Logger) *RoomManager {
	return &RoomManager{
		queue:    make(chan *Player),
		stopchan: make(chan struct{}),
		log:      logger,
	}
}

func (rm *RoomManager) Run() {
	rm.log.Println("starting room manager...")
	for {
		select {
		default:
			player1 := <-rm.queue
			player2 := <-rm.queue
			rm.log.Printf("find game: %s vs %s \n", player1.user.Username, player2.user.Username)
			room := NewRoom(player1, player2, rm)
			go room.Start()
		case <-rm.stopchan:
			return
		}
	}
}
