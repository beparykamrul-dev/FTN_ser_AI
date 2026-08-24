package backend

import (
    "fmt"
    "sync"
)

type TelemetryCapability string
const (
    TelemetryTurnKey TelemetryCapability = "turnkey"
    TelemetryOpenRC TelemetryCapability = "openrc"
    TelemetryJSDelivr TelemetryCapability = "jsdelivr"
    TelemetryNitefood TelemetryCapability = "nitefood"
    TelemetryVXLAN TelemetryCapability = "vxlan"
    TelemetryIPSec TelemetryCapability = "ipsec-tunnel"
    TelemetryNetSA TelemetryCapability = "netsa"
    TelemetryCERTCCSiLK TelemetryCapability = "certcc-silk"
    TelemetryIPFIX TelemetryCapability = "ipfix"
    TelemetryNetFlowV5 TelemetryCapability = "netflow-v5"
    TelemetryNetFlowV9 TelemetryCapability = "netflow-v9"
    TelemetryPerlPython TelemetryCapability = "perl-python"
    TelemetrySolaris TelemetryCapability = "solaris"
    TelemetryOpenBSD TelemetryCapability = "openbsd"
    TelemetryCygwin TelemetryCapability = "cygwin"
    TelemetryRwfowpack TelemetryCapability = "rwflowpack"
    TelemetrySiLKAnalysis TelemetryCapability = "silk-analysis"
    TelemetryYAFDpacketPlugin TelemetryCapability = "yaf-dpacketplugin"
)

type TelemetryProfile struct { ID string; NodeID string; Capabilities []TelemetryCapability; Authorized bool }
type TelemetryRegistry struct { mu sync.RWMutex; profiles map[string]TelemetryProfile }
func NewTelemetryRegistry()*TelemetryRegistry{return &TelemetryRegistry{profiles:map[string]TelemetryProfile{}}}
func (r *TelemetryRegistry) Upsert(p TelemetryProfile) error { if p.ID==""||p.NodeID==""||!p.Authorized{return fmt.Errorf("invalid or unauthorized telemetry profile")}; if len(p.Capabilities)==0{return fmt.Errorf("telemetry capability required")}; r.mu.Lock(); r.profiles[p.ID]=p; r.mu.Unlock(); return nil }
func (r *TelemetryRegistry) Get(id string)(TelemetryProfile,bool){r.mu.RLock();defer r.mu.RUnlock();p,ok:=r.profiles[id];return p,ok}
