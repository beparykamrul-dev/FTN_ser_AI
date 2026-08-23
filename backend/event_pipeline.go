package backend

import (
	"fmt"
	"sync"
	"time"
)

type ServiceEvent struct {
	ID        string
	Kind      string
	NodeID    string
	ServiceID string
	Value     float64
	At        time.Time
}

type EventPipeline struct {
	mu     sync.RWMutex
	events []ServiceEvent
}

func NewEventPipeline() *EventPipeline { return &EventPipeline{events: make([]ServiceEvent, 0, 128)} }

func (p *EventPipeline) Publish(e ServiceEvent) error {
	if e.ID == "" || e.Kind == "" || e.At.IsZero() {
		return fmt.Errorf("invalid service event")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, e)
	return nil
}

func (p *EventPipeline) Snapshot() []ServiceEvent {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]ServiceEvent, len(p.events))
	copy(out, p.events)
	return out
}
