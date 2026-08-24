package backend

import (
    "fmt"
    "sync"
)

type FlowPipelineCapability string
const (
    FlowAetherCore FlowPipelineCapability = "aether-core"
    FlowSiLK FlowPipelineCapability = "cmusei/silk"
    FlowRWWFlowpack FlowPipelineCapability = "rwflowpack"
)

type FlowFabricProfile struct {
    ID string
    NodeID string
    Capabilities []FlowPipelineCapability
    Authorized bool
    Enabled bool
}

type FlowFabricRegistry struct { mu sync.RWMutex; profiles map[string]FlowFabricProfile }
func NewFlowFabricRegistry() *FlowFabricRegistry { return &FlowFabricRegistry{profiles:map[string]FlowFabricProfile{}} }

func (r *FlowFabricRegistry) Upsert(p FlowFabricProfile) error {
    if p.ID=="" || p.NodeID=="" || !p.Authorized { return fmt.Errorf("invalid or unauthorized flow fabric profile") }
    if len(p.Capabilities)==0 { return fmt.Errorf("flow capability is required") }
    seen:=map[FlowPipelineCapability]bool{}
    for _, c:=range p.Capabilities { if c!=FlowAetherCore && c!=FlowSiLK && c!=FlowRWWFlowpack { return fmt.Errorf("unsupported flow capability") }; if seen[c] { return fmt.Errorf("duplicate flow capability") }; seen[c]=true }
    r.mu.Lock(); r.profiles[p.ID]=p; r.mu.Unlock(); return nil
}
func (r *FlowFabricRegistry) Get(id string)(FlowFabricProfile,bool){r.mu.RLock();defer r.mu.RUnlock();p,ok:=r.profiles[id];return p,ok}
