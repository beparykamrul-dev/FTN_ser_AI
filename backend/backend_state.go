package backend

import (
	"errors"
	"sync"
	"time"
)

// ServiceState is the minimal backend state consumed by the AI reconciler.
type ServiceState struct {
	NodeID        string    `json:"node_id"`
	ServiceID     string    `json:"service_id"`
	Healthy       bool      `json:"healthy"`
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryPercent float64   `json:"memory_percent"`
	LatencyMS     float64   `json:"latency_ms"`
	PacketLoss    float64   `json:"packet_loss"`
	ObservedAt    time.Time `json:"observed_at"`
}

type StateStore struct {
	mu     sync.RWMutex
	states map[string]ServiceState
}

func NewStateStore() *StateStore {
	return &StateStore{states: make(map[string]ServiceState)}
}

func (s *StateStore) Upsert(state ServiceState) error {
	if state.NodeID == "" || state.ServiceID == "" {
		return errors.New("node_id and service_id are required")
	}
	if state.CPUPercent < 0 || state.CPUPercent > 100 || state.MemoryPercent < 0 || state.MemoryPercent > 100 {
		return errors.New("resource percentages must be between 0 and 100")
	}
	if state.LatencyMS < 0 || state.PacketLoss < 0 || state.PacketLoss > 100 {
		return errors.New("invalid latency or packet loss")
	}
	if state.ObservedAt.IsZero() {
		state.ObservedAt = time.Now().UTC()
	}
	s.mu.Lock()
	s.states[state.NodeID+":"+state.ServiceID] = state
	s.mu.Unlock()
	return nil
}

func (s *StateStore) Get(nodeID, serviceID string) (ServiceState, bool) {
	s.mu.RLock()
	state, ok := s.states[nodeID+":"+serviceID]
	s.mu.RUnlock()
	return state, ok
}
