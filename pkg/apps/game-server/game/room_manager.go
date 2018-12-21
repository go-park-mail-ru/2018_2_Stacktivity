package game

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
	log.Println("Starting room manager")
	pair := make([]*Player, 0)
	filter := make([]*Player, 0)
	for {
		select {
		case player := <-rm.singleplayer:
			log.Println("starting singleplayer...")
			PlayersPendingRoomMetric.With(labelTypeSingle).Dec() // players pending room metric update
			RoomCountMetric.With(labelTypeSingle).Inc()          // room metric update

			room := NewRoom([](*Player){player})
			player.room = room
			go room.Start()
		case p := <-rm.queue:
			if len(pair) == 0 {
				pair = append(pair, p)
			} else if pair[0].user.Username != p.user.Username {
				pair = append(pair, p)
				rm.log.Println("check ws connections")
				for _, p := range pair {
					p.mu.Lock()
					err := p.conn.WriteJSON(models.Message{Event: models.CreateConn})
					p.mu.Unlock()
					if err == nil {
						filter = append(filter, p)
					}
				}
				if len(pair) == len(filter) {
					rm.log.Printf("find game-server: %s vs %s \n", pair[0].user.Username, pair[1].user.Username)
					PlayersPendingRoomMetric.With(labelTypeMult).Sub(2) // players pending room metric update
					RoomCountMetric.With(labelTypeMult).Inc()           // room metric update

					room := NewRoom(pair)
					for _, p := range pair {
						p.room = room
					}
					rm.rooms[room.ID] = room
					go room.Start()
					pair = make([]*Player, 0)
					filter = make([]*Player, 0)
				} else {
					pair = filter
					filter = make([]*Player, 0)
				}
			} else {
				pair[0].conn = p.conn
			}
		case <-rm.stopchan:
			rm.log.Println("stopping game manager...")
			for _, room := range rm.rooms {
				room.stopchanel <- models.Close
			}

			// room metric reset
			RoomCountMetric.With(labelTypeSingle).Set(0)
			RoomCountMetric.With(labelTypeMult).Set(0)
			return
		}
	}
}
