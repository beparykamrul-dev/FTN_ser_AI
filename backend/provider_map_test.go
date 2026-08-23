package backend

import "testing"

func TestProviderMapFilter(t *testing.T) {
    m := NewProviderMap()
    if err := m.Upsert(ProviderMapPoint{ID:"dns-1", Provider:"FTN DNS", IP:"192.0.2.1", ASN:64501, Country:"BD", ServiceType:"dns", Authorized:true, Latitude:23.8, Longitude:90.4, Healthy:true}); err != nil { t.Fatal(err) }
    if err := m.Upsert(ProviderMapPoint{ID:"cache-1", Provider:"Example Cache", IP:"192.0.2.2", ASN:64502, Country:"BD", ServiceType:"cache", Authorized:false, Latitude:23.9, Longitude:90.5}); err != nil { t.Fatal(err) }
    got := m.Filter("BD", "dns", true)
    if len(got) != 1 || got[0].ID != "dns-1" { t.Fatalf("unexpected map result: %#v", got) }
}

func TestProviderMapRejectsInvalidCoordinates(t *testing.T) {
    m := NewProviderMap()
    if err := m.Upsert(ProviderMapPoint{ID:"x", Provider:"P", IP:"192.0.2.3", ASN:64503, Country:"BD", ServiceType:"dns", Latitude:100, Longitude:90}); err == nil { t.Fatal("expected coordinate validation error") }
}
