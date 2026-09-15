package backend

// ProviderRole describes how FTN should use an external platform.
type ProviderRole string

const (
	ProviderPrimaryControl ProviderRole = "primary_control"
	ProviderGlobalDNS      ProviderRole = "global_dns"
	ProviderEdgeSecurity   ProviderRole = "edge_security"
	ProviderCompute        ProviderRole = "compute"
	ProviderWebRender      ProviderRole = "web_render"
	ProviderGeoAssist      ProviderRole = "geo_assist"
)

type ProviderPolicy struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Priority int            `json:"priority"`
	Roles    []ProviderRole `json:"roles"`
	Enabled  bool           `json:"enabled"`
	ReadOnly bool           `json:"read_only"`
	Reason   string         `json:"reason"`
}

// FTN keeps Google and Cloudflare as strategic first-choice helpers while
// retaining provider portability. They are not hard-coded as the only path.
func StrategicProviderPolicy() []ProviderPolicy {
	return []ProviderPolicy{
		{ID: "cloudflare", Name: "Cloudflare", Priority: 1, Roles: []ProviderRole{ProviderPrimaryControl, ProviderGlobalDNS, ProviderEdgeSecurity, ProviderWebRender}, Enabled: true, ReadOnly: false, Reason: "primary global DNS, edge, security and platform integration"},
		{ID: "google", Name: "Google Cloud", Priority: 1, Roles: []ProviderRole{ProviderGlobalDNS, ProviderCompute, ProviderGeoAssist}, Enabled: true, ReadOnly: false, Reason: "global anycast DNS, cloud infrastructure and geo-aware services"},
		{ID: "render", Name: "Render", Priority: 2, Roles: []ProviderRole{ProviderWebRender}, Enabled: true, ReadOnly: false, Reason: "managed application and static-site deployment"},
		{ID: "ftn", Name: "FTN", Priority: 0, Roles: []ProviderRole{ProviderPrimaryControl, ProviderGlobalDNS, ProviderGeoAssist}, Enabled: true, ReadOnly: false, Reason: "FTN remains the control-plane authority and provider abstraction layer"},
	}
}
