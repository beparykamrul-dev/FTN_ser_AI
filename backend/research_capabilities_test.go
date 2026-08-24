package backend

import "testing"

func TestResearchRegistry(t *testing.T){
 r:=NewResearchRegistry()
 if err:=r.Upsert(ResearchCapability{ID:"pqc",Name:"Post-Quantum Cryptography",Domain:"cryptography",Maturity:ResearchEstablished,ProductionEligible:true,Enabled:true});err!=nil{t.Fatal(err)}
 if _,ok:=r.Get("pqc");!ok{t.Fatal("capability not found")}
}
func TestResearchRegistryBlocksSpeculativeProduction(t *testing.T){
 r:=NewResearchRegistry()
 if err:=r.Upsert(ResearchCapability{ID:"ftl",Name:"FTL theoretical tunnel",Domain:"theoretical-physics",Maturity:ResearchSpeculative,ProductionEligible:true});err==nil{t.Fatal("expected production eligibility rejection")}
}
