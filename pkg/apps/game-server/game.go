package game_server

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
	return &Game{
		NewRoomManager(logger),
		logger,
	}
}

func (g *Game) RunSinglePlayer(user *models.User, ws *websocket.Conn) {
	player := NewPlayer(user, ws)
	g.rm.singleplayer <- player
}

func (g *Game) AddPlayer(user *models.User, ws *websocket.Conn) {
	player := NewPlayer(user, ws)
	g.rm.queue <- player
}

func (g *Game) Start() {
	g.log.Println("starting game-server...")
	go g.rm.Run()
}

func (g *Game) Stop() {
	g.log.Println("stopping game-server...")
	g.rm.stopchan <- 0
}
