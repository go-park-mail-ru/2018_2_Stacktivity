package chat_server

import (
	"2018_2_Stacktivity/models"
	"log"
)

type Chat struct {
	ID       int
	Name     string
	players  []*Player
	history  []*models.ChatMessage
	chatChan chan *IncomingMessage
	stopChan chan interface{}
}

func CreateChat(ID int) *Chat {
	return &Chat{
		ID:       ID,
		Name:     "noname",
		players:  make([]*Player, 0),
		history:  make([]*models.ChatMessage, 0),
		chatChan: make(chan *IncomingMessage, 256),
		stopChan: make(chan interface{}),
	}
}

func (c *Chat) AddPlayerToChat(p *Player) {
	c.players = append(c.players, p)
}

func (c *Chat) RunBroadcast() {
	for {
		m := <-c.chatChan
		if m.Message.Event == models.Close {
			return
		}
		log.Println("broadcast...")
		for _, p := range c.players {
			p.Send(m.Message)
		}
	}
}
