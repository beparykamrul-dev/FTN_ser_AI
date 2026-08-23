package backend

import (
	"fmt"
	"sync"
	"time"
)

type CertificateState struct {
	ID        string
	Subject   string
	Issuer    string
	ExpiresAt time.Time
	Revoked   bool
}

type GatewayPolicy struct {
	ID             string
	Enabled        bool
	MaxConnections int
	MaxMbps        int
	RequireMTLS    bool
}

type GatewayDecision struct {
	Allowed bool
	Reason  string
}

type PKIGatewayStore struct {
	mu      sync.RWMutex
	certs   map[string]CertificateState
	policies map[string]GatewayPolicy
}

func NewPKIGatewayStore() *PKIGatewayStore {
	return &PKIGatewayStore{certs: map[string]CertificateState{}, policies: map[string]GatewayPolicy{}}
}

func (s *PKIGatewayStore) UpsertCertificate(c CertificateState) error {
	if c.ID == "" || c.Subject == "" || c.ExpiresAt.IsZero() {
		return fmt.Errorf("invalid certificate state")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.certs[c.ID] = c
	return nil
}

func (s *PKIGatewayStore) PutGatewayPolicy(p GatewayPolicy) error {
	if p.ID == "" || p.MaxConnections < 0 || p.MaxMbps < 0 {
		return fmt.Errorf("invalid gateway policy")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[p.ID] = p
	return nil
}

func (s *PKIGatewayStore) EvaluateGateway(policyID, certID string, now time.Time) GatewayDecision {
	s.mu.RLock()
	p, pok := s.policies[policyID]
	c, cok := s.certs[certID]
	s.mu.RUnlock()
	if !pok || !p.Enabled {
		return GatewayDecision{Reason: "gateway policy disabled or missing"}
	}
	if !cok || c.Revoked || !c.ExpiresAt.After(now) {
		return GatewayDecision{Reason: "certificate invalid or expired"}
	}
	if p.RequireMTLS && c.Issuer == "" {
		return GatewayDecision{Reason: "mTLS trust identity missing"}
	}
	return GatewayDecision{Allowed: true, Reason: "policy and certificate valid"}
}
