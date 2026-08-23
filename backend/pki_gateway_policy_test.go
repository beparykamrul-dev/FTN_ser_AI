package backend

import (
	"testing"
	"time"
)

func TestGatewayRejectsExpiredCertificate(t *testing.T) {
	s := NewPKIGatewayStore()
	now := time.Now()
	if err := s.UpsertCertificate(CertificateState{ID: "c1", Subject: "node-1", Issuer: "ftn-ca", ExpiresAt: now.Add(-time.Minute)}); err != nil { t.Fatal(err) }
	if err := s.PutGatewayPolicy(GatewayPolicy{ID: "g1", Enabled: true, RequireMTLS: true}); err != nil { t.Fatal(err) }
	if got := s.EvaluateGateway("g1", "c1", now); got.Allowed { t.Fatal("expired certificate must be rejected") }
}

func TestGatewayAllowsValidCertificate(t *testing.T) {
	s := NewPKIGatewayStore()
	now := time.Now()
	if err := s.UpsertCertificate(CertificateState{ID: "c1", Subject: "node-1", Issuer: "ftn-ca", ExpiresAt: now.Add(time.Hour)}); err != nil { t.Fatal(err) }
	if err := s.PutGatewayPolicy(GatewayPolicy{ID: "g1", Enabled: true, RequireMTLS: true}); err != nil { t.Fatal(err) }
	if got := s.EvaluateGateway("g1", "c1", now); !got.Allowed { t.Fatalf("valid certificate rejected: %s", got.Reason) }
}
