package backend

import (
    "errors"
    "strings"
)

// ValidateProviderZoneOperation keeps FTN-authoritative zones under FTN control.
// External providers may only operate on zones explicitly delegated to them.
func ValidateProviderZoneOperation(providerID, zone string, write bool) error {
    providerID = strings.ToLower(strings.TrimSpace(providerID))
    zone = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(zone), "."))
    if providerID == "" || zone == "" {
        return errors.New("provider and zone are required")
    }
    for _, authority := range AuthorityZones() {
        if authority.Zone != zone {
            continue
        }
        if authority.Mode == AuthorityFTNAuthoritative && providerID != "ftn" {
            return errors.New("zone is FTN-authoritative; external provider operation denied")
        }
        if write && authority.ReadOnlyMirror {
            return errors.New("zone is configured as a read-only provider mirror")
        }
        if authority.ProviderID != providerID && write {
            return errors.New("provider is not the configured write authority for this zone")
        }
    }
    return nil
}
