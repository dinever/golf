package Golf

type SessionManager struct {
	data map[string]interface{}
}

func (session *SessionManager) Get(key string) interface{} {
	return 5
}
