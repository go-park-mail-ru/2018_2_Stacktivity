package game

import (
	"2018_2_Stacktivity/models"

	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	user *models.User
	conn *websocket.Conn
}

func NewPlayer(user *models.User, conn *websocket.Conn) *Player {
	log.Println("creating player...")
	return &Player{user, conn}
}
