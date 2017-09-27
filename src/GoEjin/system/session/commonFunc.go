package session

import (
	"fmt"
	"errors"
)

var providers = make(map[string]Provider)

func Register(name string, provider Provider) error {
	if provider == nil {
		return errors.New("session:Register provider is nil")
	}
	if _,dup := providers[name];dup {
		return errors.New("session:Register called twice for provider" + name)
	}
	providers[name] = provider
	return nil
}

func NewManager(providerName, cookieName string, maxLifetime int64) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session:unknow provider %q", providerName)
	}
	return &Manager{provider:provider, cookieName:cookieName, maxLifetime:maxLifetime},nil
}
