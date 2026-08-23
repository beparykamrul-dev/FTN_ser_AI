package backend

import (
 "net/netip"
 "testing"
)

func TestProviderRegistryFiltersDNSByCountry(t *testing.T) { r:=NewProviderRegistry(); if err:=r.Register(ProviderEndpoint{ID:"dns-bd-1",Provider:"example-dns",Class:ProviderDNS,Country:"BD",Address:netip.MustParseAddr("192.0.2.1"),ASN:64500,Official:true,Authorized:true,Enabled:true}); err!=nil {t.Fatal(err)}; if err:=r.Register(ProviderEndpoint{ID:"cache-bd-1",Provider:"example-cache",Class:ProviderCache,Country:"BD",Address:netip.MustParseAddr("192.0.2.2"),ASN:64501,Official:false,Authorized:true,Enabled:true}); err!=nil {t.Fatal(err)}; got:=r.List(ProviderDNS,"BD"); if len(got)!=1 || got[0].ID!="dns-bd-1" {t.Fatalf("unexpected providers: %#v",got)} }

func TestProviderRegistryRejectsUnauthorized(t *testing.T) { r:=NewProviderRegistry(); if err:=r.Register(ProviderEndpoint{ID:"x",Provider:"x",Class:ProviderDNS,Address:netip.MustParseAddr("192.0.2.9"),ASN:64510,Authorized:false}); err==nil {t.Fatal("expected unauthorized endpoint rejection")} }
