package backend

import ("encoding/json"; "net/http"; "time")

type MonitorAPI struct { Monitor *MonitorStore; Mesh *DNSMeshStore; DB *DBStore }
func (a *MonitorAPI) Handler() http.Handler { mux:=http.NewServeMux(); mux.HandleFunc("/api/v1/monitor/nodes",a.nodes); mux.HandleFunc("/api/v1/monitor/alerts",a.alerts); mux.HandleFunc("/api/v1/dns/mesh",a.mesh); mux.HandleFunc("/api/v1/db/health",a.db); return mux }
func jsonOut(w http.ResponseWriter,v any){w.Header().Set("Content-Type","application/json"); _=json.NewEncoder(w).Encode(v)}
func (a *MonitorAPI) nodes(w http.ResponseWriter,_ *http.Request){jsonOut(w,a.Monitor.Nodes())}
func (a *MonitorAPI) alerts(w http.ResponseWriter,_ *http.Request){jsonOut(w,a.Monitor.Alerts())}
func (a *MonitorAPI) mesh(w http.ResponseWriter,_ *http.Request){jsonOut(w,a.Mesh.Get())}
func (a *MonitorAPI) db(w http.ResponseWriter,_ *http.Request){jsonOut(w,a.DB.Get())}
var _ = time.Time{}
