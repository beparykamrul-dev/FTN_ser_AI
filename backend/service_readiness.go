package backend

import (
    "fmt"
    "sync"
    "time"
)

type ReadinessCheck struct { Name string; Passed bool; Detail string }
type ServiceReadiness struct { ServiceID string; Ready bool; Checks []ReadinessCheck; CheckedAt time.Time }

type ServiceReadinessEngine struct { mu sync.RWMutex; results map[string]ServiceReadiness }
func NewServiceReadinessEngine() *ServiceReadinessEngine { return &ServiceReadinessEngine{results: map[string]ServiceReadiness{}} }

func (e *ServiceReadinessEngine) Evaluate(s LiveService, now time.Time) (ServiceReadiness, error) {
    if s.ID == "" || s.Name == "" || s.NodeID == "" { return ServiceReadiness{}, fmt.Errorf("incomplete service identity") }
    checks := []ReadinessCheck{
        {Name:"authorization", Passed:s.Authorized, Detail:"service authorization"},
        {Name:"runtime-state", Passed:s.State == ServiceLive, Detail:"service is live"},
        {Name:"health-timestamp", Passed:!s.LastHealthCheck.IsZero(), Detail:"health check recorded"},
    }
    ready := true; for _, c := range checks { if !c.Passed { ready = false; break } }
    result := ServiceReadiness{ServiceID:s.ID, Ready:ready, Checks:checks, CheckedAt:now}
    e.mu.Lock(); e.results[s.ID] = result; e.mu.Unlock(); return result, nil
}

func (e *ServiceReadinessEngine) Get(id string) (ServiceReadiness, bool) { e.mu.RLock(); defer e.mu.RUnlock(); r, ok := e.results[id]; return r, ok }
