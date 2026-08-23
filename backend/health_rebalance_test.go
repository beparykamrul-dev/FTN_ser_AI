package backend

import "testing"

func TestBuildRebalancePlan(t *testing.T) {
    plan, err := BuildRebalancePlan("svc", []ServiceCandidate{
        {NodeID:"n1", ServiceID:"svc", Status:HealthHealthy, LatencyMs:20, CapacityPct:30},
        {NodeID:"n2", ServiceID:"svc", Status:HealthHealthy, LatencyMs:40, CapacityPct:90},
        {NodeID:"n3", ServiceID:"svc", Status:HealthDegraded, LatencyMs:5, CapacityPct:10},
    })
    if err != nil { t.Fatal(err) }
    if plan.FromNode != "n2" || plan.ToNode != "n1" { t.Fatalf("unexpected plan: %#v", plan) }
}

func TestBuildRebalancePlanNeedsTwoHealthyNodes(t *testing.T) {
    if _, err := BuildRebalancePlan("svc", []ServiceCandidate{{NodeID:"n1", ServiceID:"svc", Status:HealthHealthy}}); err == nil { t.Fatal("expected insufficient candidate error") }
}
