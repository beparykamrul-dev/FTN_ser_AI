package backend

type AccessDecision struct { Actor string `json:"actor"`; Action string `json:"action"`; Resource string `json:"resource"`; Allowed bool `json:"allowed"`; RequiresApproval bool `json:"requires_approval"` }
func IsDestructive(action string) bool { switch action { case "drop","truncate","delete_all","restore","rotate_root_credentials","withdraw_all_routes": return true }; return false }
