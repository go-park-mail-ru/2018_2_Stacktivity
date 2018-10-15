package session

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Session struct {
	ID int
}

type SessionID struct {
	ID uuid.UUID
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[SessionID]*Session
}

type SessionManagerI interface {
	Create(*Session) (*SessionID, error)
	Check(*SessionID) (bool, *Session)
	Delete(*SessionID)
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu: sync.RWMutex{},
		// TODO save sessions into key-value storage
		sessions: map[SessionID]*Session{},
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	sm.mu.Lock()
	ID, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "Can't create session ID")
	}
	sessionID := SessionID{ID}
	sm.mu.Unlock()
	sm.sessions[sessionID] = in
	return &sessionID, nil
}

func (sm *SessionManager) Check(in *SessionID) (bool, *Session) {
	sm.mu.RLock()
	session := sm.sessions[*in]
	sm.mu.RUnlock()
	return session != nil, session
}

func (sm *SessionManager) Delete(in *SessionID) {
	sm.mu.Lock()
	delete(sm.sessions, *in)
	sm.mu.Unlock()
}
