package idempotency

import "sync"

type Store struct {
	mu    sync.RWMutex
	items map[string]any
}

func NewStore() *Store {
	return &Store{items: make(map[string]any)}
}

func (s *Store) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.items[key]
	return value, ok
}

func (s *Store) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
}
