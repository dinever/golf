package Golf

import (
	"fmt"
	"net/http"
)

// Session is an interface for session instance, a session instance contains
// data needed.
type Session interface {
	Set(key string, value interface{}) error
	Get(key string) interface{}
	Delete(key interface{}) error
	SessionID() string
}

// MemorySession is an memory based implementaion of Session.
type MemorySession struct {
	sid  string
	data map[string]interface{}
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
	return nil, fmt.Errorf("key %q in session (id %d) not found", key, s.sid)
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
