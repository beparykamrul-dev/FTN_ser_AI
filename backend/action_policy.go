package backend

import (
    "crypto/subtle"
    "errors"
    "os"
    "strings"
)

func actionNeedsApproval(kind string) bool {
    switch strings.ToLower(strings.TrimSpace(kind)) {
    case "provider.dns.record.upsert", "provider.dns.record.delete", "provider.zone.create", "provider.zone.delete":
        return true
    default:
        return false
    }
}

func approvalConfigured() bool { return strings.TrimSpace(os.Getenv("FTN_ACTION_APPROVAL_SECRET")) != "" }

func validateApprovalSecret(got string) error {
    expected := strings.TrimSpace(os.Getenv("FTN_ACTION_APPROVAL_SECRET"))
    if expected == "" { return errors.New("action approval is not configured") }
    a, b := []byte(got), []byte(expected)
    if len(a) != len(b) || subtle.ConstantTimeCompare(a, b) != 1 { return errors.New("invalid approval credential") }
    return nil
}
