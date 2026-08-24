package backend

import "testing"
import "time"

func TestServiceReadiness(t *testing.T) {
    e := NewServiceReadinessEngine()
    s := LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"n1", State:ServiceLive, Authorized:true, LastHealthCheck:time.Now()}
    r, err := e.Evaluate(s, time.Now()); if err != nil { t.Fatal(err) }
    if !r.Ready { t.Fatalf("expected ready: %#v", r) }
}

func TestServiceReadinessRejectsUnhealthyState(t *testing.T) {
    e := NewServiceReadinessEngine()
    s := LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"n1", State:ServiceDegraded, Authorized:true, LastHealthCheck:time.Now()}
    r, err := e.Evaluate(s, time.Now()); if err != nil { t.Fatal(err) }
    if r.Ready { t.Fatal("expected not ready") }
}
