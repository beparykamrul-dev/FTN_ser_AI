package backend

import "math"

type RouteCandidate struct {
	NodeID      string  `json:"node_id"`
	Healthy     bool    `json:"healthy"`
	LatencyMS   float64 `json:"latency_ms"`
	LossPercent float64 `json:"loss_percent"`
	LoadPercent float64 `json:"load_percent"`
}

type RouteScore struct {
	NodeID string  `json:"node_id"`
	Score  float64 `json:"score"`
}

func ScoreRoute(c RouteCandidate) float64 {
	if !c.Healthy {
		return math.Inf(1)
	}
	return c.LatencyMS + (c.LossPercent * 20) + (c.LoadPercent * 0.1)
}

func SelectBestRoute(candidates []RouteCandidate) (RouteCandidate, bool) {
	var best RouteCandidate
	bestScore := math.Inf(1)
	found := false
	for _, c := range candidates {
		score := ScoreRoute(c)
		if score < bestScore {
			best, bestScore, found = c, score, true
		}
	}
	return best, found
}
