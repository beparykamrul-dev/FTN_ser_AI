package backend

import "testing"

func TestManualEndpointIdentity(t *testing.T) {
    s := NewEndpointIdentityStore()
    e := EndpointIdentity{ID:"node-1", Name:"FTN Node", IP:"192.0.2.10", MAC:"02:00:00:00:00:10", ASN:64510, Country:"BD", Authorized:true}
    if err := s.Add(e); err != nil { t.Fatal(err) }
    got, ok := s.Get("node-1")
    if !ok || got.IP != e.IP || got.MAC != e.MAC { t.Fatalf("unexpected identity: %#v", got) }
}

func TestManualEndpointIdentityRejectsInvalidIPMAC(t *testing.T) {
    s := NewEndpointIdentityStore()
    if err := s.Add(EndpointIdentity{ID:"bad", Name:"bad", IP:"not-an-ip", MAC:"not-a-mac", ASN:64510, Country:"BD"}); err == nil { t.Fatal("expected validation error") }
}
