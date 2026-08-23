package backend

import "testing"

func TestExecuteService(t *testing.T) {
	s := LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"n1", Authorized:true}
	got, err := ExecuteService(DryRunExecutor{}, s, ServiceLive)
	if err != nil { t.Fatal(err) }
	if !got.Applied || got.ServiceID != "svc1" || got.Desired != ServiceLive { t.Fatalf("unexpected result: %#v", got) }
}

func TestExecuteServiceRejectsUnauthorized(t *testing.T) {
	s := LiveService{ID:"svc1", Name:"FTN DNS", NodeID:"n1", Authorized:false}
	if _, err := ExecuteService(DryRunExecutor{}, s, ServiceLive); err == nil { t.Fatal("expected authorization error") }
}
