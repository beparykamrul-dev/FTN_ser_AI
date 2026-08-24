package backend

import (
    "testing"
    "time"
)

func TestApprovedReadyRolloutCanGoLive(t *testing.T) {
    r:=NewRolloutRegistry()
    if err:=r.Advance(ServiceRollout{ID:"r1",ServiceID:"svc1",Stage:RolloutLive,Approved:true,Ready:true},time.Now());err!=nil{t.Fatal(err)}
    x,ok:=r.Get("r1");if !ok||x.Stage!=RolloutLive{t.Fatalf("unexpected rollout: %#v",x)}
}
func TestLiveRolloutRequiresApprovalAndReadiness(t *testing.T){
    r:=NewRolloutRegistry(); now:=time.Now()
    if err:=r.Advance(ServiceRollout{ID:"r1",ServiceID:"svc1",Stage:RolloutLive,Approved:false,Ready:true},now);err==nil{t.Fatal("expected approval error")}
    if err:=r.Advance(ServiceRollout{ID:"r2",ServiceID:"svc1",Stage:RolloutLive,Approved:true,Ready:false},now);err==nil{t.Fatal("expected readiness error")}
}
