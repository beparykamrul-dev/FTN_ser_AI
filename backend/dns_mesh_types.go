package backend

type AnycastRoute struct {
	Prefix  string `json:"prefix"`
	NodeID  string `json:"node_id"`
	Enabled bool   `json:"enabled"`
	Healthy bool   `json:"healthy"`
}

type DNSMeshStatus struct {
	Zones  []string      `json:"zones"`
	Nodes  []DNSMeshNode `json:"nodes"`
	Routes []AnycastRoute `json:"routes"`
	DNSSEC bool          `json:"dnssec"`
}
