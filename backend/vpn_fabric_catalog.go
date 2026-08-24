package backend

// VPNFabricCapability describes a governed FTN VPN/secure-overlay integration.
// These are integration capabilities; runtime activation remains policy- and
// authorization-controlled.
type VPNFabricCapability string

const (
    VPNWireGuard VPNFabricCapability = "wireguard"
    VPNAmneziaWG VPNFabricCapability = "amneziawg"
    VPNAetherCore VPNFabricCapability = "aether-core"
    VPNStrongSwanIKEv2 VPNFabricCapability = "strongswan-ikev2-ipsec"
    VPNOpenVPN VPNFabricCapability = "openvpn"
    VPNHeadscale VPNFabricCapability = "headscale"
    VPNTailscale VPNFabricCapability = "tailscale"
    VPNHysteria2 VPNFabricCapability = "hysteria2"
    VPNMASQUE VPNFabricCapability = "masque-http3-quic"
    VPNIPsec VPNFabricCapability = "ipsec"
    VPNZeroTrustGateway VPNFabricCapability = "zero-trust-gateway"
)

// RecommendedVPNFabric returns a deliberately small, complementary set.
// WireGuard is the default high-performance tunnel; Aether-Core orchestrates
// paths/policy; strongSwan provides standards-based IKEv2/IPsec interoperability;
// OpenVPN covers legacy/enterprise interoperability; Headscale/Tailscale provide
// identity-aware mesh control; Hysteria2/MASQUE are optional QUIC transport
// adapters rather than mandatory parallel tunnels.
func RecommendedVPNFabric() []VPNFabricCapability {
    return []VPNFabricCapability{
        VPNWireGuard, VPNAmneziaWG, VPNAetherCore,
        VPNStrongSwanIKEv2, VPNOpenVPN,
        VPNHeadscale, VPNTailscale,
        VPNHysteria2, VPNMASQUE, VPNIPsec,
        VPNZeroTrustGateway,
    }
}
