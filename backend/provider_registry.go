package backend

import (
 "fmt"
 "net/netip"
 "sync"
)

type ProviderClass string
const (
 ProviderDNS ProviderClass = "dns"
 ProviderCache ProviderClass = "cache"
 ProviderTransit ProviderClass = "transit"
 ProviderCDN ProviderClass = "cdn"
)

type ProviderEndpoint struct { ID string; Provider string; Class ProviderClass; Country string; Address netip.Addr; ASN uint32; Official bool; Authorized bool; Enabled bool; Notes string }

type ProviderRegistry struct { mu sync.RWMutex; items map[string]ProviderEndpoint }
func NewProviderRegistry() *ProviderRegistry { return &ProviderRegistry{items: map[string]ProviderEndpoint{}} }
func (r *ProviderRegistry) Register(p ProviderEndpoint) error { if p.ID=="" || p.Provider=="" || !p.Address.IsValid() { return fmt.Errorf("invalid provider endpoint") }; if p.ASN==0 { return fmt.Errorf("asn required") }; if !p.Authorized { return fmt.Errorf("endpoint is not authorized") }; r.mu.Lock(); defer r.mu.Unlock(); r.items[p.ID]=p; return nil }
func (r *ProviderRegistry) List(class ProviderClass, country string) []ProviderEndpoint { r.mu.RLock(); defer r.mu.RUnlock(); out:=make([]ProviderEndpoint,0); for _,p:=range r.items { if class!="" && p.Class!=class {continue}; if country!="" && p.Country!=country {continue}; out=append(out,p) }; return out }
