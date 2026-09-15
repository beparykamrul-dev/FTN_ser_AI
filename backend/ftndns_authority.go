package backend

import (
	"sort"
	"strings"
)

type AuthorityMode string

const (
	AuthorityCloudflare AuthorityMode = "cloudflare"
	AuthorityFTN         AuthorityMode = "ftn"
)

type ZoneAuthority struct {
	Zone       string        `json:"zone"`
	Authority  AuthorityMode `json:"authority"`
	ProviderID string        `json:"provider_id,omitempty"`
	DNSSEC     bool          `json:"dnssec"`
	Anycast    bool          `json:"anycast"`
}

var FTNZoneAuthorities = []ZoneAuthority{
	{Zone: "familytimenet.com", Authority: AuthorityCloudflare, ProviderID: "cloudflare", DNSSEC: true, Anycast: false},
	{Zone: "ftndns.com", Authority: AuthorityFTN, DNSSEC: true, Anycast: true},
	{Zone: "ftnddns.net", Authority: AuthorityFTN, DNSSEC: true, Anycast: true},
}

func NormalizeZone(zone string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(zone)), ".")
}

func ZoneAuthorityFor(zone string) (ZoneAuthority, bool) {
	zone = NormalizeZone(zone)
	for _, z := range FTNZoneAuthorities {
		if z.Zone == zone {
			return z, true
		}
	}
	return ZoneAuthority{}, false
}

func AuthorityZones() []ZoneAuthority {
	out := append([]ZoneAuthority(nil), FTNZoneAuthorities...)
	sort.Slice(out, func(i, j int) bool { return out[i].Zone < out[j].Zone })
	return out
}
