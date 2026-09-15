package backend

import "time"

type AuditEvent struct { ID string `json:"id"`; Actor string `json:"actor"`; Action string `json:"action"`; Resource string `json:"resource"`; Success bool `json:"success"`; Timestamp time.Time `json:"timestamp"` }
