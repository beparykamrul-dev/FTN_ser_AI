package backend

type DNSMeshNode struct { ID string `json:"id"`; Region string `json:"region"`; Address string `json:"address"`; Healthy bool `json:"healthy"`; AnycastReady bool `json:"anycast_ready"` }
type AnycastRoute struct { Prefix string `json:"prefix"`; NodeID string `json:"node_id"`; Enabled bool `json:"enabled"`; Healthy bool `json:"healthy"` }
type DNSMeshStatus struct { Zones []string `json:"zones"`; Nodes []DNSMeshNode `json:"nodes"`; Routes []AnycastRoute `json:"routes"`; DNSSEC bool `json:"dnssec"` }
