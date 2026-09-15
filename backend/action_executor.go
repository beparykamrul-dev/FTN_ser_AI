package backend

import (
 "context"
 "encoding/json"
 "fmt"
 "strings"
 "sync"
 "time"
)

type ActionDefinition struct { Kind string `json:"kind"`; Capability string `json:"capability"`; Risk string `json:"risk"`; ApprovalRequired bool `json:"approval_required"` }
type ActionAudit struct { ActionID string `json:"action_id"`; Kind string `json:"kind"`; Target string `json:"target"`; Status string `json:"status"`; Actor string `json:"actor,omitempty"`; Time time.Time `json:"time"`; Detail string `json:"detail,omitempty"` }
type ActionExecutor struct { mu sync.RWMutex; audit []ActionAudit; agents *AgentHub; adapters *AdapterEngine }
func NewActionExecutor(agents *AgentHub, adapters *AdapterEngine) *ActionExecutor { return &ActionExecutor{audit:make([]ActionAudit,0,256),agents:agents,adapters:adapters} }
func ActionCatalog() []ActionDefinition { return []ActionDefinition{
 {Kind:"device.health",Capability:"monitoring",Risk:"low"},{Kind:"device.snapshot",Capability:"monitoring",Risk:"low"},{Kind:"network.latency",Capability:"monitoring",Risk:"low"},{Kind:"network.route",Capability:"monitoring",Risk:"low"},{Kind:"dns.query",Capability:"dns_mesh",Risk:"low"},{Kind:"dns.validate",Capability:"dns_mesh",Risk:"low"},{Kind:"ddns.validate",Capability:"dns_mesh",Risk:"low"},{Kind:"provider.health",Capability:"providers",Risk:"low"},{Kind:"provider.zone.list",Capability:"providers",Risk:"low"},{Kind:"provider.dns.records.list",Capability:"providers",Risk:"low"},{Kind:"database.health",Capability:"database",Risk:"low"},{Kind:"monitoring.refresh",Capability:"monitoring",Risk:"low"},
} }
func (e *ActionExecutor) Execute(ctx context.Context, action WebAction) WebAction {
 action.UpdatedAt=time.Now().UTC(); def,ok:=actionDefinition(action.Kind)
 if !ok { action.Status="rejected"; action.Error="unsupported action"; e.record(action,"rejected","action not in capability catalog"); return action }
 if strings.TrimSpace(action.Target)=="" { action.Status="rejected"; action.Error="target required"; e.record(action,"rejected","target required"); return action }
 if def.ApprovalRequired { action.Status="pending_approval"; e.record(action,"pending_approval","approval required before execution"); return action }
 if strings.HasPrefix(action.Kind,"provider.") && e.adapters!=nil { if err:=e.executeProvider(ctx,&action); err!=nil { action.Status="failed"; action.Error=err.Error(); e.record(action,"failed","provider operation failed"); return action }; action.Status="completed"; e.record(action,"completed","typed provider operation executed server-side"); return action }
 if isAgentAction(action.Kind) && e.agents!=nil { payload,_:=json.Marshal(map[string]any{"kind":action.Kind,"target":action.Target,"action_id":action.ID}); err:=e.agents.Send(action.Target,AgentMessage{Type:"action.request",ID:action.ID,DeviceID:action.Target,Timestamp:time.Now().UTC(),Payload:payload}); if err!=nil { action.Status="queued"; action.Error="agent not connected; action retained for later delivery"; e.record(action,"queued","typed action queued; no arbitrary shell execution"); return action }; action.Status="dispatched"; e.record(action,"dispatched","typed action dispatched over authenticated agent WebSocket"); return action }
 action.Status="accepted"; e.record(action,"accepted","typed action accepted; concrete operation adapter pending"); return action
}
func (e *ActionExecutor) executeProvider(ctx context.Context, action *WebAction) error {
 switch action.Kind {
 case "provider.health": return e.adapters.ExecuteHealth(ctx,strings.ToLower(action.Target))
 case "provider.zone.list":
  if strings.ToLower(strings.TrimSpace(action.Target))!="cloudflare" { return fmt.Errorf("provider zone listing currently supports cloudflare only") }
  a,ok:=e.adapters.Get("cloudflare"); if !ok{return fmt.Errorf("cloudflare adapter not found")}; cf,ok:=a.(*HTTPProviderAdapter); if !ok{return fmt.Errorf("cloudflare adapter type mismatch")}; v,err:=cf.CloudflareZones(ctx); if err!=nil{return err}; b,_:=json.Marshal(v); action.Result=b; return nil
 case "provider.dns.records.list":
  p:=strings.SplitN(action.Target,"|",2); if len(p)!=2||strings.ToLower(strings.TrimSpace(p[0]))!="cloudflare" {return fmt.Errorf("target must be cloudflare|zone.example")}; zone:=strings.TrimSuffix(strings.ToLower(strings.TrimSpace(p[1])),"."); if !isManagedProviderZoneRead(zone){return fmt.Errorf("zone is not allowed for provider read: %s",zone)}; a,ok:=e.adapters.Get("cloudflare"); if !ok{return fmt.Errorf("cloudflare adapter not found")}; cf,ok:=a.(*HTTPProviderAdapter); if !ok{return fmt.Errorf("cloudflare adapter type mismatch")}; v,err:=cf.CloudflareRecords(ctx,zone); if err!=nil{return err}; b,_:=json.Marshal(v); action.Result=b; return nil
 default:return fmt.Errorf("unsupported provider operation: %s",action.Kind)
 }
}
func isManagedProviderZoneRead(zone string) bool { for _,z:=range AuthorityZones(){if z.Mode==AuthorityCloudflarePrimary&&z.ProviderID=="cloudflare"&&z.ReadOnlyMirror&&strings.EqualFold(z.Zone,zone){return true}}; return false }
func isAgentAction(kind string) bool { return strings.HasPrefix(kind,"device.")||strings.HasPrefix(kind,"network.") }
func actionDefinition(kind string)(ActionDefinition,bool){ for _,a:=range ActionCatalog(){if a.Kind==kind{return a,true}};return ActionDefinition{},false }
func (e *ActionExecutor) record(a WebAction,status,detail string){ e.mu.Lock();defer e.mu.Unlock();e.audit=append([]ActionAudit{{ActionID:a.ID,Kind:a.Kind,Target:a.Target,Status:status,Actor:a.RequestedBy,Time:time.Now().UTC(),Detail:detail}},e.audit...);if len(e.audit)>256{e.audit=e.audit[:256]} }
func (e *ActionExecutor) Audit() []ActionAudit { e.mu.RLock();defer e.mu.RUnlock();return append([]ActionAudit(nil),e.audit...) }
