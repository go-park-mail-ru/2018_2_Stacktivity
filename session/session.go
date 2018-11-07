package session

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Session struct {
	ID int
}

type SessionID struct {
	ID uuid.UUID
}

type SessionManager struct {
	sessions redis.Conn
	log      *log.Logger
}

type SessionManagerI interface {
	Create(*Session) (*SessionID, error)
	Check(*SessionID) (bool, *Session)
	Delete(*SessionID)
}

func NewSessionManager(logger *log.Logger, conn redis.Conn) *SessionManager {
	return &SessionManager{
		sessions: conn,
		log:      logger,
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	ID, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "Can't create session ID")
	}
	sessionID := SessionID{ID}
	dataSerialized, err := json.Marshal(in)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal session")
	}
	mkey := "sessions:" + sessionID.ID.String()
	result, err := redis.String(sm.sessions.Do("SET", mkey, dataSerialized, "EX", 86400))
	if err != nil {
		return nil, errors.Wrap(err, "can't insert valuer into redis")
	}
	if result != "OK" {
		return nil, errors.New("result from redis is not OK: " + result)
	}
	return &sessionID, nil
}

func (sm *SessionManager) Check(in *SessionID) (bool, *Session) {
	mkey := "sessions:" + in.ID.String()
	data, err := redis.Bytes(sm.sessions.Do("GET", mkey))
	if err != nil {
		sm.log.Println("can't get data:", err)
		return false, nil
	}
	session := &Session{}
	err = json.Unmarshal(data, session)
	if err != nil {
		sm.log.Println("can't unpack session data:", err)
		return false, nil
	}
	return session != nil, session
}

func (sm *SessionManager) Delete(in *SessionID) {
	mkey := "sessions:" + in.ID.String()
	_, err := redis.Int(sm.sessions.Do("DEL", mkey))
	if err != nil {
		sm.log.Println("redis error:", err)
	}
}
