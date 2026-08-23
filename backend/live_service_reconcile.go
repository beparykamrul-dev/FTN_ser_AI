package backend

import (
    "fmt"
    "sync"
    "time"
)

type RuntimeAction struct { ID string; ServiceID string; Desired ServiceRuntimeState; Applied bool; At time.Time }

type ServiceReconciler struct { mu sync.Mutex; registry *LiveServiceRegistry; actions map[string]RuntimeAction }

func NewServiceReconciler(r *LiveServiceRegistry) (*ServiceReconciler, error) { if r == nil { return nil, fmt.Errorf("registry is required") }; return &ServiceReconciler{registry:r, actions:map[string]RuntimeAction{}}, nil }

func (r *ServiceReconciler) Reconcile(a RuntimeAction) error {
    if a.ID == "" || a.ServiceID == "" || a.At.IsZero() { return fmt.Errorf("invalid runtime action") }
    r.mu.Lock(); defer r.mu.Unlock()
    if _, exists := r.actions[a.ID]; exists { return fmt.Errorf("duplicate runtime action") }
    s, ok := r.registry.Get(a.ServiceID); if !ok { return fmt.Errorf("service not found") }
    if !s.Authorized { return fmt.Errorf("service is not authorized") }
    if a.Desired != ServiceLive && a.Desired != ServiceDegraded && a.Desired != ServiceStopped { return fmt.Errorf("unsupported desired state") }
    s.State = a.Desired; s.LastHealthCheck = a.At
    r.registry.mu.Lock(); r.registry.services[s.ID] = s; r.registry.mu.Unlock()
    a.Applied = true; r.actions[a.ID] = a
    return nil
}
