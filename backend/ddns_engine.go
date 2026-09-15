package backend

import (
	"errors"
	"net/netip"
	"strings"
	"sync"
	"time"
)

type DDNSRecord struct {
	Hostname  string    `json:"hostname"`
	Address   string    `json:"address"`
	TTL       uint32    `json:"ttl"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DDNSStore struct { mu sync.RWMutex; records map[string]DDNSRecord }

func NewDDNSStore() *DDNSStore { return &DDNSStore{records: map[string]DDNSRecord{}} }

func (s *DDNSStore) Upsert(hostname, address string, ttl uint32) (DDNSRecord, error) {
	hostname = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(hostname), "."))
	if hostname == "" || !strings.Contains(hostname, ".") { return DDNSRecord{}, errors.New("invalid hostname") }
	if _, err := netip.ParseAddr(strings.TrimSpace(address)); err != nil { return DDNSRecord{}, errors.New("invalid IP address") }
	if ttl == 0 { ttl = 60 }
	if ttl < 30 { return DDNSRecord{}, errors.New("ttl must be at least 30 seconds") }
	r := DDNSRecord{Hostname: hostname, Address: strings.TrimSpace(address), TTL: ttl, UpdatedAt: time.Now().UTC()}
	s.mu.Lock(); s.records[hostname] = r; s.mu.Unlock()
	return r, nil
}

func (s *DDNSStore) Get(hostname string) (DDNSRecord, bool) {
	hostname = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(hostname), "."))
	s.mu.RLock(); defer s.mu.RUnlock(); r, ok := s.records[hostname]; return r, ok
}
