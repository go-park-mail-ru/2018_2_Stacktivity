package chat_server

import (
	"2018_2_Stacktivity/models"
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	user *models.User
	conn *websocket.Conn
}

type IncomingMessage struct {
	Player  *Player
	Message *models.ChatMessage
}

func NewPlayer(user *models.User, conn *websocket.Conn) *Player {
	return &Player{user: user, conn: conn}
}

func (p *Player) SetConnection(conn *websocket.Conn) {
	p.conn = conn
}

func (p *Player) Listen(ch chan *IncomingMessage) {
	for {
		m := &models.ChatMessage{}
		err := p.conn.ReadJSON(m)
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("player %d was disconnected", p.user.ID)
			p.conn = nil
			return
		}
		im := &IncomingMessage{
			Player:  p,
			Message: m,
		}

		ch <- im
	}
}

func (p *Player) Send(s *models.ChatMessage) {
	if p.conn != nil {
		if err := p.conn.WriteJSON(s); err != nil {
			log.Printf("can't send message to player %s\n", p.user.Username)
			p.conn = nil
		}
	}
}

func (p *Player) isOnline() bool {
	return p.conn == nil
}
