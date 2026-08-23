package backend

import (
    "fmt"
    "sort"
)

type ServiceCandidate struct { NodeID string; ServiceID string; Status HealthStatus; LatencyMs float64; PacketLossPct float64; CapacityPct float64 }

type RebalancePlan struct { ServiceID string; FromNode string; ToNode string; Reason string }

func BuildRebalancePlan(serviceID string, candidates []ServiceCandidate) (RebalancePlan, error) {
    if serviceID == "" { return RebalancePlan{}, fmt.Errorf("service id required") }
    eligible := make([]ServiceCandidate, 0, len(candidates))
    for _, c := range candidates { if c.ServiceID == serviceID && c.Status == HealthHealthy && c.CapacityPct >= 0 && c.CapacityPct <= 100 { eligible = append(eligible, c) } }
    if len(eligible) < 2 { return RebalancePlan{}, fmt.Errorf("insufficient healthy candidates") }
    sort.Slice(eligible, func(i,j int) bool { if eligible[i].CapacityPct != eligible[j].CapacityPct { return eligible[i].CapacityPct < eligible[j].CapacityPct }; return eligible[i].LatencyMs < eligible[j].LatencyMs })
    from := eligible[len(eligible)-1]
    to := eligible[0]
    if from.NodeID == to.NodeID { return RebalancePlan{}, fmt.Errorf("no rebalance target") }
    return RebalancePlan{ServiceID: serviceID, FromNode: from.NodeID, ToNode: to.NodeID, Reason: "healthy lower-capacity/latency target selected"}, nil
}
