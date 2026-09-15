package backend

import "sync"

type DBStore struct { mu sync.RWMutex; health DBHealth }
func NewDBStore() *DBStore { return &DBStore{health:DBHealth{Instances:[]DBInstance{}}} }
func (s *DBStore) Get() DBHealth { s.mu.RLock(); defer s.mu.RUnlock(); return s.health }
func (s *DBStore) Set(v DBHealth) { s.mu.Lock(); s.health=v; s.mu.Unlock() }
