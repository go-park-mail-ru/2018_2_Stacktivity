package game_server

import (
	"2018_2_Stacktivity/models"

	log "github.com/sirupsen/logrus"
)

type RoomManager struct {
	rooms        map[string]*Room
	queue        chan *Player
	stopchan     chan interface{}
	singleplayer chan *Player
	log          *log.Logger
}

func NewRoomManager(logger *log.Logger) *RoomManager {
	return &RoomManager{
		rooms:        make(map[string]*Room),
		queue:        make(chan *Player),
		stopchan:     make(chan interface{}),
		singleplayer: make(chan *Player),
		log:          logger,
	}
}

func (rm *RoomManager) Run() {
	pair := make([]*Player, 0)
	for {
		select {
		case player := <-rm.singleplayer:
			log.Println("starting singleplayer...")
			room := NewRoom([](*Player){player}, rm)
			go room.Start()
		case p := <-rm.queue:
			pair = append(pair, p)
			if len(pair) == 2 {
				rm.log.Printf("find game-server: %s vs %s \n", pair[0].user.Username, pair[1].user.Username)
				room := NewRoom(pair, rm)
				for _, p := range pair {
					p.room = room
				}
				rm.rooms[room.ID] = room
				go room.Start()
				pair = make([]*Player, 0)
			}
		case <-rm.stopchan:
			rm.log.Println("stopping room manager...")
			for _, room := range rm.rooms {
				room.stopchanel <- models.Close
			}
			return
		}
	}
}
