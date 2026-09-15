package backend

import "time"

type Job struct { ID string `json:"id"`; Kind string `json:"kind"`; Status string `json:"status"`; CreatedAt time.Time `json:"created_at"`; StartedAt *time.Time `json:"started_at,omitempty"`; FinishedAt *time.Time `json:"finished_at,omitempty"`; Error string `json:"error,omitempty"` }
