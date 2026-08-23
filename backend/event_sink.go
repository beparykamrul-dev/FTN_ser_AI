package backend

import (
	"fmt"
	"sync"
)

// EventSink is the stable boundary for persistent/streaming backends.
// Implementations can target PostgreSQL, CockroachDB, Kafka, or other approved stores.
type EventSink interface {
	Append(ServiceEvent) error
}

type MemoryEventSink struct {
	mu     sync.RWMutex
	events []ServiceEvent
}

func NewMemoryEventSink() *MemoryEventSink { return &MemoryEventSink{} }

func (s *MemoryEventSink) Append(e ServiceEvent) error {
	if e.ID == "" || e.Kind == "" || e.At.IsZero() {
		return fmt.Errorf("invalid event")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
	return nil
}

func (s *MemoryEventSink) Snapshot() []ServiceEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]ServiceEvent, len(s.events))
	copy(out, s.events)
	return out
}

// Persistent adapters should implement EventSink without changing callers.
type EventPublisher struct { sink EventSink }

func NewEventPublisher(sink EventSink) (*EventPublisher, error) {
	if sink == nil { return nil, fmt.Errorf("event sink is required") }
	return &EventPublisher{sink: sink}, nil
}

func (p *EventPublisher) Publish(e ServiceEvent) error { return p.sink.Append(e) }
