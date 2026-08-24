package backend

import "testing"

func TestFlowFabricCapabilities(t *testing.T) {
    r:=NewFlowFabricRegistry()
    p:=FlowFabricProfile{ID:"flow1",NodeID:"node1",Authorized:true,Enabled:true,Capabilities:[]FlowPipelineCapability{FlowAetherCore,FlowSiLK,FlowRWWFlowpack}}
    if err:=r.Upsert(p);err!=nil{t.Fatal(err)}
    got,ok:=r.Get("flow1");if !ok||len(got.Capabilities)!=3{t.Fatalf("unexpected profile: %#v",got)}
}

func TestFlowFabricRejectsUnknownCapability(t *testing.T){
    r:=NewFlowFabricRegistry()
    if err:=r.Upsert(FlowFabricProfile{ID:"flow1",NodeID:"node1",Authorized:true,Capabilities:[]FlowPipelineCapability{"unknown"}});err==nil{t.Fatal("expected unsupported capability error")}
}
