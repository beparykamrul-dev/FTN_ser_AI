package backend

import (
    "fmt"
    "strings"
    "sync"
)

type RouterVendor string
const (
    RouterMikroTik RouterVendor = "mikrotik"
    RouterFTN RouterVendor = "ftn"
)

type RouterProfile struct {
    ID string
    Name string
    Vendor RouterVendor
    Address string
    APIVersion string
    Authorized bool
    Enabled bool
    Capabilities []PhaseFabricCapability
}

type CoreRouterRegistry struct { mu sync.RWMutex; routers map[string]RouterProfile }
func NewCoreRouterRegistry() *CoreRouterRegistry { return &CoreRouterRegistry{routers:map[string]RouterProfile{}} }

func (r *CoreRouterRegistry) Upsert(p RouterProfile) error {
    if p.ID=="" || p.Name=="" || p.Address=="" || !p.Authorized { return fmt.Errorf("invalid or unauthorized router profile") }
    p.Address = strings.TrimSpace(p.Address)
    if p.Vendor!=RouterMikroTik && p.Vendor!=RouterFTN { return fmt.Errorf("unsupported router vendor") }
    r.mu.Lock(); r.routers[p.ID]=p; r.mu.Unlock(); return nil
}
func (r *CoreRouterRegistry) Get(id string)(RouterProfile,bool){r.mu.RLock();defer r.mu.RUnlock();p,ok:=r.routers[id];return p,ok}
