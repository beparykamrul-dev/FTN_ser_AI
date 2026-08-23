package backend

import (
	"testing"
	"time"
)

func TestEventPipelinePublishesAndSnapshots(t *testing.T) {
	p := NewEventPipeline()
	if err := p.Publish(ServiceEvent{ID: "e1", Kind: "health", NodeID: "n1", ServiceID: "s1", At: time.Now()}); err != nil { t.Fatal(err) }
	got := p.Snapshot()
	if len(got) != 1 || got[0].ID != "e1" { t.Fatalf("unexpected snapshot: %#v", got) }
}

func TestEventPipelineRejectsInvalidEvent(t *testing.T) {
	p := NewEventPipeline()
	if err := p.Publish(ServiceEvent{ID: "", Kind: "health", At: time.Now()}); err == nil { t.Fatal("expected invalid event error") }
}
