package backend

import "testing"

func TestNetworkEdgeCapabilities(t *testing.T) {
    r := NewNetworkEdgeRegistry()
    p := NetworkEdgeProfile{ID:"edge1", Name:"FTN Edge", NodeID:"node1", Authorized:true, Enabled:true, Capabilities:[]NetworkEdgeCapability{EdgeNginx, EdgeDynamicMesh, EdgeDNSStealth, EdgeEBPFXDP, EdgeVPNZeroTrust}}
    if err := r.Upsert(p); err != nil { t.Fatal(err) }
    got, ok := r.Get("edge1")
    if !ok || len(got.Capabilities) != 5 { t.Fatalf("unexpected profile: %#v", got) }
}

func TestNetworkEdgeRejectsUnauthorized(t *testing.T) {
    r := NewNetworkEdgeRegistry()
    if err := r.Upsert(NetworkEdgeProfile{ID:"edge1", Name:"FTN Edge", NodeID:"node1", Capabilities:[]NetworkEdgeCapability{EdgeNginx}}); err == nil { t.Fatal("expected authorization error") }
}
