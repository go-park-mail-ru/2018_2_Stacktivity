package game

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/storage"
	"encoding/json"
	"log"

	"sync"

	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	mu     sync.Mutex
	user   *models.User
	enemy  *Player
	room   *Room
	conn   *websocket.Conn
	logic  models.PlayerLogic
	isOpen bool
}

type IncomingMessage struct {
	Player  *Player
	Message *models.Message
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func NewPlayer(user *models.User, conn *websocket.Conn) *Player {
	return &Player{mu: sync.Mutex{}, user: user, conn: conn, isOpen: true}
}

func (p *Player) CheckConn() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
	}()

	for {
		<-ticker.C
		p.mu.Lock()
		p.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			p.mu.Unlock()
			p.isOpen = false
			return
		}
		p.mu.Unlock()
	}
}

func (p *Player) Listen() {
	defer p.conn.Close()
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		if p.room != nil {
			m := &models.Message{}
			err := p.conn.ReadJSON(m)
			if websocket.IsUnexpectedCloseError(err) {
				log.Printf("player %d was disconnected", p.user.ID)
				p.room.Unregister <- p
				p.isOpen = false
				log.Println("player deleted")
				return
			}
			im := &IncomingMessage{
				Player:  p,
				Message: m,
			}
			p.room.Message <- im
		} else {
			_, _, err := p.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					p.isOpen = false
					return
				}
				break
			}
		}
	}
}

func (p *Player) Send(s *models.Message) {
	p.mu.Lock()
	p.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := p.conn.WriteJSON(s)
	p.mu.Unlock()
	if err != nil {
		log.Printf("can't send message to player %s\n", p.user.Username)
		log.Println(err.Error())
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
