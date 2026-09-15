package backend

import (
    "context"
    "fmt"
    "strings"
    "sync"
    "time"
)

type ActionDefinition struct {
    Kind string `json:"kind"`
    Capability string `json:"capability"`
    Risk string `json:"risk"`
    ApprovalRequired bool `json:"approval_required"`
}

type ActionAudit struct {
    ActionID string `json:"action_id"`
    Kind string `json:"kind"`
    Target string `json:"target"`
    Status string `json:"status"`
    Actor string `json:"actor,omitempty"`
    Time time.Time `json:"time"`
    Detail string `json:"detail,omitempty"`
}

type ActionExecutor struct { mu sync.RWMutex; audit []ActionAudit }

func NewActionExecutor() *ActionExecutor { return &ActionExecutor{audit: make([]ActionAudit,0,256)} }

func ActionCatalog() []ActionDefinition {
    return []ActionDefinition{
        {Kind:"device.health",Capability:"monitoring",Risk:"low"},
        {Kind:"device.snapshot",Capability:"monitoring",Risk:"low"},
        {Kind:"network.latency",Capability:"monitoring",Risk:"low"},
        {Kind:"network.route",Capability:"monitoring",Risk:"low"},
        {Kind:"dns.query",Capability:"dns_mesh",Risk:"low"},
        {Kind:"dns.validate",Capability:"dns_mesh",Risk:"low"},
        {Kind:"ddns.validate",Capability:"dns_mesh",Risk:"low"},
        {Kind:"provider.health",Capability:"providers",Risk:"low"},
        {Kind:"database.health",Capability:"database",Risk:"low"},
        {Kind:"monitoring.refresh",Capability:"monitoring",Risk:"low"},
    }
}

func (e *ActionExecutor) Execute(ctx context.Context, action WebAction) WebAction {
    _ = ctx
    action.UpdatedAt=time.Now().UTC()
    if !allowedAction(action.Kind) { action.Status="rejected"; action.Error="unsupported action"; e.record(action,"rejected","action not in capability catalog"); return action }
    if strings.TrimSpace(action.Target)=="" { action.Status="rejected"; action.Error="target required"; e.record(action,"rejected","target required"); return action }
    action.Status="completed"
    e.record(action,"completed","validated web operation; provider/agent execution adapter can consume this action")
    return action
}

func allowedAction(kind string) bool { for _,a:=range ActionCatalog(){if a.Kind==kind{return true}}; return false }
func (e *ActionExecutor) record(a WebAction,status,detail string){e.mu.Lock();defer e.mu.Unlock();e.audit=append([]ActionAudit{{ActionID:a.ID,Kind:a.Kind,Target:a.Target,Status:status,Actor:a.RequestedBy,Time:time.Now().UTC(),Detail:detail}},e.audit...);if len(e.audit)>256{e.audit=e.audit[:256]}}
func (e *ActionExecutor) Audit() []ActionAudit {e.mu.RLock();defer e.mu.RUnlock();return append([]ActionAudit(nil),e.audit...)}
func actionError(kind string) error {return fmt.Errorf("unsupported action: %s",kind)}
