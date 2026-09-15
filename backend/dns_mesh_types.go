package backend

type DNSMeshNode struct {
	ID        string  `json:"id"`
	Endpoint  string  `json:"endpoint"`
	Healthy   bool    `json:"healthy"`
	LatencyMS float64 `json:"latency_ms"`
}

type AnycastRoute struct {
	Prefix  string `json:"prefix"`
	NodeID  string `json:"node_id"`
	Enabled bool   `json:"enabled"`
	Healthy bool   `json:"healthy"`
}

type DNSMeshStatus struct {
	Zones  []string       `json:"zones"`
	Nodes  []DNSMeshNode  `json:"nodes"`
	Routes []AnycastRoute `json:"routes"`
	DNSSEC bool           `json:"dnssec"`
}
