package backend

func MeshNodeEligible(n DNSMeshNode) bool { return n.Address!="" && n.Healthy }
func RouteEligible(r AnycastRoute) bool { return r.Prefix!="" && r.Enabled && r.Healthy }
