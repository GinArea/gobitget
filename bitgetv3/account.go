package bitgetv3

import (
	"github.com/msw-x/moon/ujson"
)

// GetAccountInfo - request for GET /api/v3/account/info (UTA)
// https://www.bitget.com/api-doc/uta/account/Get-Account-Info
// Rate limit: 5/sec/UID, no request parameters
type GetAccountInfo struct{}

func (o GetAccountInfo) Do(c *Client) Response[AccountInfo] {
	return Get(c, "account/info", o, forward[AccountInfo])
}

func (o *Client) GetAccountInfo() Response[AccountInfo] {
	return GetAccountInfo{}.Do(o)
}

// AccountInfo - response for GET /api/v3/account/info (UTA)
// https://www.bitget.com/api-doc/uta/account/Get-Account-Info
type AccountInfo struct {
	// UserId - User ID
	UserId ujson.Int64
	// InviterId - Inviter UID (API returns a string, or null when absent)
	InviterId ujson.Int64
	// ParentId - Parent account UID. Only has a value when the calling account is a sub-account
	// (API returns a JSON number, or null for the main account)
	ParentId ujson.Int64
	// ChannelCode - Channel invitation code
	ChannelCode string
	// Channel - Channel
	Channel string
	// Ips - IP whitelist
	Ips string
	// PermType - Permission type: read-only, read-and-write (actual API returns underscore variants: read_and_write)
	PermType string
	// Permissions - Permissions list: uta_mgt, uta_trade, withdraw, copy_futures_position, copy_futures_order
	Permissions []string
	// RegisTime - Account registration time (Unix timestamp in milliseconds)
	RegisTime ujson.Int64
}
