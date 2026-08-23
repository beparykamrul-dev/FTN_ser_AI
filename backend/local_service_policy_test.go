package backend

import "testing"

func TestSelectEndpointPrefersLocal(t *testing.T) {
    got, err := SelectEndpoint([]ServiceEndpoint{
        {ID:"global-1", Scope:ScopeGlobal, Healthy:true, Authorized:true, CapacityPct:20, LatencyMs:5},
        {ID:"local-1", Scope:ScopeLocal, Healthy:true, Authorized:true, CapacityPct:50, LatencyMs:20},
    })
    if err != nil || got.EndpointID != "local-1" || got.Scope != ScopeLocal { t.Fatalf("unexpected selection: %#v %v", got, err) }
}

func TestSelectEndpointFallsBackToGlobal(t *testing.T) {
    got, err := SelectEndpoint([]ServiceEndpoint{
        {ID:"local-1", Scope:ScopeLocal, Healthy:false, Authorized:true, CapacityPct:10},
        {ID:"global-1", Scope:ScopeGlobal, Healthy:true, Authorized:true, CapacityPct:60, LatencyMs:15},
    })
    if err != nil || got.EndpointID != "global-1" { t.Fatalf("unexpected fallback: %#v %v", got, err) }
}
