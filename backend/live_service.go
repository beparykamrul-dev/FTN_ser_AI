package backend

import (
    "fmt"
    "sync"
    "time"
)

type ServiceRuntimeState string
const (
    ServiceStopped ServiceRuntimeState = "stopped"
    ServiceStarting ServiceRuntimeState = "starting"
    ServiceLive ServiceRuntimeState = "live"
    ServiceDegraded ServiceRuntimeState = "degraded"
    ServiceFailed ServiceRuntimeState = "failed"
)

type LiveService struct {
    ID string
    Name string
    NodeID string
    State ServiceRuntimeState
    LastHealthCheck time.Time
    Authorized bool
}

type LiveServiceRegistry struct { mu sync.RWMutex; services map[string]LiveService }
func NewLiveServiceRegistry() *LiveServiceRegistry { return &LiveServiceRegistry{services: map[string]LiveService{}} }

func (r *LiveServiceRegistry) Activate(s LiveService, now time.Time) error {
    if s.ID == "" || s.Name == "" || s.NodeID == "" || now.IsZero() { return fmt.Errorf("invalid live service identity") }
    if !s.Authorized { return fmt.Errorf("service is not authorized") }
    s.State = ServiceLive
    s.LastHealthCheck = now
    r.mu.Lock(); defer r.mu.Unlock(); r.services[s.ID] = s
    return nil
}

func (r *LiveServiceRegistry) SetHealth(id string, healthy bool, now time.Time) error {
    r.mu.Lock(); defer r.mu.Unlock()
    s, ok := r.services[id]
    if !ok { return fmt.Errorf("service not found") }
    s.LastHealthCheck = now
    if healthy { s.State = ServiceLive } else { s.State = ServiceDegraded }
    r.services[id] = s
    return nil
}

func (r *LiveServiceRegistry) Get(id string) (LiveService, bool) { r.mu.RLock(); defer r.mu.RUnlock(); s, ok := r.services[id]; return s, ok }
