package backend

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// AIAction is an auditable action proposed by FTN AI.
type AIAction struct {
	ID string `json:"id"`
	Kind string `json:"kind"`
	Target string `json:"target"`
	Reason string `json:"reason"`
	Parameters map[string]string `json:"parameters,omitempty"`
	RequiresAuth bool `json:"requires_auth"`
	Approved bool `json:"approved"`
	CreatedAt time.Time `json:"created_at"`
}

type NodeObservation struct {
	NodeID string `json:"node_id"`
	CPUPercent float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	LatencyMS float64 `json:"latency_ms"`
	PacketLoss float64 `json:"packet_loss"`
	Healthy bool `json:"healthy"`
	TrafficInBPS uint64 `json:"traffic_in_bps"`
	TrafficOutBPS uint64 `json:"traffic_out_bps"`
	ObservedAt time.Time `json:"observed_at"`
}

func Decide(obs NodeObservation) (AIAction, bool) {
	if obs.NodeID == "" || !obs.Healthy { return AIAction{}, false }
	if obs.PacketLoss >= 5 || obs.LatencyMS >= 250 || obs.CPUPercent >= 90 || obs.MemoryPercent >= 90 {
		return AIAction{ID: "rebalance-" + obs.NodeID, Kind: "service.rebalance", Target: obs.NodeID, Reason: "node crossed FTN health/capacity policy threshold", RequiresAuth: true, CreatedAt: time.Now().UTC()}, true
	}
	return AIAction{}, false
}

type ActionStore struct { mu sync.RWMutex; items map[string]AIAction }
func NewActionStore() *ActionStore { return &ActionStore{items: make(map[string]AIAction)} }
func (s *ActionStore) Put(a AIAction) error { if a.ID == "" { return errors.New("action id is required") }; s.mu.Lock(); defer s.mu.Unlock(); if _, exists := s.items[a.ID]; exists { return errors.New("action already exists") }; s.items[a.ID] = a; return nil }
func (s *ActionStore) List() []AIAction { s.mu.RLock(); defer s.mu.RUnlock(); out := make([]AIAction,0,len(s.items)); for _,a := range s.items { out=append(out,a) }; sort.Slice(out,func(i,j int)bool{return out[i].ID<out[j].ID}); return out }
func (s *ActionStore) Approve(id string) error { s.mu.Lock(); defer s.mu.Unlock(); a,ok:=s.items[id]; if !ok{return errors.New("action not found")}; a.Approved=true; s.items[id]=a; return nil }
