package backend

import "testing"

func TestDecideHealthyNode(t *testing.T) {
	action, ok := Decide(NodeObservation{NodeID:"node-1", Healthy:true, CPUPercent:30, MemoryPercent:40, LatencyMS:20, PacketLoss:0})
	if ok || action.ID != "" { t.Fatalf("unexpected action: %+v", action) }
}
func TestDecideOverloadedNode(t *testing.T) {
	action, ok := Decide(NodeObservation{NodeID:"node-2", Healthy:true, CPUPercent:95, MemoryPercent:40, LatencyMS:20, PacketLoss:0})
	if !ok { t.Fatal("expected rebalance action") }
	if action.Kind != "service.rebalance" || !action.RequiresAuth { t.Fatalf("unexpected action: %+v", action) }
}
func TestActionStoreApproval(t *testing.T) {
	s:=NewActionStore(); action:=AIAction{ID:"a1",Kind:"service.rebalance",Target:"node-1",RequiresAuth:true}
	if err:=s.Put(action);err!=nil{t.Fatal(err)}
	if err:=s.Approve("a1");err!=nil{t.Fatal(err)}
	items:=s.List(); if len(items)!=1 || !items[0].Approved{t.Fatalf("approval not persisted: %+v",items)}
}
