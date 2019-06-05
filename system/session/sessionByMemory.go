package session

import (
	"time"
	"sync"
	"container/list"
	"errors"
)

var mProvider = &MyProvider{sessions:make(map[string]*list.Element),list:list.New()}
var SessionManager *Manager

func init()  {
	Register("memory", mProvider)
	SessionManager,_  = NewManager("memory", "goejincookie",3600)
	SessionManager.GC()
}

type SessionStore struct {
	sid          string
	timeAccessed time.Time //最后访问时间
	value        map[interface{}]interface{}
}

func (this *SessionStore) Set(key, value interface{}) error {
	this.value[key] = value
	mProvider.SessionUpdate(this.sid)
	return nil
}

func (this *SessionStore) Get(key interface{}) interface{} {
	mProvider.SessionUpdate(this.sid)
	if v, ok := this.value[key]; ok {
		return v
	}
	return nil
}

func (this *SessionStore) Delete(key interface{}) {
	delete(this.value, key)
	mProvider.SessionUpdate(this.sid)
}

func (this *SessionStore) SessionId() string {
	return this.sid
}

type MyProvider struct {
	lock sync.Mutex
	//TODO 细看
	sessions map[string]*list.Element
	list     *list.List //用来做GC??
}

func (this *MyProvider) SessionInit(sid string) (Session, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newSession := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := this.list.PushBack(newSession)
	this.sessions[sid] = element
	return newSession, nil
}

func (this *MyProvider) SessionRead(sid string) (Session, error) {
	if element, ok := this.sessions[sid]; ok {
		s := element.Value.(*SessionStore)

		if (s.timeAccessed.Unix() + SessionManager.maxLifetime < time.Now().Unix()) {
			this.SessionDestroy(sid)
			return nil,errors.New("session:session is expired")
		} else {
			return element.Value.(*SessionStore), nil
		}
	}
	return this.SessionInit(sid)
}

func (this *MyProvider) SessionDestroy(sid string) {
	if element, ok := this.sessions[sid]; ok {
		delete(this.sessions, sid)
		this.list.Remove(element)
	}
}

func (this *MyProvider) SessionGC(maxLifetime int64) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for {
		element := this.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxLifetime) < time.Now().Unix() {
			this.list.Remove(element)
			delete(this.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (this *MyProvider) SessionUpdate(sid string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if element, ok := this.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		this.list.MoveToFront(element)
	}
	return nil
}
