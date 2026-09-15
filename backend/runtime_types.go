package backend

type RuntimeSnapshot struct { GoVersion string `json:"go_version"`; Goroutines int `json:"goroutines"`; HeapAlloc uint64 `json:"heap_alloc"`; HeapSys uint64 `json:"heap_sys"`; NumGC uint32 `json:"num_gc"` }
