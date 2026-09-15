package backend

import "sort"

type GeoFlowTarget struct {
	ID        string   `json:"id"`
	Provider  string   `json:"provider"`
	Region    string   `json:"region"`
	Countries []string `json:"countries,omitempty"`
	Endpoint  string   `json:"endpoint"`
	Healthy   bool     `json:"healthy"`
	LatencyMS float64  `json:"latency_ms"`
	Weight    float64  `json:"weight"`
}

type GeoFlowPolicy struct {
	ID              string         `json:"id"`
	Service         string         `json:"service"`
	Mode            string         `json:"mode"`
	Targets         []GeoFlowTarget `json:"targets"`
	FailoverEnabled bool           `json:"failover_enabled"`
}

func ScoreGeoFlowTarget(t GeoFlowTarget) float64 {
	if !t.Healthy {
		return 1e12
	}
	weight := t.Weight
	if weight <= 0 {
		weight = 1
	}
	return t.LatencyMS / weight
}

func RankGeoFlowTargets(targets []GeoFlowTarget) []GeoFlowTarget {
	out := append([]GeoFlowTarget(nil), targets...)
	sort.SliceStable(out, func(i, j int) bool {
		return ScoreGeoFlowTarget(out[i]) < ScoreGeoFlowTarget(out[j])
	})
	return out
}

func DefaultGeoFlowPolicy() GeoFlowPolicy {
	return GeoFlowPolicy{
		ID: "global-routing",
		Service: "ftn",
		Mode: "geo-latency-health",
		FailoverEnabled: true,
		Targets: []GeoFlowTarget{},
	}
}
