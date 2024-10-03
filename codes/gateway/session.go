package main

import (
	// "net"
	"fmt"
	"net"
	"sync"

	"github.com/xtaci/smux"
)

type Session struct {
	ClientID string
	// Connection net.Conn
	Connection *smux.Session
}

type SessionManager struct {
	sessionsMu sync.Mutex
	sessions   map[string]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (mgr *SessionManager) GetSessionByClientID(clientID string) (net.Conn, error) {
	mgr.sessionsMu.Lock()
	defer mgr.sessionsMu.Unlock()
	sess := mgr.sessions[clientID]
	if sess == nil {
		return nil, fmt.Errorf("client id %s not connected", clientID)
	}

	stream, err := sess.Connection.OpenStream()

	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (sm *SessionManager) CreateSession(clientID string, conn net.Conn) (*Session, error) {
	sm.sessionsMu.Lock()
	defer sm.sessionsMu.Unlock()

	mux, err := smux.Server(conn, nil)
	if err != nil {
		return nil, err
	}

	old := sm.sessions[clientID]
	if old != nil {
		return nil, fmt.Errorf("client id %s already connected", clientID)
	}

	sess := &Session{
		ClientID:   clientID,
		Connection: mux,
	}
	sm.sessions[clientID] = sess
	return sess, nil
}

func (sm *SessionManager) CloseSession(clientID string) {
	sm.sessionsMu.Lock()
	defer sm.sessionsMu.Unlock()
	sess := sm.sessions[clientID]
	if sess == nil {
		return
	}

	sess.Connection.Close()
	delete(sm.sessions, clientID)
}

func (sm *SessionManager) Range(f func(k string, v *Session) bool) {
	sm.sessionsMu.Lock()
	defer sm.sessionsMu.Unlock()

	for k, v := range sm.sessions {
		ok := f(k, v)
		if !ok {
			delete(sm.sessions, k)
		}
	}
}
