package backend

import (
    "fmt"
    "sync"
)

type LocalServiceClass string
const (
    ServiceRouter LocalServiceClass = "router"
    ServiceMobile LocalServiceClass = "mobile"
    ServiceCCTV LocalServiceClass = "cctv"
    ServiceTV LocalServiceClass = "tv"
    ServiceLocalDNS LocalServiceClass = "local-dns"
    ServiceLocalCache LocalServiceClass = "local-cache"
)

type LocalServiceProfile struct { ID, Name, NodeID string; Class LocalServiceClass; QoSClass string; Authorized bool; Enabled bool }
type LocalFabricRegistry struct { mu sync.RWMutex; services map[string]LocalServiceProfile }
func NewLocalFabricRegistry()*LocalFabricRegistry{return &LocalFabricRegistry{services:map[string]LocalServiceProfile{}}}
func(r *LocalFabricRegistry) Upsert(s LocalServiceProfile)error{if s.ID==""||s.Name==""||s.NodeID==""||!s.Authorized{return fmt.Errorf("invalid or unauthorized local service")};switch s.Class{case ServiceRouter,ServiceMobile,ServiceCCTV,ServiceTV,ServiceLocalDNS,ServiceLocalCache:default:return fmt.Errorf("unsupported local service class")};r.mu.Lock();r.services[s.ID]=s;r.mu.Unlock();return nil}
func(r *LocalFabricRegistry)Get(id string)(LocalServiceProfile,bool){r.mu.RLock();defer r.mu.RUnlock();s,ok:=r.services[id];return s,ok}
