package backend

type IPPool struct { ID string `json:"id"`; CIDR string `json:"cidr"`; Family int `json:"family"`; Description string `json:"description,omitempty"` }
type IPAssignment struct { PoolID string `json:"pool_id"`; Address string `json:"address"`; NodeID string `json:"node_id,omitempty"`; Status string `json:"status"` }
