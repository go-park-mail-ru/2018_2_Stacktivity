package chat_server

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/storage"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ChatManager struct {
	chats       map[int]*Chat
	stopchan    chan interface{}
	log         *log.Logger
	members     map[string]*Player
	allMessages chan *IncomingMessage
	chatStorage storage.ChatStorageI
}

func NewChatManager(logger *log.Logger) *ChatManager {
	return &ChatManager{
		chats:       make(map[int]*Chat),
		members:     make(map[string]*Player),
		stopchan:    make(chan interface{}),
		allMessages: make(chan *IncomingMessage, 512),
		log:         logger,
		chatStorage: storage.GetChatStorage(),
	}
}

func (cm *ChatManager) Run() {
	log.Println("Starting chat manager")
	chats, err := cm.chatStorage.GetAllChats()
	if err != nil {
		cm.log.Println("can't get chats")
		return
	}
	for i := range *chats {
		log.Println((*chats)[i].ID)
		cm.chats[(*chats)[i].ID] = CreateChat((*chats)[i].ID)
	}
	for {
		select {
		case m := <-cm.allMessages:
			cm.log.Println("New message:")
			switch m.Message.Event {
			case models.SendMessage:
				cm.log.Println("Sending message to chat ", m.Message.Chat)
				chat, ok := cm.chats[m.Message.Chat]
				if ok {
					cm.log.Println("OK")
					m.Message = cm.NewMessage(chat, m)
					m.Message.Event = models.SendMessageSuccess
					chat.chatChan <- m
				}
			case models.ConnectToChat:
				cm.log.Println("New connect to chat")
				chat, ok := cm.chats[m.Message.Chat]
				if ok {
					p, ok := cm.members[m.Message.NewUsername]
					if ok {
						cm.AddPlayerToChat(chat, p)
						m.Message.Event = models.ConnectToChatSuccess
						m.Message.Chat = chat.ID
						chat.chatChan <- m
					}
				}
			case models.CreateChat:
				cm.log.Println("New chat")
				chat := cm.CreateChat()
				if chat != nil {
					cm.chats[chat.ID] = chat
					go chat.RunBroadcast()
					cm.AddPlayerToChat(chat, m.Player)
					m.Message.Chat = chat.ID
					m.Message.Event = models.ConnectToChatSuccess
					chat.chatChan <- m
				}
			}
		case <-cm.stopchan:
			cm.log.Println("stopping chat manager...")
			for _, chat := range cm.chats {
				chat.stopChan <- models.Close
			}
			return
		}
	}
}

func (cm *ChatManager) CreatePlayer(user *models.User, conn *websocket.Conn) {
	player, ok := cm.members[user.Username]
	if !ok {
		cm.log.Println("create new player")
		player = NewPlayer(user, conn)
	}
	player.conn = conn
	cm.members[user.Username] = player
	chats, err := cm.chatStorage.GetUsersChat(user.Username)
	if err != nil {
		cm.log.Println("can't get users chats")
		return
	}
	for i := range *chats {
		messages, err := cm.chatStorage.GetMessageHisory((*chats)[i].ID)
		if err != nil {
			cm.log.Println("can't get messages history")
			return
		}
		(*chats)[i].History = *messages
	}
	go player.Listen(cm.allMessages)
	message := &models.ChatMessage{
		Event: models.ConnectToChatSuccess,
		Data:  chats,
	}
	player.Send(message)
}

func (cm *ChatManager) CreateChat() *Chat {
	newChat, err := cm.chatStorage.CreateChat()
	if err != nil {
		cm.log.Println("can't create new chat")
		return nil
	}
	return CreateChat(newChat.ID)
}

func (cm *ChatManager) AddPlayerToChat(c *Chat, p *Player) {
	if err := cm.chatStorage.AddUserToChat(c.ID, p.user.Username); err != nil {
		cm.log.Println("can't add user to chat")
		return
	}
	c.AddPlayerToChat(p)
}

func (cm *ChatManager) NewMessage(c *Chat, m *IncomingMessage) *models.ChatMessage {
	newMessage, err := cm.chatStorage.AddMessage(m.Message, int(m.Player.user.ID), c.ID)
	if err != nil {
		cm.log.Println("can't add new message")
		return nil
	}
	newMessage.Event = models.SendMessageSuccess
	return newMessage
}
