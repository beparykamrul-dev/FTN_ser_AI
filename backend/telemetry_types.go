package backend

import "time"

type TelemetryEnvelope struct { NodeID string `json:"node_id"`; Kind string `json:"kind"`; Payload any `json:"payload"`; Timestamp time.Time `json:"timestamp"` }
