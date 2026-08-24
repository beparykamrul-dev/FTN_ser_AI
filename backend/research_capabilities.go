package backend

import (
    "fmt"
    "sync"
)

type ResearchMaturity string
const (
    ResearchEstablished ResearchMaturity = "established"
    ResearchExperimental ResearchMaturity = "experimental"
    ResearchTheoretical ResearchMaturity = "theoretical"
    ResearchSpeculative ResearchMaturity = "speculative"
)

type ResearchCapability struct { ID, Name, Domain string; Maturity ResearchMaturity; ProductionEligible bool; Enabled bool }
type ResearchRegistry struct { mu sync.RWMutex; items map[string]ResearchCapability }
func NewResearchRegistry()*ResearchRegistry{return &ResearchRegistry{items:map[string]ResearchCapability{}}}
func(r *ResearchRegistry) Upsert(x ResearchCapability)error{if x.ID==""||x.Name==""||x.Domain==""{return fmt.Errorf("incomplete research capability")};switch x.Maturity{case ResearchEstablished,ResearchExperimental,ResearchTheoretical,ResearchSpeculative:default:return fmt.Errorf("unsupported research maturity")};if x.Maturity!=ResearchEstablished&&x.ProductionEligible{return fmt.Errorf("non-established research cannot be production eligible")};r.mu.Lock();r.items[x.ID]=x;r.mu.Unlock();return nil}
func(r *ResearchRegistry)Get(id string)(ResearchCapability,bool){r.mu.RLock();defer r.mu.RUnlock();x,ok:=r.items[id];return x,ok}
