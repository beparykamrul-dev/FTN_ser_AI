package backend

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type FlowRecord struct { ExporterID string `json:"exporter_id"`; Protocol string `json:"protocol"`; SrcIP string `json:"src_ip"`; DstIP string `json:"dst_ip"`; SrcPort uint16 `json:"src_port"`; DstPort uint16 `json:"dst_port"`; Packets uint64 `json:"packets"`; Bytes uint64 `json:"bytes"`; StartedAt time.Time `json:"started_at"`; EndedAt time.Time `json:"ended_at"` }
type FlowExporterRegistry struct { mu sync.RWMutex; exporters map[string]bool }
func NewFlowExporterRegistry() *FlowExporterRegistry { return &FlowExporterRegistry{exporters: make(map[string]bool)} }
func (r *FlowExporterRegistry) Authorize(id string) { r.mu.Lock(); defer r.mu.Unlock(); r.exporters[id]=true }
func (r *FlowExporterRegistry) Revoke(id string) { r.mu.Lock(); defer r.mu.Unlock(); delete(r.exporters,id) }
func (r *FlowExporterRegistry) Accept(record FlowRecord) error { if record.ExporterID=="" {return errors.New("flow exporter is required")}; r.mu.RLock(); ok:=r.exporters[record.ExporterID]; r.mu.RUnlock(); if !ok{return fmt.Errorf("unauthorized flow exporter: %s",record.ExporterID)}; if record.Bytes==0 && record.Packets==0{return errors.New("flow record contains no traffic counters")}; return nil }
