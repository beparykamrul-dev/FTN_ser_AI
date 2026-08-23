package backend

import (
    "testing"
    "time"
)

func TestAuthorizedServiceBecomesLive(t *testing.T) {
    r := NewLiveServiceRegistry()
    now := time.Now()
    if err := r.Activate(LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"node1", Authorized:true}, now); err != nil { t.Fatal(err) }
    s, ok := r.Get("svc1")
    if !ok || s.State != ServiceLive { t.Fatalf("unexpected service state: %#v", s) }
}

func TestUnauthorizedServiceCannotBecomeLive(t *testing.T) {
    r := NewLiveServiceRegistry()
    if err := r.Activate(LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"node1", Authorized:false}, time.Now()); err == nil { t.Fatal("expected authorization error") }
}

func TestLiveServiceHealthTransition(t *testing.T) {
    r := NewLiveServiceRegistry()
    now := time.Now()
    if err := r.Activate(LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"node1", Authorized:true}, now); err != nil { t.Fatal(err) }
    if err := r.SetHealth("svc1", false, now.Add(time.Second)); err != nil { t.Fatal(err) }
    s, _ := r.Get("svc1")
    if s.State != ServiceDegraded { t.Fatalf("expected degraded, got %s", s.State) }
}
