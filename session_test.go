package golf

import (
	"testing"
	"time"
)

func TestMemorySessionCRUD(t *testing.T) {
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

func TestMemorySessionExpire(t *testing.T) {
	mgr := NewMemorySessionManager()
	sid, _ := mgr.sessionID()
	s := MemorySession{sid: sid, data: make(map[string]interface{}), createdAt: time.Now().AddDate(0, 0, -1)}
	mgr.sessions[sid] = &s
	mgr.GarbageCollection()
	_, err := mgr.Session(sid)
	if err == nil {
		t.Errorf("Could not correctly recycle expired sessions.")
	}
}

func TestMemorySessionNotExpire(t *testing.T) {
	mgr := NewMemorySessionManager()
	sid, _ := mgr.sessionID()
	s := MemorySession{sid: sid, data: make(map[string]interface{}), createdAt: time.Now()}
	mgr.sessions[sid] = &s
	mgr.GarbageCollection()
	_, err := mgr.Session(sid)
	if err != nil {
		t.Errorf("Falsely recycled non-expired sessions.")
	}
}

func TestMemorySessionCount(t *testing.T) {
	mgr := NewMemorySessionManager()
	for i := 0; i < 100; i++ {
		sid, _ := mgr.sessionID()
		s := MemorySession{sid: sid, data: make(map[string]interface{}), createdAt: time.Now().AddDate(0, 0, -1)}
		mgr.sessions[sid] = &s
	}
	if mgr.Count() != 100 {
		t.Errorf("Could not correctly get session count: %v != %v", mgr.Count(), 100)
	}
	mgr.GarbageCollection()
	if mgr.Count() != 0 {
		t.Errorf("Could not correctly get session count after GC: %v != %v", mgr.Count(), 0)
	}
}
