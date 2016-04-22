package golf

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

const sessionIDLength = 64
const sessionExpireTime = 3600
const gcTimeInterval = 30

// SessionManager manages a map of sessions.
type SessionManager interface {
	sessionID() (string, error)
	NewSession() (Session, error)
	Session(string) (Session, error)
	GarbageCollection()
	Count() int
}

// MemorySessionManager is a implementation of Session Manager, which stores data in memory.
type MemorySessionManager struct {
	sessions map[string]*MemorySession
	lock     sync.RWMutex
}

// NewMemorySessionManager creates a new session manager.
func NewMemorySessionManager() *MemorySessionManager {
	mgr := new(MemorySessionManager)
	mgr.sessions = make(map[string]*MemorySession)
	mgr.GarbageCollection()
	return mgr
}

func (mgr *MemorySessionManager) sessionID() (string, error) {
	b := make([]byte, sessionIDLength)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG.")
	}
	return hex.EncodeToString(b), nil
}

// Session gets the session instance by indicating a session id.
func (mgr *MemorySessionManager) Session(sid string) (Session, error) {
	mgr.lock.RLock()
	if s, ok := mgr.sessions[sid]; ok {
		s.createdAt = time.Now()
		mgr.lock.RUnlock()
		return s, nil
	}
	mgr.lock.RUnlock()
	return nil, fmt.Errorf("Can not retrieve session with id %s.", sid)
}

// NewSession creates and returns a new session.
func (mgr *MemorySessionManager) NewSession() (Session, error) {
	sid, err := mgr.sessionID()
	if err != nil {
		return nil, err
	}
	s := MemorySession{sid: sid, data: make(map[string]interface{}), createdAt: time.Now()}
	mgr.lock.Lock()
	mgr.sessions[sid] = &s
	mgr.lock.Unlock()
	return &s, nil
}

// GarbageCollection recycles expired sessions, delete them from the session manager.
func (mgr *MemorySessionManager) GarbageCollection() {
	for k, v := range mgr.sessions {
		if v.isExpired() {
			delete(mgr.sessions, k)
		}
	}
	time.AfterFunc(time.Duration(gcTimeInterval)*time.Second, mgr.GarbageCollection)
}

// Count returns the number of the current session stored in the session manager.
func (mgr *MemorySessionManager) Count() int {
	return len(mgr.sessions)
}

// Session is an interface for session instance, a session instance contains
// data needed.
type Session interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	SessionID() string
	isExpired() bool
}

// MemorySession is an memory based implementaion of Session.
type MemorySession struct {
	sid       string
	data      map[string]interface{}
	createdAt time.Time
}

func (s *MemorySession) isExpired() bool {
	return (s.createdAt.Unix() + sessionExpireTime) <= time.Now().Unix()
}

// Set method sets the key value pair in the session.
func (s *MemorySession) Set(key string, value interface{}) error {
	s.data[key] = value
	return nil
}

// Get method gets the value by given a key in the session.
func (s *MemorySession) Get(key string) (interface{}, error) {
	if value, ok := s.data[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("key %q in session (id %s) not found", key, s.sid)
}

// Delete method deletes the value by given a key in the session.
func (s *MemorySession) Delete(key string) error {
	delete(s.data, key)
	return nil
}

// SessionID returns the current ID of the session.
func (s *MemorySession) SessionID() string {
	return s.sid
}
