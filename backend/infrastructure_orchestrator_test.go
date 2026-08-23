package backend

import (
    "testing"
    "time"
)

func TestInfrastructurePlacement(t *testing.T) {
    o:=NewInfrastructureOrchestrator()
    if err:=o.RegisterNode(NodeCapability{NodeID:"n1",CPUCapacityPct:30,MemoryCapacityPct:20,StorageCapacityPct:40,NetworkCapacityPct:25,Healthy:true,Authorized:true});err!=nil{t.Fatal(err)}
    if err:=o.RegisterNode(NodeCapability{NodeID:"n2",CPUCapacityPct:70,MemoryCapacityPct:60,StorageCapacityPct:80,NetworkCapacityPct:75,Healthy:true,Authorized:true});err!=nil{t.Fatal(err)}
    p,err:=o.PlanPlacement("svc1",time.Now());if err!=nil{t.Fatal(err)}
    if p.NodeID!="n1"{t.Fatalf("expected n1, got %#v",p)}
}

func TestInfrastructureRejectsUnauthorizedNode(t *testing.T){
    o:=NewInfrastructureOrchestrator()
    if err:=o.RegisterNode(NodeCapability{NodeID:"n1",Healthy:true,Authorized:false});err==nil{t.Fatal("expected authorization error")}
}

func TestInfrastructureNeedsEligibleNode(t *testing.T){
    o:=NewInfrastructureOrchestrator()
    if _,err:=o.PlanPlacement("svc1",time.Now());err==nil{t.Fatal("expected no eligible node error")}
}
