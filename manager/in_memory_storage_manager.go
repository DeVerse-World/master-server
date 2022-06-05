package manager

import (
	"time"

	cache "github.com/Shopify/go-cache"
)

type InMemoryStorageManager struct {
	client          cache.Client
	timeoutDuration time.Duration
}

func NewInMemoryStorageManager() *InMemoryStorageManager {
	return &InMemoryStorageManager{
		client:          cache.NewMemoryClient(),
		timeoutDuration: 3000 * time.Second,
	}
}

func (m *InMemoryStorageManager) Get(key string, data interface{}) error {
	return m.client.Get(key, data)
}

func (m *InMemoryStorageManager) Set(key string, data interface{}) error {
	return m.client.Set(key, data, time.Now().Add(m.timeoutDuration))
}

func (m *InMemoryStorageManager) Delete(key string) error {
	return m.client.Delete(key)
}
