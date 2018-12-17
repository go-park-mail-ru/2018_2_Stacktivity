package game

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/storage"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	user  *models.User
	enemy *Player
	room  *Room
	conn  *websocket.Conn
	logic models.PlayerLogic
}

type IncomingMessage struct {
	Player  *Player
	Message *models.Message
}

func NewPlayer(user *models.User, conn *websocket.Conn) *Player {
	return &Player{user: user, conn: conn}
}

func (p *Player) Listen() {
	for {
		m := &models.Message{}
		err := p.conn.ReadJSON(m)
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("player %d was disconnected", p.user.ID)
			p.room.Unregister <- p
			log.Println("player deleted")
			return
		}
		im := &IncomingMessage{
			Player:  p,
			Message: m,
		}
		p.room.Message <- im
	}
}

func (p *Player) Send(s *models.Message) {
	if err := p.conn.WriteJSON(s); err != nil {
		log.Printf("can't send message to player %s\n", p.user.Username)
		p.room.Unregister <- p
	}
}

func (p *Player) StartMultiplayer() {
	players := make([]string, 2)
	players[0] = p.user.Username
	players[1] = p.enemy.user.Username

	log.Println("Get level ", p.room.levelNum)
	var level models.Level

	dbLevel, err := storage.GetUserStorage().GetLevelByNumber(p.room.levelNum)
	if err != nil {
		log.Println("PIZDA RULIU")
		log.Println(err.Error())
		return
	}
	if err := json.Unmarshal([]byte(dbLevel.Level), &level); err != nil {
		log.Println("HUITA KAKAYA-TO")
		log.Println(err.Error())
		return
	}

	m := &models.Message{
		Event:   models.StartGame,
		Players: &players,
		Level:   &level,
	}
	p.Send(m)
}
