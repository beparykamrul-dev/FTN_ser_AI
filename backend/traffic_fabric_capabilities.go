package backend

import (
    "fmt"
    "strings"
    "sync"
)

type TrafficFabricCapability string
const (
    TrafficECMP TrafficFabricCapability = "ecmp-instant-failover"
    TrafficMACDiscovery TrafficFabricCapability = "mac-fingerprinting-auto-discovery"
    TrafficGeo TrafficFabricCapability = "p2location-geomap"
    TrafficTPM2 TrafficFabricCapability = "tpm2-tools-trousers"
    TrafficLogs TrafficFabricCapability = "fluentbit-vector"
    TrafficMASQUE TrafficFabricCapability = "masque-http3-quic"
    TrafficShadowsocks TrafficFabricCapability = "shadowsocks"
    TrafficXray TrafficFabricCapability = "xray-vless-reality-xtls"
    TrafficGugiGilong TrafficFabricCapability = "gugigilong-adapter"
)

type TrafficFabricProfile struct { ID, Name, NodeID string; Authorized bool; Enabled bool; Capabilities []TrafficFabricCapability }
type TrafficFabricRegistry struct { mu sync.RWMutex; profiles map[string]TrafficFabricProfile }
func NewTrafficFabricRegistry() *TrafficFabricRegistry { return &TrafficFabricRegistry{profiles:map[string]TrafficFabricProfile{}} }
func (r *TrafficFabricRegistry) Upsert(p TrafficFabricProfile) error {
    if p.ID=="" || p.Name=="" || p.NodeID=="" || !p.Authorized { return fmt.Errorf("invalid or unauthorized traffic fabric profile") }
    if len(p.Capabilities)==0 { return fmt.Errorf("traffic fabric capability is required") }
    seen:=map[TrafficFabricCapability]bool{}
    for _,c:=range p.Capabilities { if strings.TrimSpace(string(c))=="" || seen[c] { return fmt.Errorf("invalid or duplicate capability") }; seen[c]=true }
    r.mu.Lock(); defer r.mu.Unlock(); r.profiles[p.ID]=p; return nil
}
func (r *TrafficFabricRegistry) Get(id string)(TrafficFabricProfile,bool){r.mu.RLock();defer r.mu.RUnlock();p,ok:=r.profiles[id];return p,ok}
func (r *TrafficFabricRegistry) Snapshot()[]TrafficFabricProfile{r.mu.RLock();defer r.mu.RUnlock();out:=make([]TrafficFabricProfile,0,len(r.profiles));for _,p:=range r.profiles{out=append(out,p)};return out}
