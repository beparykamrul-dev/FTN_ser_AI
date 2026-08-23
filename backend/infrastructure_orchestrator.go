package backend

import (
    "fmt"
    "sort"
    "sync"
    "time"
)

type NodeCapability struct { NodeID string; CPUCapacityPct float64; MemoryCapacityPct float64; StorageCapacityPct float64; NetworkCapacityPct float64; Services []string; Healthy bool; Authorized bool }

type PlacementPlan struct { ServiceID string; NodeID string; Score float64; Reason string; CreatedAt time.Time }

type InfrastructureOrchestrator struct { mu sync.RWMutex; nodes map[string]NodeCapability; plans map[string]PlacementPlan }

func NewInfrastructureOrchestrator() *InfrastructureOrchestrator { return &InfrastructureOrchestrator{nodes:map[string]NodeCapability{}, plans:map[string]PlacementPlan{}} }

func (o *InfrastructureOrchestrator) RegisterNode(n NodeCapability) error {
    if n.NodeID == "" || !n.Authorized { return fmt.Errorf("node identity or authorization required") }
    for _, v := range []float64{n.CPUCapacityPct,n.MemoryCapacityPct,n.StorageCapacityPct,n.NetworkCapacityPct} { if v < 0 || v > 100 { return fmt.Errorf("capacity must be between 0 and 100") } }
    o.mu.Lock(); defer o.mu.Unlock(); o.nodes[n.NodeID]=n; return nil
}

func (o *InfrastructureOrchestrator) PlanPlacement(serviceID string, now time.Time) (PlacementPlan, error) {
    if serviceID == "" || now.IsZero() { return PlacementPlan{}, fmt.Errorf("service id and timestamp required") }
    o.mu.RLock(); defer o.mu.RUnlock()
    candidates := make([]NodeCapability,0)
    for _, n := range o.nodes { if n.Healthy && n.Authorized { candidates=append(candidates,n) } }
    if len(candidates)==0 { return PlacementPlan{}, fmt.Errorf("no eligible nodes") }
    sort.Slice(candidates,func(i,j int) bool { return capacityScore(candidates[i]) < capacityScore(candidates[j]) })
    n:=candidates[0]
    p:=PlacementPlan{ServiceID:serviceID,NodeID:n.NodeID,Score:capacityScore(n),Reason:"healthy authorized node with highest available capacity",CreatedAt:now}
    return p,nil
}

func capacityScore(n NodeCapability) float64 { return (n.CPUCapacityPct+n.MemoryCapacityPct+n.StorageCapacityPct+n.NetworkCapacityPct)/4 }
