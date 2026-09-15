package backend

type AuthorityMode string

const (
	AuthorityCloudflarePrimary AuthorityMode = "cloudflare_primary"
	AuthorityFTNAuthoritative  AuthorityMode = "ftn_authoritative"
)

type ZoneAuthority struct {
	Zone            string        `json:"zone"`
	Mode            AuthorityMode `json:"mode"`
	ProviderID      string        `json:"provider_id"`
	DNSSEC          bool          `json:"dnssec"`
	AnycastReady    bool          `json:"anycast_ready"`
	AnycastEnabled  bool          `json:"anycast_enabled"`
	BGPAdvertised   bool          `json:"bgp_advertised"`
	ReadOnlyMirror  bool          `json:"read_only_mirror"`
	Managed         bool          `json:"managed"`
}

type AnycastPrefix struct {
	Prefix        string `json:"prefix"`
	ASN           uint32 `json:"asn"`
	Enabled       bool   `json:"enabled"`
	BGPAdvertised bool   `json:"bgp_advertised"`
	Healthy       bool   `json:"healthy"`
	RouteCount    int    `json:"route_count"`
}

func AuthorityZones() []ZoneAuthority {
	return []ZoneAuthority{
		{Zone: "familytimenet.com", Mode: AuthorityCloudflarePrimary, ProviderID: "cloudflare", DNSSEC: true, AnycastReady: false, AnycastEnabled: false, BGPAdvertised: false, ReadOnlyMirror: true, Managed: true},
		{Zone: "ftndns.com", Mode: AuthorityFTNAuthoritative, ProviderID: "ftn", DNSSEC: true, AnycastReady: true, AnycastEnabled: false, BGPAdvertised: false, ReadOnlyMirror: false, Managed: true},
		{Zone: "ftnddns.net", Mode: AuthorityFTNAuthoritative, ProviderID: "ftn", DNSSEC: true, AnycastReady: true, AnycastEnabled: false, BGPAdvertised: false, ReadOnlyMirror: false, Managed: true},
	}
}

func DefaultAnycastPrefixes() []AnycastPrefix {
	return []AnycastPrefix{
		{Prefix: "", ASN: 0, Enabled: false, BGPAdvertised: false, Healthy: false, RouteCount: 0},
	}
}
