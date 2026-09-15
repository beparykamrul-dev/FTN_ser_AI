package backend

type DNSProvider struct { ID string `json:"id"`; Name string `json:"name"`; Kind string `json:"kind"`; Enabled bool `json:"enabled"`; ReadOnly bool `json:"read_only"` }
type DNSRecord struct { Zone string `json:"zone"`; Name string `json:"name"`; Type string `json:"type"`; Content string `json:"content"`; TTL uint32 `json:"ttl"` }
