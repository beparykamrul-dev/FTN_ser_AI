package backend

import "strings"

// PhaseFabricCapability identifies an integration that may participate in the
// FTN Phase 1-3 service fabric. It is a capability declaration, not a command
// to install, intercept traffic, or modify a host.
type PhaseFabricCapability string

const (
    CapKeepalive PhaseFabricCapability = "keepalive"
    CapHysteria2 PhaseFabricCapability = "hysteria2"
    CapWireGuardMesh PhaseFabricCapability = "wireguard-mesh"
    CapHeadscale PhaseFabricCapability = "headscale"
    CapTailscaleMesh PhaseFabricCapability = "tailscale-mesh"
    CapBATMANAdv PhaseFabricCapability = "batman-adv"
    CapBabel PhaseFabricCapability = "babel"
    CapOLSR PhaseFabricCapability = "olsr"
    CapYggdrasil PhaseFabricCapability = "yggdrasil"
    CapVPP PhaseFabricCapability = "vpp"
    CapFDio PhaseFabricCapability = "fdio"
    CapOPNsense PhaseFabricCapability = "opnsense"
    CapNetmaker PhaseFabricCapability = "netmaker"
    CapFirezone PhaseFabricCapability = "firezone"
    CapOpenVPNEnterprise PhaseFabricCapability = "openvpn-enterprise"
    CapFortiSSLVPN PhaseFabricCapability = "fortisslvpn"
    CapGlobalProtect PhaseFabricCapability = "globalprotect"
    CapCiscoSecureClient PhaseFabricCapability = "cisco-secure-client"
    CapSSLVPN PhaseFabricCapability = "ssl-vpn"
    CapNginxReverseProxy PhaseFabricCapability = "nginx-reverse-proxy"
    CapHTTPSOffload PhaseFabricCapability = "https-offload"
    CapTrafficShaping PhaseFabricCapability = "traffic-shaping"
    CapTCPBBR PhaseFabricCapability = "tcp-bbr"
    CapNetFlowIPFIX PhaseFabricCapability = "netflow-ipfix"
    CapFluentBit PhaseFabricCapability = "fluent-bit"
    CapVector PhaseFabricCapability = "vector"
    CapZeek PhaseFabricCapability = "zeek"
    CapWiresharkTshark PhaseFabricCapability = "wireshark-tshark"
    CapLUKS PhaseFabricCapability = "luks"
    CapAnsible PhaseFabricCapability = "ansible"
    CapTPM2 PhaseFabricCapability = "tpm2"
    CapResilioSync PhaseFabricCapability = "resilio-sync"
    CapMagicWormhole PhaseFabricCapability = "magic-wormhole"
    CapThunderbolt PhaseFabricCapability = "thunderbolt"
    CapUSBAutomount PhaseFabricCapability = "usb-automount"
    CapTmpfs PhaseFabricCapability = "tmpfs"
    CapKernelTCPTuning PhaseFabricCapability = "kernel-tcp-tuning"
    CapSubfinderAmass PhaseFabricCapability = "asset-discovery"
    CapMASQUE PhaseFabricCapability = "masque-http3-quic"
    CapShadowsocks PhaseFabricCapability = "shadowsocks"
    CapXray PhaseFabricCapability = "xray"
    CapRawUDP PhaseFabricCapability = "raw-udp"
)

func NormalizePhaseCapability(v PhaseFabricCapability) PhaseFabricCapability {
    return PhaseFabricCapability(strings.ToLower(strings.TrimSpace(string(v))))
}

// PhaseFabricCatalog returns the production-oriented capabilities that can be
// surfaced in the FTN control plane. Security-sensitive tooling is represented
// as governed capabilities and requires authorization and policy checks before
// any runtime adapter can act on it.
func PhaseFabricCatalog() []PhaseFabricCapability {
    return []PhaseFabricCapability{
        CapKeepalive, CapHysteria2, CapWireGuardMesh, CapHeadscale,
        CapTailscaleMesh, CapBATMANAdv, CapBabel, CapOLSR, CapYggdrasil,
        CapVPP, CapFDio, CapOPNsense, CapNetmaker, CapFirezone,
        CapOpenVPNEnterprise, CapFortiSSLVPN, CapGlobalProtect,
        CapCiscoSecureClient, CapSSLVPN, CapNginxReverseProxy,
        CapHTTPSOffload, CapTrafficShaping, CapTCPBBR, CapNetFlowIPFIX,
        CapFluentBit, CapVector, CapZeek, CapWiresharkTshark, CapLUKS,
        CapAnsible, CapTPM2, CapResilioSync, CapMagicWormhole,
        CapThunderbolt, CapUSBAutomount, CapTmpfs, CapKernelTCPTuning,
        CapSubfinderAmass, CapMASQUE, CapShadowsocks, CapXray, CapRawUDP,
    }
}
