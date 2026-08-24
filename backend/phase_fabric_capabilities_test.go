package backend

import "testing"

func TestPhaseFabricCatalog(t *testing.T) {
    catalog := PhaseFabricCatalog()
    required := map[PhaseFabricCapability]bool{
        CapHysteria2:false, CapWireGuardMesh:false, CapHeadscale:false,
        CapVPP:false, CapOPNsense:false, CapNetmaker:false,
        CapTrafficShaping:false, CapTCPBBR:false, CapNetFlowIPFIX:false,
        CapZeek:false, CapTPM2:false, CapMASQUE:false, CapShadowsocks:false,
        CapXray:false,
    }
    for _, c := range catalog { if _, ok := required[c]; ok { required[c] = true } }
    for c, ok := range required { if !ok { t.Fatalf("missing capability %q", c) } }
}

func TestNormalizePhaseCapability(t *testing.T) {
    if got := NormalizePhaseCapability("  WireGuard-Mesh "); got != CapWireGuardMesh { t.Fatalf("got %q", got) }
}
