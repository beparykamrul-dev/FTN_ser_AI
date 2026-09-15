package backend

import "sync"

type DNSMeshStore struct { mu sync.RWMutex; status DNSMeshStatus }

func NewDNSMeshStore() *DNSMeshStore {
	return &DNSMeshStore{status: DNSMeshStatus{
		Zones: []string{"ftndns.com", "ftnddns.net"},
		DNSSEC: true,
	}}
}

func (s *DNSMeshStore) Get() DNSMeshStatus { s.mu.RLock(); defer s.mu.RUnlock(); return s.status }
func (s *DNSMeshStore) Set(v DNSMeshStatus) { s.mu.Lock(); s.status=v; s.mu.Unlock() }
