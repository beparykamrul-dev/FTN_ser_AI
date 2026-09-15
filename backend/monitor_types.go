package backend

import "time"

type MetricSample struct { Name string `json:"name"`; Value float64 `json:"value"`; Unit string `json:"unit,omitempty"`; Labels map[string]string `json:"labels,omitempty"`; ObservedAt time.Time `json:"observed_at"` }
type NodeMetrics struct { NodeID string `json:"node_id"`; Hostname string `json:"hostname,omitempty"`; Samples []MetricSample `json:"samples"`; ObservedAt time.Time `json:"observed_at"` }
type Alert struct { ID string `json:"id"`; Severity string `json:"severity"`; Source string `json:"source"`; Message string `json:"message"`; Active bool `json:"active"`; ObservedAt time.Time `json:"observed_at"` }
