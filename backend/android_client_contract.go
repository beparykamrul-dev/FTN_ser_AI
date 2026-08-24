package backend

import "fmt"

type AndroidClientCapability string
const (
    AndroidFTNVPN AndroidClientCapability = "ftnvpn"
    AndroidRouterControl AndroidClientCapability = "router-control"
    AndroidServiceHealth AndroidClientCapability = "service-health"
    AndroidNetworkDiagnostics AndroidClientCapability = "network-diagnostics"
)

type AndroidClientProfile struct {
    DeviceID string
    UserID string
    Authorized bool
    Capabilities []AndroidClientCapability
}

func ValidateAndroidClient(p AndroidClientProfile) error {
    if p.DeviceID=="" || p.UserID=="" { return fmt.Errorf("android client identity is required") }
    if !p.Authorized { return fmt.Errorf("android client is not authorized") }
    if len(p.Capabilities)==0 { return fmt.Errorf("android client capability is required") }
    return nil
}
