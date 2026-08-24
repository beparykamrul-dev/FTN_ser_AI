package backend

import "testing"

func TestCoreRouterRegistrySupportsMikroTik(t *testing.T) {
    r:=NewCoreRouterRegistry()
    p:=RouterProfile{ID:"r1",Name:"FTN-MikroTik-01",Vendor:RouterMikroTik,Address:"10.0.0.1",APIVersion:"routeros",Authorized:true,Enabled:true}
    if err:=r.Upsert(p);err!=nil{t.Fatal(err)}
    got,ok:=r.Get("r1");if !ok||got.Vendor!=RouterMikroTik{t.Fatalf("unexpected router: %#v",got)}
}

func TestCoreRouterRejectsUnauthorized(t *testing.T){
    r:=NewCoreRouterRegistry()
    if err:=r.Upsert(RouterProfile{ID:"r1",Name:"bad",Vendor:RouterMikroTik,Address:"10.0.0.1"});err==nil{t.Fatal("expected authorization error")}
}
