package backend

type CapabilitySet struct { Monitoring bool `json:"monitoring"`; Database bool `json:"database"`; DNSMesh bool `json:"dns_mesh"`; Anycast bool `json:"anycast"`; IPAM bool `json:"ipam"`; Providers bool `json:"providers"`; Audit bool `json:"audit"` }
func DefaultCapabilities() CapabilitySet { return CapabilitySet{Monitoring:true,Database:true,DNSMesh:true,Anycast:true,IPAM:true,Providers:true,Audit:true} }
