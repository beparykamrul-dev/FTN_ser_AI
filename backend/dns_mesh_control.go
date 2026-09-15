package backend

import (
    "net/http"
    "sync"
    "time"
)

type DNSMeshNode struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    Region        string    `json:"region"`
    Address       string    `json:"address"`
    IPv6          string    `json:"ipv6,omitempty"`
    Anycast       bool      `json:"anycast"`
    Authoritative bool      `json:"authoritative"`
    Healthy       bool      `json:"healthy"`
    LatencyMS     float64   `json:"latency_ms"`
    PacketLoss    float64   `json:"packet_loss"`
    LastSeen      time.Time `json:"last_seen"`
}

type DNSMeshState struct {
    mu            sync.RWMutex
    Domains       []string
    Nodes         map[string]DNSMeshNode
    AnycastRoutes []string
}

func NewDNSMeshState() *DNSMeshState {
    return &DNSMeshState{
        Domains: []string{"familytimenet.com", "ftnddns.net"},
        Nodes: make(map[string]DNSMeshNode),
        AnycastRoutes: []string{},
    }
}

func (s *DNSMeshState) Snapshot() map[string]any {
    s.mu.RLock(); defer s.mu.RUnlock()
    nodes := make([]DNSMeshNode, 0, len(s.Nodes))
    for _, n := range s.Nodes { nodes = append(nodes, n) }
    return map[string]any{
        "status": "controller-ready",
        "domains": append([]string(nil), s.Domains...),
        "nodes": nodes,
        "anycast_routes": append([]string(nil), s.AnycastRoutes...),
        "policy": map[string]any{
            "ipv4": true, "ipv6": true,
            "health_withdrawal": true,
            "dnssec": true,
            "external_providers_optional": true,
        },
    }
}

func (s *APIServer) dnsMesh(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet { w.Header().Set("Allow", http.MethodGet); writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error":"method not allowed"}); return }
    writeJSON(w, http.StatusOK, s.mesh.Snapshot())
}
