package backend

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestObservationPollerPublishes(t *testing.T) {
	var count atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	p, err := NewObservationPoller(time.Millisecond, func(context.Context) (ServiceEvent, error) {
		return ServiceEvent{ID: "poll", Kind: "health", At: time.Now()}, nil
	}, func(ServiceEvent) error {
		count.Add(1)
		return nil
	})
	if err != nil { t.Fatal(err) }
	_ = p.Start(ctx)
	if count.Load() == 0 { t.Fatal("poller did not publish an observation") }
}

func TestObservationPollerRejectsInvalidConfig(t *testing.T) {
	if _, err := NewObservationPoller(0, nil, nil); err == nil { t.Fatal("expected invalid configuration error") }
}
