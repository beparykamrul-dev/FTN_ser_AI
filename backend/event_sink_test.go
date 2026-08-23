package backend

import (
	"testing"
	"time"
)

func TestEventPublisherUsesSink(t *testing.T) {
	sink := NewMemoryEventSink()
	pub, err := NewEventPublisher(sink)
	if err != nil { t.Fatal(err) }
	e := ServiceEvent{ID: "e1", Kind: "metric", NodeID: "n1", ServiceID: "s1", Value: 42, At: time.Now()}
	if err := pub.Publish(e); err != nil { t.Fatal(err) }
	if got := sink.Snapshot(); len(got) != 1 || got[0].ID != e.ID { t.Fatalf("unexpected events: %#v", got) }
}

func TestEventPublisherRequiresSink(t *testing.T) {
	if _, err := NewEventPublisher(nil); err == nil { t.Fatal("expected missing sink error") }
}
