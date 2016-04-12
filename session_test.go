package Golf

import (
	"testing"
)

func TestMemorySession(t *testing.T) {
	cases := []struct {
		key, value string
	}{
		{"foo", "bar"},
		{"abc", "123"},
	}

	mgr := NewMemorySessionManager()
	sid, err := mgr.sessionID()
	if err != nil {
		t.Errorf("Could not generate session ID.")
	}
	s := MemorySession{sid: sid, data: make(map[string]interface{})}
	for _, c := range cases {
		s.Set(c.key, c.value)
		value, err := s.Get(c.key)
		if value != c.value {
			t.Errorf("Could not set memory sesseion k-v pair.")
		}

		err = s.Delete(c.key)
		if err != nil {
			t.Errorf("Could not delete session key.")
		}
		_, err = s.Get(c.key)
		if err == nil {
			t.Errorf("Could not correctly delete session key.")
		}
	}
}

func TestMemorySessionManager(t *testing.T) {
	mgr := NewMemorySessionManager()
	s, err := mgr.NewSession()
	if err != nil {
		t.Errorf("Could not create a new session.")
	}
	sid := s.SessionID()
	newSession, _ := mgr.Session(sid)
	if newSession.SessionID() != s.SessionID() {
		t.Errorf("Memory session manager could not retrieve a previously generated session.")
	}
}
