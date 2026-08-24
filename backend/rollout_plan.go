package backend

import (
    "fmt"
    "sync"
    "time"
)

type RolloutStage string
const (
    RolloutValidate RolloutStage = "validate"
    RolloutCanary RolloutStage = "canary"
    RolloutLive RolloutStage = "live"
    RolloutRollback RolloutStage = "rollback"
)

type ServiceRollout struct { ID, ServiceID string; Stage RolloutStage; Approved bool; Ready bool; UpdatedAt time.Time }

type RolloutRegistry struct { mu sync.RWMutex; items map[string]ServiceRollout }
func NewRolloutRegistry() *RolloutRegistry { return &RolloutRegistry{items:map[string]ServiceRollout{}} }

func (r *RolloutRegistry) Advance(x ServiceRollout, now time.Time) error {
    if x.ID=="" || x.ServiceID=="" || now.IsZero() { return fmt.Errorf("invalid rollout") }
    if !x.Approved { return fmt.Errorf("rollout requires approval") }
    if x.Stage==RolloutLive && !x.Ready { return fmt.Errorf("service is not ready") }
    if x.Stage!=RolloutValidate && x.Stage!=RolloutCanary && x.Stage!=RolloutLive && x.Stage!=RolloutRollback { return fmt.Errorf("unsupported rollout stage") }
    x.UpdatedAt=now; r.mu.Lock(); r.items[x.ID]=x; r.mu.Unlock(); return nil
}
func (r *RolloutRegistry) Get(id string)(ServiceRollout,bool){r.mu.RLock();defer r.mu.RUnlock();x,ok:=r.items[id];return x,ok}
