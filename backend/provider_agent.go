package backend

import (
    "fmt"
    "sync"
)

type ProviderAgentMode string
const (
    AgentCatalog ProviderAgentMode = "catalog"
    AgentConnected ProviderAgentMode = "connected"
    AgentPending ProviderAgentMode = "pending"
)

type ProviderAgent struct {
    ID string
    Provider string
    ServiceType string
    Mode ProviderAgentMode
    Authorized bool
    Enabled bool
    Endpoint string
    AccountRef string
}

type ProviderAgentRegistry struct { mu sync.RWMutex; agents map[string]ProviderAgent }

func NewProviderAgentRegistry() *ProviderAgentRegistry { return &ProviderAgentRegistry{agents: map[string]ProviderAgent{}} }

func (r *ProviderAgentRegistry) Register(a ProviderAgent) error {
    if a.ID == "" || a.Provider == "" || a.ServiceType == "" { return fmt.Errorf("provider agent identity incomplete") }
    if a.Mode == "" { a.Mode = AgentCatalog }
    r.mu.Lock(); defer r.mu.Unlock(); r.agents[a.ID] = a; return nil
}

func (r *ProviderAgentRegistry) SetEnabled(id string, enabled bool) error {
    r.mu.Lock(); defer r.mu.Unlock()
    a, ok := r.agents[id]; if !ok { return fmt.Errorf("provider agent not found") }
    if enabled && !a.Authorized { return fmt.Errorf("provider agent is not authorized") }
    a.Enabled = enabled
    if enabled { a.Mode = AgentConnected } else { a.Mode = AgentPending }
    r.agents[id] = a
    return nil
}

func (r *ProviderAgentRegistry) Get(id string) (ProviderAgent, bool) { r.mu.RLock(); defer r.mu.RUnlock(); a, ok := r.agents[id]; return a, ok }
