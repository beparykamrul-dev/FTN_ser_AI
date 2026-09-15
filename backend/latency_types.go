package backend

type LatencyProbe struct { Source string `json:"source"`; Target string `json:"target"`; LatencyMS float64 `json:"latency_ms"`; PacketLoss float64 `json:"packet_loss"`; Healthy bool `json:"healthy"` }
type LatencyMatrix struct { Probes []LatencyProbe `json:"probes"` }
