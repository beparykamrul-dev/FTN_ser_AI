package backend

import (
    "context"
    "errors"
    "sync"
    "time"
)

type AdapterKind string

const (
    AdapterGitHub AdapterKind = "github"
    AdapterCloudflare AdapterKind = "cloudflare"
    AdapterGoogle AdapterKind = "google"
    AdapterRender AdapterKind = "render"
    AdapterFTN AdapterKind = "ftn"
    AdapterAgent AdapterKind = "agent"
)

type Adapter interface {
    ID() string
    Kind() AdapterKind
    Enabled() bool
    Health(context.Context) error
}

type AdapterStatus struct {
    ID string `json:"id"`
    Kind AdapterKind `json:"kind"`
    Enabled bool `json:"enabled"`
    Healthy bool `json:"healthy"`
    LastCheck time.Time `json:"last_check,omitempty"`
    Error string `json:"error,omitempty"`
}

type AdapterEngine struct {
    mu sync.RWMutex
    adapters map[string]Adapter
    status map[string]AdapterStatus
}

func NewAdapterEngine() *AdapterEngine {
    return &AdapterEngine{adapters: map[string]Adapter{}, status: map[string]AdapterStatus{}}
}

func (e *AdapterEngine) Register(a Adapter) error {
    if a == nil || a.ID() == "" { return errors.New("invalid adapter") }
    e.mu.Lock(); defer e.mu.Unlock()
    e.adapters[a.ID()] = a
    e.status[a.ID()] = AdapterStatus{ID:a.ID(), Kind:a.Kind(), Enabled:a.Enabled()}
    return nil
}

func (e *AdapterEngine) Check(ctx context.Context, id string) AdapterStatus {
    e.mu.RLock(); a := e.adapters[id]; e.mu.RUnlock()
    if a == nil { return AdapterStatus{ID:id, Error:"adapter not found"} }
    st := AdapterStatus{ID:a.ID(), Kind:a.Kind(), Enabled:a.Enabled(), LastCheck:time.Now().UTC()}
    if !a.Enabled() { st.Error="adapter disabled"; return st }
    if err:=a.Health(ctx); err!=nil { st.Error=err.Error() } else { st.Healthy=true }
    e.mu.Lock(); e.status[id]=st; e.mu.Unlock()
    return st
}

// ExecuteHealth exposes only the typed provider health operation to the control plane.
// Credentials remain inside the adapter and are never returned to the web client.
func (e *AdapterEngine) ExecuteHealth(ctx context.Context, id string) error {
    e.mu.RLock(); a := e.adapters[id]; e.mu.RUnlock()
    if a == nil { return errors.New("adapter not found") }
    if !a.Enabled() { return errors.New("adapter disabled or not configured") }
    return a.Health(ctx)
}

func (e *AdapterEngine) List() []AdapterStatus {
    e.mu.RLock(); defer e.mu.RUnlock()
    out:=make([]AdapterStatus,0,len(e.status)); for _,s:=range e.status { out=append(out,s) }; return out
}
