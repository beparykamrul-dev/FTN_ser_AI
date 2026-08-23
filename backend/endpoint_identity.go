package backend

import (
    "fmt"
    "net"
    "strings"
    "sync"
)

type EndpointIdentity struct {
    ID string
    Name string
    IP string
    MAC string
    ASN uint32
    Country string
    Authorized bool
}

type EndpointIdentityStore struct { mu sync.RWMutex; values map[string]EndpointIdentity }

func NewEndpointIdentityStore() *EndpointIdentityStore { return &EndpointIdentityStore{values: map[string]EndpointIdentity{}} }

func validateMAC(v string) bool { if v == "" { return false }; _, err := net.ParseMAC(v); return err == nil }

func (s *EndpointIdentityStore) Add(e EndpointIdentity) error {
    if e.ID == "" || e.Name == "" || net.ParseIP(e.IP) == nil || !validateMAC(e.MAC) { return fmt.Errorf("invalid endpoint identity") }
    if e.ASN == 0 || strings.TrimSpace(e.Country) == "" { return fmt.Errorf("provider identity incomplete") }
    s.mu.Lock(); defer s.mu.Unlock(); s.values[e.ID] = e; return nil
}

func (s *EndpointIdentityStore) Get(id string) (EndpointIdentity, bool) { s.mu.RLock(); defer s.mu.RUnlock(); e, ok := s.values[id]; return e, ok }

func (s *EndpointIdentityStore) Snapshot() []EndpointIdentity { s.mu.RLock(); defer s.mu.RUnlock(); out := make([]EndpointIdentity,0,len(s.values)); for _, e := range s.values { out = append(out,e) }; return out }
