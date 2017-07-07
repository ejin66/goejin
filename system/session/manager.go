package session

import (
	"sync"
	"io"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"
)

func init() {

}

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifetime int64
}

func (this *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (this *Manager) SessionStart(w *http.ResponseWriter, r *http.Request) (session Session) {
	this.lock.Lock()
	defer this.lock.Unlock()
	cookie, err := r.Cookie(this.cookieName)

	//若cookie有值，尝试取出session。 若没有cookie或者没有取到session, 则重新创建
	if err == nil && cookie.Value != "" {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, err = this.provider.SessionRead(sid)
		if err == nil {
			return
		}
	}

	if err != nil || cookie.Value == "" {
		sid := this.sessionId()
		session, _ = this.provider.SessionInit(sid)
		//TODO 等下细读
		cookieNew := http.Cookie{Name: this.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(this.maxLifetime)}
		http.SetCookie(*w, &cookieNew)
	}
	return
}

func (this *Manager) GC() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.provider.SessionGC(this.maxLifetime)
	time.AfterFunc(time.Duration(this.maxLifetime) * time.Second,this.GC)
}
