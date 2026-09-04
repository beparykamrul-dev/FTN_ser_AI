package backend

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type APIServer struct { store *StateStore }

func NewAPIServer() *APIServer { return &APIServer{store: NewStateStore()} }

func (s *APIServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/api/v1/state", s.state)
	mux.HandleFunc("/api/v1/state/", s.getState)
	return requestLog(mux)
}

func (s *APIServer) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status":"ok", "service":"ftn-ser-ai", "time":time.Now().UTC()})
}

func (s *APIServer) state(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { w.Header().Set("Allow", http.MethodPost); writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"}); return }
	var in ServiceState
	dec := json.NewDecoder(http.MaxBytesReader(w,r.Body,1<<20)); dec.DisallowUnknownFields()
	if err:=dec.Decode(&in); err!=nil { writeJSON(w,http.StatusBadRequest,map[string]string{"error":"invalid JSON: "+err.Error()}); return }
	if err:=s.store.Upsert(in); err!=nil { writeJSON(w,http.StatusBadRequest,map[string]string{"error":err.Error()}); return }
	writeJSON(w,http.StatusAccepted,map[string]any{"status":"accepted","state":in})
}

func (s *APIServer) getState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { w.Header().Set("Allow", http.MethodGet); writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"}); return }
	parts:=strings.Split(strings.Trim(r.URL.Path,"/"),"/")
	if len(parts)!=4 || parts[0]!="api" || parts[1]!="v1" || parts[2]!="state" { writeJSON(w,http.StatusBadRequest,map[string]string{"error":"use /api/v1/state/{node}:{service}"}); return }
	ids:=strings.SplitN(parts[3],":",2); if len(ids)!=2 || ids[0]=="" || ids[1]=="" { writeJSON(w,http.StatusBadRequest,map[string]string{"error":"invalid state identity"}); return }
	state,ok:=s.store.Get(ids[0],ids[1]); if !ok { writeJSON(w,http.StatusNotFound,map[string]string{"error":"state not found"}); return }; writeJSON(w,http.StatusOK,state)
}

func writeJSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
func requestLog(next http.Handler) http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){start:=time.Now();next.ServeHTTP(w,r);log.Printf("%s %s %s",r.Method,r.URL.Path,time.Since(start))})}

func RunHTTP(ctx context.Context, addr string) error {
	if addr=="" { addr="127.0.0.1:8080" }
	srv:=&http.Server{Addr:addr,Handler:NewAPIServer().Handler(),ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:15*time.Second,IdleTimeout:60*time.Second}
	go func(){<-ctx.Done(); shutdownCtx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel();_ = srv.Shutdown(shutdownCtx)}()
	log.Printf("FTN SER AI listening on %s",addr)
	if err:=srv.ListenAndServe();err!=nil && !errors.Is(err,http.ErrServerClosed){return err};return nil
}
