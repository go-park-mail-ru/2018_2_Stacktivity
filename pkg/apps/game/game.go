package game

import (
	"2018_2_Stacktivity/models"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Game struct {
	rm  *RoomManager
	log *log.Logger
}

func NewGame(logger *log.Logger) *Game {
	log.Println("creating game...")
	return &Game{
		NewRoomManager(logger),
		logger,
	}
}

func (g *Game) AddPlayer(user *models.User, ws *websocket.Conn) {
	g.log.Println("adding player...")
	player := NewPlayer(user, ws)
	g.rm.queue <- player
}

func (g *Game) Start() {
	g.log.Println("starting game...")
	go g.rm.Run()
}
