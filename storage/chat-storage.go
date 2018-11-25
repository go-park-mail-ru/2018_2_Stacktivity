package storage

import (
	"2018_2_Stacktivity/models"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ChatStorageI interface {
	CreateChat() (*models.Chat, error)
	GetChatByID(ID int) (*models.Chat, error)
	GetChatUsers(ID int) (*[]models.User, error)
	GetUsersChat(Username string) (*[]models.Chat, error)
	AddUserToChat(ChatID int, Username string) error

	AddMessage(message *models.ChatMessage, userID int, chatID int) (*models.ChatMessage, error)
	GetMessageHisory(ChatID int) (*[]models.ChatMessage, error)

	GetAllChats() (*[]models.Chat, error)
}

type ChatStorage struct {
	DB *sqlx.DB
}

func GetChatStorage() *ChatStorage {
	storage := &ChatStorage{
		DB: db,
	}
	return storage
}

var createChat = `INSERT INTO chat(chat_name, chat_type) VALUES ($1, $2) RETURNING id;`

func (cs *ChatStorage) CreateChat() (*models.Chat, error) {
	var id int
	if err := cs.DB.QueryRow(createChat, "noname", "public").Scan(&id); err != nil {
		return nil, errors.Wrap(err, "can't select from chat")
	}
	return &models.Chat{
			ID:   id,
			Name: "noname",
		},
		nil
}

var getChatByID = `SELECT id, chat_name FROM chat WHERE id = $1;`

func (cs *ChatStorage) GetChatByID(ID int) (*models.Chat, error) {
	var chat models.Chat
	if err := cs.DB.QueryRow(getChatByID, ID).Scan(&chat.ID, &chat.Name); err != nil {
		return nil, err
	}
	return &chat, nil
}

var getChatUsersID = `SELECT user_id FROM chat_user WHERE chat_ID = $1;`
var getUserByID = `SELECT username, email, score FROM "user" WHERE uid = $1;`

func (cs *ChatStorage) GetChatUsers(ID int) (*[]models.User, error) {
	var users []models.User
	rows, err := cs.DB.Query(getChatUsersID, ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID); err != nil {
			return nil, err
		}
		if err := cs.DB.QueryRow(getUserByID, user.ID).Scan(&user.Username, &user.Email, &user.Score); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}

var addUserToChat = `INSERT INTO chat_user(chat_id, user_id) VALUES ($1, $2);`
var getUserIDbyUsername = `SELECT uid FROM "user" WHERE username = $1;`

func (cs *ChatStorage) AddUserToChat(ChatID int, Username string) error {
	var id int
	if err := cs.DB.QueryRow(getUserIDbyUsername, Username).Scan(&id); err != nil {
		return err
	}
	if _, err := cs.DB.Exec(addUserToChat, ChatID, id); err != nil {
		return err
	}
	return nil
}

var addMessage = `INSERT INTO chat_message(user_id, chat_id, chat_text) VALUES ($1, $2, $3) RETURNING (created_at);`

func (cs *ChatStorage) AddMessage(message *models.ChatMessage, userID int, chatID int) (*models.ChatMessage, error) {
	newMessage := *message
	if err := cs.DB.QueryRow(addMessage, userID, chatID, message.Text).Scan(&newMessage.Created); err != nil {
		return nil, err
	}
	return &newMessage, nil
}

var getMessages = `SELECT user_id, chat_text, created_at FROM chat_message WHERE chat_id = $1;`

func (cs *ChatStorage) GetMessageHisory(ChatID int) (*[]models.ChatMessage, error) {
	messages := make([]models.ChatMessage, 0)
	rows, err := cs.DB.Query(getMessages, ChatID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var message models.ChatMessage
		var userID int
		if err := rows.Scan(&userID, &message.Text, &message.Created); err != nil {
			return nil, err
		}
		var user models.User
		if err := cs.DB.QueryRow(getUserByID, userID).Scan(&user.Username, &user.Email, &user.Score); err != nil {
			return nil, err
		}
		message.User = &user
		messages = append(messages, message)
	}
	return &messages, nil
}

var getUsersChat = `SELECT chat_id FROM chat_user WHERE user_id = $1;`

func (cs *ChatStorage) GetUsersChat(Username string) (*[]models.Chat, error) {
	var userID int
	chats := make([]models.Chat, 0)
	if err := cs.DB.QueryRow(getUserIDbyUsername, Username).Scan(&userID); err != nil {
		return nil, err
	}
	rows, err := cs.DB.Query(getUsersChat, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID); err != nil {
			return nil, err
		}
		if err := cs.DB.QueryRow(getChatByID, chat.ID).Scan(&chat.ID, &chat.Name); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return &chats, nil
}

var getAllChats = `SELECT id, chat_name FROM chat;`

func (cs *ChatStorage) GetAllChats() (*[]models.Chat, error) {
	chats := make([]models.Chat, 0)
	rows, err := cs.DB.Query(getAllChats)
	if err != nil {

	}
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.Name); err != nil {
			return nil, err
		}
		users, err := cs.GetChatUsers(chat.ID)
		if err != nil {
			return nil, err
		}
		chat.Members = make([]string, 0)
		for _, user := range *users {
			chat.Members = append(chat.Members, user.Username)
		}
		chats = append(chats, chat)
	}
	return &chats, nil
}
