package backend

import "testing"

func TestTelemetryRegistry(t *testing.T){
 r:=NewTelemetryRegistry()
 p:=TelemetryProfile{ID:"t1",NodeID:"n1",Authorized:true,Capabilities:[]TelemetryCapability{TelemetryTurnKey,TelemetryOpenRC,TelemetryVXLAN,TelemetryIPSec,TelemetryNetSA,TelemetryIPFIX,TelemetryNetFlowV5,TelemetryNetFlowV9,TelemetryRwfowpack,TelemetrySiLKAnalysis,TelemetryYAFDpacketPlugin,TelemetryPerlPython,TelemetrySolaris,TelemetryOpenBSD,TelemetryCygwin}}
 if err:=r.Upsert(p);err!=nil{t.Fatal(err)}
 got,ok:=r.Get("t1");if !ok||len(got.Capabilities)!=15{t.Fatalf("unexpected profile: %#v",got)}
}

func TestTelemetryRegistryRejectsUnauthorized(t *testing.T){
 r:=NewTelemetryRegistry()
 if err:=r.Upsert(TelemetryProfile{ID:"t1",NodeID:"n1",Capabilities:[]TelemetryCapability{TelemetryIPFIX}});err==nil{t.Fatal("expected authorization error")}
}
