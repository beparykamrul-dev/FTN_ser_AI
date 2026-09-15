package backend

type DBInstance struct { ID string `json:"id"`; Engine string `json:"engine"`; Address string `json:"address"`; Healthy bool `json:"healthy"` }
type DBHealth struct { Instances []DBInstance `json:"instances"`; Connections int `json:"connections"`; Locks int `json:"locks"`; WALBytes uint64 `json:"wal_bytes"` }
