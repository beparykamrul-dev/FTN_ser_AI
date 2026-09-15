package backend

type ControlSnapshot struct { Capabilities CapabilitySet `json:"capabilities"`; Nodes []NodeMetrics `json:"nodes"`; Alerts []Alert `json:"alerts"`; Mesh DNSMeshStatus `json:"dns_mesh"`; DB DBHealth `json:"db"` }
func BuildSnapshot(m *MonitorStore, mesh *DNSMeshStore, db *DBStore) ControlSnapshot { return ControlSnapshot{Capabilities:DefaultCapabilities(),Nodes:m.Nodes(),Alerts:m.Alerts(),Mesh:mesh.Get(),DB:db.Get()} }
