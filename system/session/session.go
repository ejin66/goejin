package session

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(keu interface{})
	SessionId() string
}
