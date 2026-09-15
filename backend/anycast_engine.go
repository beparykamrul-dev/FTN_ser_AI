package backend

import "sort"

type AnycastNode struct {
	NodeID       string  `json:"node_id"`
	Prefix       string  `json:"prefix"`
	ASN          uint32  `json:"asn"`
	Enabled      bool    `json:"enabled"`
	Healthy      bool    `json:"healthy"`
	LatencyMS    float64 `json:"latency_ms"`
	PacketLoss   float64 `json:"packet_loss_pct"`
	Load         float64 `json:"load_pct"`
	Score        float64 `json:"score"`
}

func ScoreAnycastNode(n AnycastNode) float64 {
	if !n.Enabled || !n.Healthy { return 1e12 }
	loss := n.PacketLoss * 20
	load := n.Load * 0.10
	return n.LatencyMS + loss + load
}

func RankAnycastNodes(nodes []AnycastNode) []AnycastNode {
	out := append([]AnycastNode(nil), nodes...)
	for i := range out { out[i].Score = ScoreAnycastNode(out[i]) }
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score < out[j].Score })
	return out
}

func BestAnycastNode(nodes []AnycastNode) (AnycastNode, bool) {
	ranked := RankAnycastNodes(nodes)
	if len(ranked) == 0 || ranked[0].Score >= 1e12 { return AnycastNode{}, false }
	return ranked[0], true
}
