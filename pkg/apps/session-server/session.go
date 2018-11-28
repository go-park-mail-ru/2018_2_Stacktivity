package session_server

import (
	"2018_2_Stacktivity/pkg/session"
	"context"
	"encoding/json"

	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type SessionManager struct {
	sessions redis.Conn
}

func NewSessionManager(conn *redis.Conn) *SessionManager {
	return &SessionManager{
		sessions: *conn,
	}
}

func (sm *SessionManager) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
	log.Println("create session")
	ID, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "Can't create session-server ID")
	}
	sessionID := session.SessionID{ID: ID.String()}
	dataSerialized, err := json.Marshal(in)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal session-server")
	}
	mkey := "sessions:" + sessionID.ID
	result, err := redis.String(sm.sessions.Do("SET", mkey, dataSerialized, "EX", 86400))
	if err != nil {
		return nil, errors.Wrap(err, "can't insert valuer into redis")
	}
	if result != "OK" {
		return nil, errors.New("result from redis is not OK: " + result)
	}
	return &sessionID, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *session.SessionID) (*session.Session, error) {
	log.Println("check session ", in.ID)
	mkey := "sessions:" + in.ID
	data, err := redis.Bytes(sm.sessions.Do("GET", mkey))
	if err != nil {
		return nil, errors.Wrap(err, "can't get data")
	}
	sess := &session.Session{}
	if err = sess.UnmarshalJSON(data); err != nil {
		return nil, errors.Wrap(err, "can't unpack session-server data")
	}
	return sess, nil
}

func (sm *SessionManager) Delete(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	log.Println("delete session")
	mkey := "sessions:" + in.ID
	_, err := redis.Int(sm.sessions.Do("DEL", mkey))
	if err != nil {
		return nil, errors.Wrap(err, "can't delete session")
	}
	return &session.Nothing{Dummy: true}, nil
}
