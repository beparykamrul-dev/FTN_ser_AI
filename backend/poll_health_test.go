package backend

import (
    "testing"
    "time"
)

func TestEvaluateHealth(t *testing.T) {
    now := time.Now()
    status, err := EvaluateHealth(HealthObservation{NodeID:"n1", ServiceID:"s1", LatencyMs:10, PacketLossPct:0, At:now}, HealthPolicy{MaxLatencyMs:50, MaxPacketLossPct:1})
    if err != nil || status != HealthHealthy { t.Fatalf("expected healthy, got %s: %v", status, err) }
    status, err = EvaluateHealth(HealthObservation{NodeID:"n1", ServiceID:"s1", LatencyMs:100, PacketLossPct:0, At:now}, HealthPolicy{MaxLatencyMs:50, MaxPacketLossPct:1})
    if err != nil || status != HealthDegraded { t.Fatalf("expected degraded, got %s: %v", status, err) }
}

func TestHealthStore(t *testing.T) {
    s := NewHealthStore()
    o := HealthObservation{NodeID:"n1", ServiceID:"s1", Status:HealthHealthy, At:time.Now()}
    if err := s.Put(o); err != nil { t.Fatal(err) }
    if got, ok := s.Get("n1", "s1"); !ok || got.Status != HealthHealthy { t.Fatalf("unexpected health: %#v", got) }
}
