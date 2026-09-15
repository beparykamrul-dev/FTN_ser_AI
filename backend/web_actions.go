package backend

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WebAction struct {
	ID string `json:"id"`
	Kind string `json:"kind"`
	Target string `json:"target"`
	Status string `json:"status"`
	RequestedBy string `json:"requested_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Error string `json:"error,omitempty"`
}

type ActionRequest struct {
	Kind string `json:"kind"`
	Target string `json:"target"`
	RequestedBy string `json:"requested_by,omitempty"`
	Confirm bool `json:"confirm"`
}

type WebActionStore struct { mu sync.RWMutex; items []WebAction }
func NewWebActionStore() *WebActionStore { return &WebActionStore{items: make([]WebAction,0,128)} }
func (s *WebActionStore) Create(in ActionRequest) (WebAction,error) {
	if strings.TrimSpace(in.Kind)=="" || strings.TrimSpace(in.Target)=="" { return WebAction{},http.ErrMissingFile }
	if !in.Confirm { return WebAction{},http.ErrNotSupported }
	idb:=make([]byte,8); if _,err:=rand.Read(idb);err!=nil{return WebAction{},err}; id:=hex.EncodeToString(idb)
	now:=time.Now().UTC(); a:=WebAction{ID:id,Kind:strings.TrimSpace(in.Kind),Target:strings.TrimSpace(in.Target),Status:"accepted",RequestedBy:strings.TrimSpace(in.RequestedBy),CreatedAt:now,UpdatedAt:now}
	s.mu.Lock(); s.items=append([]WebAction{a},s.items...); if len(s.items)>128{s.items=s.items[:128]}; s.mu.Unlock(); return a,nil
}
func (s *WebActionStore) List() []WebAction { s.mu.RLock(); defer s.mu.RUnlock(); return append([]WebAction(nil),s.items...) }

func (s *APIServer) webActions(w http.ResponseWriter,r *http.Request) {
	if r.Method==http.MethodGet { writeJSON(w,http.StatusOK,map[string]any{"actions":s.actions.List()}); return }
	if r.Method!=http.MethodPost { w.Header().Set("Allow","GET, POST"); writeJSON(w,http.StatusMethodNotAllowed,map[string]string{"error":"method not allowed"}); return }
	var in ActionRequest; dec:=json.NewDecoder(http.MaxBytesReader(w,r.Body,64<<10)); if err:=dec.Decode(&in);err!=nil { writeJSON(w,http.StatusBadRequest,map[string]string{"error":"invalid JSON"}); return }
	a,err:=s.actions.Create(in); if err!=nil { writeJSON(w,http.StatusBadRequest,map[string]string{"error":"action requires kind, target and confirm=true"}); return }
	writeJSON(w,http.StatusAccepted,map[string]any{"ok":true,"action":a})
}
