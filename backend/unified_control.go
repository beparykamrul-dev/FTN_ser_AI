package backend

import (
	"net/http"
	"strings"
	"time"
)

type FTNResource struct { ID string `json:"id"`; Kind string `json:"kind"`; Name string `json:"name"`; Location string `json:"location,omitempty"`; Endpoint string `json:"endpoint,omitempty"`; Local bool `json:"local"`; Healthy bool `json:"healthy"`; LatencyMS float64 `json:"latency_ms,omitempty"` }
type DataPlaneStatus struct { RAMFirst bool `json:"ram_first"`; RedisHotState bool `json:"redis_hot_state"`; PostgresSource bool `json:"postgres_source_of_truth"`; AsyncPersist bool `json:"async_persist"`; DBOnRequestPath bool `json:"db_on_request_path"`; CacheTTL time.Duration `json:"cache_ttl"` }

func (s *APIServer) controlSnapshot(w http.ResponseWriter,r *http.Request){
	if r.Method!=http.MethodGet{w.Header().Set("Allow",http.MethodGet);writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"});return}
	p:=DefaultPersistencePolicy(); writeJSON(w,http.StatusOK,map[string]any{"ok":true,"resources":[]FTNResource{{ID:"ftndns",Kind:"dns",Name:"FTNDNS",Location:"global",Healthy:true},{ID:"ddns",Kind:"ddns",Name:"DDNS",Location:"global",Healthy:true},{ID:"agents",Kind:"device-gateway",Name:"FTN Agents",Location:"global",Healthy:true}},"authorities":AuthorityZones(),"providers":StrategicProviderPolicy(),"geoflow":DefaultGeoFlowPolicy(),"data_plane":DataPlaneStatus{RAMFirst:p.RAMFirst,RedisHotState:p.RedisHotState,PostgresSource:p.PostgresDurable,AsyncPersist:p.AsyncPersist,DBOnRequestPath:!p.ReadPathDBFree,CacheTTL:p.CacheTTL}})
}
func (s *APIServer) controlHealth(w http.ResponseWriter,r *http.Request){if r.Method!=http.MethodGet{w.Header().Set("Allow",http.MethodGet);writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"});return};writeJSON(w,http.StatusOK,map[string]any{"ok":true,"status":"healthy","latency_policy":"geo-latency-health","bdix_dependency":false})}
func (s *APIServer) controlCapabilities(w http.ResponseWriter,r *http.Request){if r.Method!=http.MethodGet{w.Header().Set("Allow",http.MethodGet);writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"});return};writeJSON(w,http.StatusOK,DefaultCapabilities())}
func (s *APIServer) dnsMesh(w http.ResponseWriter,r *http.Request){if r.Method!=http.MethodGet{w.Header().Set("Allow",http.MethodGet);writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"});return};writeJSON(w,http.StatusOK,s.mesh.Get())}
func (s *APIServer) controlResource(w http.ResponseWriter,r *http.Request){
	if r.Method!=http.MethodGet{w.Header().Set("Allow",http.MethodGet);writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"});return}
	path:=strings.TrimPrefix(r.URL.Path,"/api/v1/control/")
	switch path{
	case "nodes":writeJSON(w,http.StatusOK,map[string]any{"nodes":s.monitor.Nodes()})
	case "alerts":writeJSON(w,http.StatusOK,map[string]any{"alerts":s.monitor.Alerts()})
	case "dns":writeJSON(w,http.StatusOK,s.mesh.Get())
	case "db":writeJSON(w,http.StatusOK,s.db.Get())
	case "providers":writeJSON(w,http.StatusOK,map[string]any{"providers":s.providers.List(),"strategy":StrategicProviderPolicy()})
	case "authorities":writeJSON(w,http.StatusOK,map[string]any{"zones":AuthorityZones()})
	case "anycast":writeJSON(w,http.StatusOK,map[string]any{"status":"anycast-ready","routing":"geo-latency-health","prefixes":DefaultAnycastPrefixes()})
	case "ddns":writeJSON(w,http.StatusOK,map[string]any{"zone":"ftnddns.net","status":"engine-ready"})
	case "dataplane":writeJSON(w,http.StatusOK,DefaultPersistencePolicy())
	case "geoflow":writeJSON(w,http.StatusOK,DefaultGeoFlowPolicy())
	case "agents":writeJSON(w,http.StatusOK,map[string]any{"agents":s.agents.List(),"websocket":"/api/v1/ws/agent"})
	case "ipam":writeJSON(w,http.StatusOK,map[string]any{"status":"registry-ready"})
	case "jobs":writeJSON(w,http.StatusOK,map[string]any{"status":"registry-ready"})
	case "audit":writeJSON(w,http.StatusOK,map[string]any{"status":"registry-ready"})
	default:writeJSON(w,http.StatusNotFound,map[string]string{"error":"control resource not found"})
	}
}
