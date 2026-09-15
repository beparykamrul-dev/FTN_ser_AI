package backend

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type AgentMessage struct {
	Type string `json:"type"`
	ID string `json:"id,omitempty"`
	DeviceID string `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type AgentSession struct {
	DeviceID string `json:"device_id"`
	Connected bool `json:"connected"`
	ConnectedAt time.Time `json:"connected_at,omitempty"`
	LastSeen time.Time `json:"last_seen,omitempty"`
	RemoteAddr string `json:"remote_addr,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
	conn *websocket.Conn
	mu sync.Mutex
}

type AgentHub struct { mu sync.RWMutex; sessions map[string]*AgentSession }
func NewAgentHub() *AgentHub { return &AgentHub{sessions: map[string]*AgentSession{}} }
func (h *AgentHub) List() []AgentSession { h.mu.RLock(); defer h.mu.RUnlock(); out:=make([]AgentSession,0,len(h.sessions)); for _,s:=range h.sessions { c:=*s; c.conn=nil; out=append(out,c) }; return out }
func (h *AgentHub) Register(id,remote string,conn *websocket.Conn)*AgentSession { h.mu.Lock(); defer h.mu.Unlock(); if old:=h.sessions[id]; old!=nil { old.mu.Lock(); if old.conn!=nil {_=old.conn.Close()}; old.mu.Unlock() }; now:=time.Now().UTC(); s:=&AgentSession{DeviceID:id,Connected:true,ConnectedAt:now,LastSeen:now,RemoteAddr:remote,conn:conn}; h.sessions[id]=s; return s }
func (h *AgentHub) Unregister(id string,conn *websocket.Conn) { h.mu.Lock(); defer h.mu.Unlock(); s:=h.sessions[id]; if s==nil || s.conn!=conn{return}; s.mu.Lock(); s.Connected=false; s.LastSeen=time.Now().UTC(); s.conn=nil; s.mu.Unlock() }
func (h *AgentHub) Send(id string,msg AgentMessage) error { h.mu.RLock(); s:=h.sessions[id]; h.mu.RUnlock(); if s==nil || !s.Connected || s.conn==nil{return http.ErrNotSupported}; s.mu.Lock(); defer s.mu.Unlock(); return s.conn.WriteJSON(msg) }

func validAgentAuth(r *http.Request) bool {
	secret:=os.Getenv("FTN_AGENT_SECRET"); if secret=="" { return false }
	got:=strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "))
	if got=="" { return false }
	return subtle.ConstantTimeCompare([]byte(got),[]byte(secret))==1
}

func (s *APIServer) agentsList(w http.ResponseWriter,r *http.Request) { if r.Method!=http.MethodGet { w.Header().Set("Allow",http.MethodGet); writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"}); return }; writeJSON(w,http.StatusOK,map[string]any{"agents":s.agents.List()}) }
func (s *APIServer) agentWS(w http.ResponseWriter,r *http.Request) {
	if !validAgentAuth(r) { writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"invalid agent credentials"}); return }
	deviceID:=strings.TrimSpace(r.Header.Get("X-FTN-Device-ID")); if deviceID=="" || len(deviceID)>128 { writeJSON(w,http.StatusBadRequest,map[string]string{"error":"invalid X-FTN-Device-ID"}); return }
	upgrader:=websocket.Upgrader{CheckOrigin:func(_ *http.Request)bool{return false},ReadBufferSize:64*1024,WriteBufferSize:64*1024}
	conn,err:=upgrader.Upgrade(w,r,nil); if err!=nil{return}; session:=s.agents.Register(deviceID,r.RemoteAddr,conn)
	defer func(){s.agents.Unregister(deviceID,conn);_=conn.Close()}()
	_=conn.WriteJSON(AgentMessage{Type:"hello_ack",DeviceID:deviceID,Timestamp:time.Now().UTC()})
	conn.SetReadLimit(1<<20); _=conn.SetReadDeadline(time.Now().Add(90*time.Second)); conn.SetPongHandler(func(string)error{_=conn.SetReadDeadline(time.Now().Add(90*time.Second));return nil})
	for { var msg AgentMessage; if err:=conn.ReadJSON(&msg);err!=nil{return}; if msg.DeviceID==""{msg.DeviceID=deviceID}; session.mu.Lock();session.LastSeen=time.Now().UTC();session.mu.Unlock(); if msg.Type=="heartbeat"{_=conn.WriteJSON(AgentMessage{Type:"heartbeat_ack",DeviceID:deviceID,Timestamp:time.Now().UTC()})} }
}
