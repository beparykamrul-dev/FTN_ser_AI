package backend

import (
    "fmt"
    "sync"
    "time"
)

type HealthStatus string

const (
    HealthUnknown HealthStatus = "unknown"
    HealthHealthy HealthStatus = "healthy"
    HealthDegraded HealthStatus = "degraded"
    HealthUnhealthy HealthStatus = "unhealthy"
)

type HealthObservation struct {
    NodeID string
    ServiceID string
    Status HealthStatus
    LatencyMs float64
    PacketLossPct float64
    At time.Time
}

type HealthPolicy struct {
    MaxLatencyMs float64
    MaxPacketLossPct float64
}

func EvaluateHealth(o HealthObservation, p HealthPolicy) (HealthStatus, error) {
    if o.NodeID == "" || o.ServiceID == "" || o.At.IsZero() { return HealthUnknown, fmt.Errorf("invalid observation") }
    if o.LatencyMs < 0 || o.PacketLossPct < 0 || o.PacketLossPct > 100 { return HealthUnknown, fmt.Errorf("invalid network metrics") }
    if p.MaxLatencyMs < 0 || p.MaxPacketLossPct < 0 { return HealthUnknown, fmt.Errorf("invalid health policy") }
    if o.PacketLossPct > p.MaxPacketLossPct || o.LatencyMs > p.MaxLatencyMs { return HealthDegraded, nil }
    return HealthHealthy, nil
}

type HealthStore struct { mu sync.RWMutex; values map[string]HealthObservation }

func NewHealthStore() *HealthStore { return &HealthStore{values: map[string]HealthObservation{}} }
func (s *HealthStore) Put(o HealthObservation) error { if o.NodeID == "" || o.ServiceID == "" { return fmt.Errorf("identity required") }; s.mu.Lock(); defer s.mu.Unlock(); s.values[o.NodeID+":"+o.ServiceID] = o; return nil }
func (s *HealthStore) Get(nodeID, serviceID string) (HealthObservation, bool) { s.mu.RLock(); defer s.mu.RUnlock(); o, ok := s.values[nodeID+":"+serviceID]; return o, ok }
