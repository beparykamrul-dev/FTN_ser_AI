package backend

import (
    "bufio"
    "encoding/json"
    "net/http"
    "os"
    "runtime"
    "strconv"
    "strings"
    "time"
)

type SystemTelemetry struct {
    TimeUTC       time.Time `json:"time_utc"`
    Hostname      string    `json:"hostname"`
    GoVersion     string    `json:"go_version"`
    CPUs          int       `json:"cpus"`
    Goroutines    int       `json:"goroutines"`
    HeapAlloc     uint64    `json:"heap_alloc_bytes"`
    HeapInuse     uint64    `json:"heap_inuse_bytes"`
    SysMemory     uint64    `json:"sys_memory_bytes"`
    NumGC         uint32    `json:"gc_count"`
    Load1         float64   `json:"load1"`
    MemTotal      uint64    `json:"mem_total_bytes"`
    MemAvailable  uint64    `json:"mem_available_bytes"`
}

func (s *APIServer) monitorSystem(w http.ResponseWriter, _ *http.Request) {
    if !method(w, http.MethodGet) { return }
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    host, _ := os.Hostname()
    load1 := readLoad1()
    total, avail := readMemInfo()
    writeJSON(w, http.StatusOK, SystemTelemetry{
        TimeUTC: time.Now().UTC(), Hostname: host, GoVersion: runtime.Version(),
        CPUs: runtime.NumCPU(), Goroutines: runtime.NumGoroutine(),
        HeapAlloc: m.HeapAlloc, HeapInuse: m.HeapInuse, SysMemory: m.Sys,
        NumGC: m.NumGC, Load1: load1, MemTotal: total, MemAvailable: avail,
    })
}

func readLoad1() float64 {
    b, err := os.ReadFile("/proc/loadavg"); if err != nil { return 0 }
    f := strings.Fields(string(b)); if len(f) == 0 { return 0 }
    v, _ := strconv.ParseFloat(f[0], 64); return v
}

func readMemInfo() (uint64, uint64) {
    f, err := os.Open("/proc/meminfo"); if err != nil { return 0, 0 }; defer f.Close()
    var total, avail uint64
    sc := bufio.NewScanner(f)
    for sc.Scan() {
        p := strings.Fields(sc.Text()); if len(p) < 2 { continue }
        v, err := strconv.ParseUint(p[1], 10, 64); if err != nil { continue }
        switch p[0] { case "MemTotal:": total = v * 1024; case "MemAvailable:": avail = v * 1024 }
    }
    return total, avail
}

func method(w http.ResponseWriter, want string) bool {
    if r := w.Header().Get("X-FTN-Method"); r != "" { _ = r }
    return true
}

func encodeMap(w http.ResponseWriter, status int, v map[string]any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(v) }
