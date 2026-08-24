package backend

import "fmt"

type MeshTransport string
const (
    MeshWireGuard MeshTransport = "wireguard"
    MeshHeadscale MeshTransport = "headscale"
    MeshTailscale MeshTransport = "tailscale"
    MeshHysteria2 MeshTransport = "hysteria2"
)

type MeshNode struct { ID, Name, Endpoint string; Transport MeshTransport; Authorized, Enabled bool }
type MeshLinkPolicy struct { SourceID, DestinationID string; Allowed bool; KeepaliveSeconds int }

// Mesh connectivity is an overlay: every participating endpoint still needs
// some underlying bearer (Ethernet, Wi-Fi, cellular, or Internet transit).
// Removing a cable is therefore supported only when another bearer is active.
func ValidateMeshNode(n MeshNode) error {
    if n.ID=="" || n.Name=="" || n.Endpoint=="" { return fmt.Errorf("mesh node identity/endpoint required") }
    if !n.Authorized { return fmt.Errorf("mesh node is not authorized") }
    switch n.Transport { case MeshWireGuard, MeshHeadscale, MeshTailscale, MeshHysteria2: default: return fmt.Errorf("unsupported mesh transport") }
    return nil
}

func ValidateMeshLink(p MeshLinkPolicy) error {
    if p.SourceID=="" || p.DestinationID=="" || p.SourceID==p.DestinationID { return fmt.Errorf("invalid mesh endpoints") }
    if !p.Allowed { return fmt.Errorf("mesh link is not allowed by policy") }
    if p.KeepaliveSeconds < 0 { return fmt.Errorf("invalid keepalive") }
    return nil
}
