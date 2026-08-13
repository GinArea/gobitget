package bitgetv3

import (
	"github.com/msw-x/moon/ujson"
)

// GetAgentCustomerList - request for POST /api/v2/broker/customer-list (classic v2 API)
// https://www.bitget.com/api-doc/classic/affiliate/customerInfo/GetCustomerList
// Rate limit: 10/sec/UID; available for agent (affiliate) accounts only:
// a non-agent account gets error 70102 "Parameter verification failed-Broker identity authentication failed"
// Supports querying clients registered at any time (no 90-day restriction), data is updated in real time
type GetAgentCustomerList struct {
	// StartTime - Start time (ms); StartTime and EndTime should be both set or both left blank, within 30 days; if not set, all registered clients are returned with no time restriction
	StartTime ujson.StringInt64 `json:",omitempty"`
	// EndTime - End time (ms)
	EndTime ujson.StringInt64 `json:",omitempty"`
	// PageNo - Page number
	PageNo ujson.StringInt64 `json:",omitempty"`
	// PageSize - Page size, 100 default, max 1000
	PageSize ujson.StringInt64 `json:",omitempty"`
	// Uid - UID; set it to confirm whether a specific UID is your referral client
	Uid string `json:",omitempty"`
	// ReferralCode - Referral code
	ReferralCode string `json:",omitempty"`
}

func (o GetAgentCustomerList) Do(c *Client) Response[[]AgentCustomer] {
	return Post(c.v2(), "broker/customer-list", o, forward[[]AgentCustomer])
}

func (o *Client) GetAgentCustomerList(v GetAgentCustomerList) Response[[]AgentCustomer] {
	return v.Do(o)
}

// AgentCustomer - item of POST /api/v2/broker/customer-list response (classic v2 API)
// https://www.bitget.com/api-doc/classic/affiliate/customerInfo/GetCustomerList
type AgentCustomer struct {
	// Uid - UID (documented as String, the response example shows a JSON number)
	Uid ujson.Int64
	// RegisterTime - Register time, unix millisecond timestamp
	RegisterTime ujson.Int64
}
