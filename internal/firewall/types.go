package firewall

// Windows Firewall Constants
const (
	NET_FW_PROFILE2_DOMAIN  = 0x0001
	NET_FW_PROFILE2_PRIVATE = 0x0002
	NET_FW_PROFILE2_PUBLIC  = 0x0004
	NET_FW_PROFILE2_ALL     = 0x7

	NET_FW_ACTION_BLOCK = 0
	NET_FW_ACTION_ALLOW = 1

	NET_FW_RULE_DIR_IN  = 1
	NET_FW_RULE_DIR_OUT = 2

	NET_FW_IP_PROTOCOL_TCP = 6
	NET_FW_IP_PROTOCOL_UDP = 17
	NET_FW_IP_PROTOCOL_ANY = 256

	PROGID_POLICY2 = "HNetCfg.FwPolicy2"
	PROGID_RULE    = "HNetCfg.FwRule"

	RULE_PREFIX_OUT = "Enodia-OUT-"
	RULE_PREFIX_IN  = "Enodia-IN-"
	RULE_PREFIX     = "Enodia-"
)

// Rule represents a firewall rule managed by Enodia
type Rule struct {
	Name      string `json:"name"`
	AppPath   string `json:"appPath"`
	Direction string `json:"direction"`
	Enabled   bool   `json:"enabled"`
}

// BlockedApp represents an application with its block status
type BlockedApp struct {
	AppPath         string `json:"appPath"`
	DisplayName     string `json:"displayName"`
	InboundBlocked  bool   `json:"inboundBlocked"`
	OutboundBlocked bool   `json:"outboundBlocked"`
}
