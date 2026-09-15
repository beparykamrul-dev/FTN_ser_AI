package backend

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

type ActionDefinition struct { Kind string `json:"kind"`; Capability string `json:"capability"`; Risk string `json:"risk"`; ApprovalRequired bool `json:"approval_required"` }
type ActionAudit struct { ActionID string `json:"action_id"`; Kind string `json:"kind"`; Target string `json:"target"`; Status string `json:"status"`; Actor string `json:"actor,omitempty"`; Time time.Time `json:"time"`; Detail string `json:"detail,omitempty"` }
type ActionExecutor struct { mu sync.RWMutex; audit []ActionAudit; agents *AgentHub }
func NewActionExecutor(agents *AgentHub) *ActionExecutor { return &ActionExecutor{audit:make([]ActionAudit,0,256),agents:agents} }
func ActionCatalog() []ActionDefinition { return []ActionDefinition{
	{Kind:"device.health",Capability:"monitoring",Risk:"low"},{Kind:"device.snapshot",Capability:"monitoring",Risk:"low"},
	{Kind:"network.latency",Capability:"monitoring",Risk:"low"},{Kind:"network.route",Capability:"monitoring",Risk:"low"},
	{Kind:"dns.query",Capability:"dns_mesh",Risk:"low"},{Kind:"dns.validate",Capability:"dns_mesh",Risk:"low"},
	{Kind:"ddns.validate",Capability:"dns_mesh",Risk:"low"},{Kind:"provider.health",Capability:"providers",Risk:"low"},
	{Kind:"database.health",Capability:"database",Risk:"low"},{Kind:"monitoring.refresh",Capability:"monitoring",Risk:"low"},
} }
func (e *ActionExecutor) Execute(ctx context.Context, action WebAction) WebAction {
	action.UpdatedAt=time.Now().UTC(); def,ok:=actionDefinition(action.Kind)
	if !ok { action.Status="rejected"; action.Error="unsupported action"; e.record(action,"rejected","action not in capability catalog"); return action }
	if strings.TrimSpace(action.Target)=="" { action.Status="rejected"; action.Error="target required"; e.record(action,"rejected","target required"); return action }
	if def.ApprovalRequired { action.Status="pending_approval"; e.record(action,"pending_approval","approval required before execution"); return action }
	if isAgentAction(action.Kind) && e.agents!=nil {
		payload,_:=json.Marshal(map[string]any{"kind":action.Kind,"target":action.Target,"action_id":action.ID})
		err:=e.agents.Send(action.Target,AgentMessage{Type:"action.request",ID:action.ID,DeviceID:action.Target,Timestamp:time.Now().UTC(),Payload:payload})
		if err!=nil { action.Status="queued"; action.Error="agent not connected; action retained for later delivery"; e.record(action,"queued","typed action queued; no arbitrary shell execution"); return action }
		action.Status="dispatched"; e.record(action,"dispatched","typed action dispatched over authenticated agent WebSocket"); return action
	}
	action.Status="accepted"; e.record(action,"accepted","typed action accepted; concrete provider/database/DNS adapter execution pending"); _=ctx; return action
}
func isAgentAction(kind string) bool { return strings.HasPrefix(kind,"device.")||strings.HasPrefix(kind,"network.") }
func actionDefinition(kind string)(ActionDefinition,bool){ for _,a:=range ActionCatalog(){if a.Kind==kind{return a,true}};return ActionDefinition{},false }
func (e *ActionExecutor) record(a WebAction,status,detail string){ e.mu.Lock();defer e.mu.Unlock();e.audit=append([]ActionAudit{{ActionID:a.ID,Kind:a.Kind,Target:a.Target,Status:status,Actor:a.RequestedBy,Time:time.Now().UTC(),Detail:detail}},e.audit...);if len(e.audit)>256{e.audit=e.audit[:256]} }
func (e *ActionExecutor) Audit() []ActionAudit { e.mu.RLock();defer e.mu.RUnlock();return append([]ActionAudit(nil),e.audit...) }
