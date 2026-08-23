package backend

import "testing"

func TestProviderAgentEnableRequiresAuthorization(t *testing.T) {
    r := NewProviderAgentRegistry()
    if err := r.Register(ProviderAgent{ID:"p1", Provider:"Example", ServiceType:"cdn", Mode:AgentCatalog}); err != nil { t.Fatal(err) }
    if err := r.SetEnabled("p1", true); err == nil { t.Fatal("expected authorization error") }
}

func TestAuthorizedProviderAgentCanBeEnabled(t *testing.T) {
    r := NewProviderAgentRegistry()
    if err := r.Register(ProviderAgent{ID:"p1", Provider:"Example", ServiceType:"cdn", Authorized:true}); err != nil { t.Fatal(err) }
    if err := r.SetEnabled("p1", true); err != nil { t.Fatal(err) }
    got, ok := r.Get("p1")
    if !ok || !got.Enabled || got.Mode != AgentConnected { t.Fatalf("unexpected provider agent: %#v", got) }
}
