package backend

import "sync"

type ProviderStore struct { mu sync.RWMutex; providers map[string]DNSProvider }
func NewProviderStore() *ProviderStore { return &ProviderStore{providers:map[string]DNSProvider{}} }
func (s *ProviderStore) Upsert(p DNSProvider){s.mu.Lock();s.providers[p.ID]=p;s.mu.Unlock()}
func (s *ProviderStore) List() []DNSProvider{s.mu.RLock();defer s.mu.RUnlock();out:=make([]DNSProvider,0,len(s.providers));for _,p:=range s.providers{out=append(out,p)};return out}
