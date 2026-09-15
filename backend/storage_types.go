package backend

type StorageDevice struct { Device string `json:"device"`; Mountpoint string `json:"mountpoint,omitempty"`; TotalBytes uint64 `json:"total_bytes"`; UsedBytes uint64 `json:"used_bytes"`; FreeBytes uint64 `json:"free_bytes"`; InodesUsed uint64 `json:"inodes_used"`; InodesTotal uint64 `json:"inodes_total"` }
