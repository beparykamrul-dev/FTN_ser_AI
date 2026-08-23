package backend

import "fmt"

type ServiceScope string

const (
    ScopeLocal  ServiceScope = "local"
    ScopeGlobal ServiceScope = "global"
)

type ServiceEndpoint struct {
    ID string
    Scope ServiceScope
    Healthy bool
    CapacityPct float64
    LatencyMs float64
    Authorized bool
}

type ServiceSelection struct { EndpointID string; Scope ServiceScope; Reason string }

// SelectEndpoint prefers healthy authorized local endpoints. Global endpoints are fallback capacity.
func SelectEndpoint(endpoints []ServiceEndpoint) (ServiceSelection, error) {
    var local, global *ServiceEndpoint
    for i := range endpoints {
        e := endpoints[i]
        if !e.Authorized || !e.Healthy || e.CapacityPct < 0 || e.CapacityPct > 100 || e.LatencyMs < 0 { continue }
        if e.Scope == ScopeLocal {
            if local == nil || e.CapacityPct < local.CapacityPct || (e.CapacityPct == local.CapacityPct && e.LatencyMs < local.LatencyMs) { local = &e }
        } else if e.Scope == ScopeGlobal {
            if global == nil || e.LatencyMs < global.LatencyMs { global = &e }
        }
    }
    if local != nil { return ServiceSelection{EndpointID: local.ID, Scope: ScopeLocal, Reason: "healthy authorized local endpoint preferred"}, nil }
    if global != nil { return ServiceSelection{EndpointID: global.ID, Scope: ScopeGlobal, Reason: "local endpoint unavailable; global fallback selected"}, nil }
    return ServiceSelection{}, fmt.Errorf("no eligible endpoint")
}
