package backend

import (
    "testing"
    "time"
)

func TestServiceReconcilerAppliesDesiredState(t *testing.T) {
    reg := NewLiveServiceRegistry()
    now := time.Now()
    if err := reg.Activate(LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"n1", Authorized:true}, now); err != nil { t.Fatal(err) }
    r, err := NewServiceReconciler(reg); if err != nil { t.Fatal(err) }
    if err := r.Reconcile(RuntimeAction{ID:"a1", ServiceID:"svc1", Desired:ServiceDegraded, At:now.Add(time.Second)}); err != nil { t.Fatal(err) }
    s, _ := reg.Get("svc1")
    if s.State != ServiceDegraded { t.Fatalf("expected degraded, got %s", s.State) }
}

func TestServiceReconcilerRejectsDuplicateAction(t *testing.T) {
    reg := NewLiveServiceRegistry(); now := time.Now()
    if err := reg.Activate(LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"n1", Authorized:true}, now); err != nil { t.Fatal(err) }
    r, _ := NewServiceReconciler(reg)
    a := RuntimeAction{ID:"a1", ServiceID:"svc1", Desired:ServiceStopped, At:now}
    if err := r.Reconcile(a); err != nil { t.Fatal(err) }
    if err := r.Reconcile(a); err == nil { t.Fatal("expected duplicate action rejection") }
}
