package backend

import "testing"

func TestAndroidClientContract(t *testing.T){
    p:=AndroidClientProfile{DeviceID:"android-1",UserID:"user-1",Authorized:true,Capabilities:[]AndroidClientCapability{AndroidFTNVPN,AndroidServiceHealth}}
    if err:=ValidateAndroidClient(p);err!=nil{t.Fatal(err)}
}
