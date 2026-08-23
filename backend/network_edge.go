package backend

import (
    "fmt"
    "strings"
    "sync"
)

type NetworkEdgeCapability string
const (
    EdgeNginx NetworkEdgeCapability = "nginx"
    EdgeDynamicMesh NetworkEdgeCapability = "dynamic-mesh"
    EdgeDNSStealth NetworkEdgeCapability = "dns-stealth"
    EdgeEBPFXDP NetworkEdgeCapability = "ebpf-xdp"
    EdgeVPNZeroTrust NetworkEdgeCapability = "vpn-zero-trust-gateway"
)

type NetworkEdgeProfile struct {
    ID string
    Name string
    NodeID string
    Capabilities []NetworkEdgeCapability
    Authorized bool
    Enabled bool
}

type NetworkEdgeRegistry struct { mu sync.RWMutex; profiles map[string]NetworkEdgeProfile }
func NewNetworkEdgeRegistry() *NetworkEdgeRegistry { return &NetworkEdgeRegistry{profiles: map[string]NetworkEdgeProfile{}} }

func (r *NetworkEdgeRegistry) Upsert(p NetworkEdgeProfile) error {
    if p.ID == "" || p.Name == "" || p.NodeID == "" || !p.Authorized { return fmt.Errorf("invalid or unauthorized network edge profile") }
    if len(p.Capabilities) == 0 { return fmt.Errorf("at least one edge capability is required") }
    seen := map[NetworkEdgeCapability]bool{}
    for _, c := range p.Capabilities { if strings.TrimSpace(string(c)) == "" || seen[c] { return fmt.Errorf("invalid or duplicate capability") }; seen[c] = true }
    r.mu.Lock(); defer r.mu.Unlock(); r.profiles[p.ID] = p; return nil
}

func (r *NetworkEdgeRegistry) Get(id string) (NetworkEdgeProfile, bool) { r.mu.RLock(); defer r.mu.RUnlock(); p, ok := r.profiles[id]; return p, ok }
func (r *NetworkEdgeRegistry) Snapshot() []NetworkEdgeProfile { r.mu.RLock(); defer r.mu.RUnlock(); out:=make([]NetworkEdgeProfile,0,len(r.profiles)); for _, p:=range r.profiles { out=append(out,p) }; return out }
