package backend

import "sync"

type MonitorStore struct { mu sync.RWMutex; nodes map[string]NodeMetrics; alerts map[string]Alert }
func NewMonitorStore() *MonitorStore { return &MonitorStore{nodes: map[string]NodeMetrics{}, alerts: map[string]Alert{}} }
func (s *MonitorStore) PutNode(v NodeMetrics) { s.mu.Lock(); defer s.mu.Unlock(); s.nodes[v.NodeID]=v }
func (s *MonitorStore) Nodes() []NodeMetrics { s.mu.RLock(); defer s.mu.RUnlock(); out:=make([]NodeMetrics,0,len(s.nodes)); for _,v:=range s.nodes { out=append(out,v) }; return out }
func (s *MonitorStore) PutAlert(v Alert) { s.mu.Lock(); defer s.mu.Unlock(); s.alerts[v.ID]=v }
func (s *MonitorStore) Alerts() []Alert { s.mu.RLock(); defer s.mu.RUnlock(); out:=make([]Alert,0,len(s.alerts)); for _,v:=range s.alerts { out=append(out,v) }; return out }
