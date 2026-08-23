package backend

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ObservationPoller struct {
	interval time.Duration
	observe  func(context.Context) (ServiceEvent, error)
	publish  func(ServiceEvent) error
	mu       sync.Mutex
	running  bool
}

func NewObservationPoller(interval time.Duration, observe func(context.Context) (ServiceEvent, error), publish func(ServiceEvent) error) (*ObservationPoller, error) {
	if interval <= 0 || observe == nil || publish == nil {
		return nil, fmt.Errorf("invalid poller configuration")
	}
	return &ObservationPoller{interval: interval, observe: observe, publish: publish}, nil
}

func (p *ObservationPoller) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return fmt.Errorf("poller already running")
	}
	p.running = true
	p.mu.Unlock()
	defer func() { p.mu.Lock(); p.running = false; p.mu.Unlock() }()

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			e, err := p.observe(ctx)
			if err != nil {
				continue
			}
			_ = p.publish(e)
		}
	}
}
