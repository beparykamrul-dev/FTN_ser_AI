package backend

import "testing"

func TestTrafficFabricCapabilities(t *testing.T) {
    r:=NewTrafficFabricRegistry()
    p:=TrafficFabricProfile{ID:"tf1",Name:"FTN Traffic Fabric",NodeID:"node1",Authorized:true,Enabled:true,Capabilities:[]TrafficFabricCapability{TrafficECMP,TrafficMACDiscovery,TrafficGeo,TrafficTPM2,TrafficLogs,TrafficMASQUE,TrafficShadowsocks,TrafficXray,TrafficGugiGilong}}
    if err:=r.Upsert(p);err!=nil{t.Fatal(err)}
    got,ok:=r.Get("tf1");if !ok||len(got.Capabilities)!=9{t.Fatalf("unexpected profile: %#v",got)}
}

func TestTrafficFabricRejectsUnauthorized(t *testing.T){
    r:=NewTrafficFabricRegistry()
    if err:=r.Upsert(TrafficFabricProfile{ID:"tf1",Name:"FTN Traffic Fabric",NodeID:"node1",Capabilities:[]TrafficFabricCapability{TrafficECMP}});err==nil{t.Fatal("expected authorization error")}
}
