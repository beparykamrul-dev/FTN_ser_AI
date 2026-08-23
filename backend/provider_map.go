package backend

import (
    "fmt"
    "sync"
)

type ProviderMapPoint struct {
    ID string
    Provider string
    IP string
    ASN uint32
    Country string
    ServiceType string
    Official bool
    Authorized bool
    Latitude float64
    Longitude float64
    Healthy bool
}

type ProviderMap struct { mu sync.RWMutex; points map[string]ProviderMapPoint }

func NewProviderMap() *ProviderMap { return &ProviderMap{points: map[string]ProviderMapPoint{}} }

func (m *ProviderMap) Upsert(p ProviderMapPoint) error {
    if p.ID == "" || p.Provider == "" || p.IP == "" || p.ASN == 0 || p.Country == "" || p.ServiceType == "" { return fmt.Errorf("provider map identity incomplete") }
    if p.Latitude < -90 || p.Latitude > 90 || p.Longitude < -180 || p.Longitude > 180 { return fmt.Errorf("invalid coordinates") }
    m.mu.Lock(); defer m.mu.Unlock(); m.points[p.ID] = p; return nil
}

func (m *ProviderMap) Snapshot() []ProviderMapPoint {
    m.mu.RLock(); defer m.mu.RUnlock()
    out := make([]ProviderMapPoint, 0, len(m.points))
    for _, p := range m.points { out = append(out, p) }
    return out
}

func (m *ProviderMap) Filter(country, serviceType string, authorizedOnly bool) []ProviderMapPoint {
    m.mu.RLock(); defer m.mu.RUnlock()
    out := make([]ProviderMapPoint, 0)
    for _, p := range m.points {
        if country != "" && p.Country != country { continue }
        if serviceType != "" && p.ServiceType != serviceType { continue }
        if authorizedOnly && !p.Authorized { continue }
        out = append(out, p)
    }
    return out
}
