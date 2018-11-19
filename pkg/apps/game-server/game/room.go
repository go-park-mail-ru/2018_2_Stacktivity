package game

import (
	"2018_2_Stacktivity/models"
	"log"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID         string
	players    []*Player
	Ticker     *time.Ticker
	Message    chan *IncomingMessage
	Broadcast  chan *models.Message
	Unregister chan *Player
	stopchanel chan interface{}
}

func NewRoom(players []*Player, rm *RoomManager) *Room {
	log.Println("creating game...")
	return &Room{ID: uuid.New().String(),
		players:    players,
		Message:    make(chan *IncomingMessage),
		Broadcast:  make(chan *models.Message),
		Unregister: make(chan *Player),
		stopchanel: make(chan interface{}),
	}
}

func (r *Room) Start() {
	switch len(r.players) {
	case 1:
		// TODO add validate game-server for singleplayer
		log.Println("Start singleplayer")
		go r.players[0].Listen()
	case 2:
		r.players[0].enemy = r.players[1]
		r.players[1].enemy = r.players[0]
		go r.RunBroadcast()
		for _, p := range r.players {
			p.StartMultiplayer()
			go p.Listen()
		}
		go r.ListenToPlayers()
	}

	r.Ticker = time.NewTicker(time.Second)
	for {
		<-r.Ticker.C
		// TODO add some work
	}
}

func (r *Room) RunBroadcast() {
	for {
		m := <-r.Broadcast
		if m.Event == models.Close {
			return
		}
		for _, p := range r.players {
			p.Send(m)
		}
	}
}

func (r *Room) ListenToPlayers() {
	for {
		select {
		case m := <-r.Message:
			log.Printf("Message from player %s: %v", m.Player.user.Username, m.Message)
			switch m.Message.Event {
			case models.UpdateCurv:
				if CheckCurve() {
					UpdateCurve()
					m.Player.enemy.Send(m.Message)
				} else {
					r.Broadcast <- &models.Message{Event: models.InvalidDrop}
				}
			case models.EndCurv:
				m.Player.enemy.Send(m.Message)
				StartCurve(m.Message.Curve)
			}
		case p := <-r.Unregister:
			log.Printf("Player %s exit", p.user.Username)
			msg := &models.Message{
				Event:  models.EndGame,
				Status: &models.StatusSuccess,
			}
			p.enemy.Send(msg)
			return
		case <-r.stopchanel:
			log.Println("Close listening")
			for _, p := range r.players {
				p.Send(&models.Message{Event: models.EndGame, Status: &models.StatusFailure})
				if err := p.conn.Close(); err != nil {
					log.Println("can't close connections", err)
					return
				}
			}
			r.Broadcast <- &models.Message{Event: models.Close}
			return
		}
	}
}
