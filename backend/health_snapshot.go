package backend

import "time"

type DependencyHealth struct { Service string `json:"service"`; Healthy bool `json:"healthy"`; CheckedAt time.Time `json:"checked_at"`; Detail string `json:"detail,omitempty"` }
