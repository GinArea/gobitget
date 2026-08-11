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

// GetDepositAddress - request for GET /api/v3/account/deposit-address (UTA)
// https://www.bitget.com/api-doc/uta/account/deposit/Get-Deposit-Address
// Rate limit: 10/sec/UID, permission: UTA mgt. (read)
type GetDepositAddress struct {
	// Coin - Coin name, e.g. USDT
	Coin string
	// Chain - Chain name, e.g. trc20 (auto-matched if omitted)
	Chain string `url:",omitempty"`
	// Size - Deposit quantity, for BTC Lightning Network only, range: 0.000001 - 0.001
	Size string `url:",omitempty"`
}

func (o GetDepositAddress) Do(c *Client) Response[DepositAddress] {
	return Get(c, "account/deposit-address", o, forward[DepositAddress])
}

func (o *Client) GetDepositAddress(v GetDepositAddress) Response[DepositAddress] {
	return v.Do(o)
}

// DepositAddress - response for GET /api/v3/account/deposit-address (UTA)
// https://www.bitget.com/api-doc/uta/account/deposit/Get-Deposit-Address
type DepositAddress struct {
	// Address - Deposit address
	Address string
	// Chain - Chain name
	Chain string
	// Coin - Coin name
	Coin string
	// Tag - Tag (API may return null)
	Tag string
	// Url - Blockchain explorer address
	Url string
}

// GetDepositRecords - request for GET /api/v3/account/deposit-records (UTA)
// https://www.bitget.com/api-doc/uta/account/deposit/Get-Deposit-Records
// Rate limit: 10/sec/UID, permission: UTA mgt. (read)
type GetDepositRecords struct {
	// StartTime - Query record start time, unix millisecond timestamp (interval with EndTime cannot be greater than 30 days)
	StartTime int64
	// EndTime - Query record end time, unix millisecond timestamp
	EndTime int64
	// Coin - Coin name; if left blank, all currency deposit records will be retrieved
	Coin string `url:",omitempty"`
	// OrderId - Order ID, used for specifying order queries
	OrderId string `url:",omitempty"`
	// Limit - Items per page, default 20, max 100
	Limit int `url:",omitempty"`
	// Cursor - Cursor ID for pagination; use the smallest orderId returned from the previous query
	Cursor string `url:",omitempty"`
}

func (o GetDepositRecords) Do(c *Client) Response[[]DepositRecord] {
	return Get(c, "account/deposit-records", o, forward[[]DepositRecord])
}

func (o *Client) GetDepositRecords(v GetDepositRecords) Response[[]DepositRecord] {
	return v.Do(o)
}

// DepositRecord - item of GET /api/v3/account/deposit-records response (UTA)
// https://www.bitget.com/api-doc/uta/account/deposit/Get-Deposit-Records
type DepositRecord struct {
	// OrderId - Order ID
	OrderId string
	// RecordId - Deposit record ID: the on-chain hash when Dest is on_chain, the order ID when Dest is internal_transfer
	RecordId string
	// ClientOid - Client order ID (not documented, API returns it, may be null)
	ClientOid string
	// Coin - Coin name
	Coin string
	// Type - Operation type: deposit
	Type string
	// Dest - Deposit type: on_chain, internal_transfer
	Dest string
	// Size - Deposit quantity
	Size ujson.Float64
	// Status - Deposit status: pending, success, fail
	Status string
	// FromAddress - Deposit initiator: the on-chain address when Dest is on_chain, the UID/email/mobile when Dest is internal_transfer
	FromAddress string
	// ToAddress - Deposit recipient: the on-chain address when Dest is on_chain, the UID/email/mobile when Dest is internal_transfer
	ToAddress string
	// Chain - Deposit network (can be ignored when Dest is internal_transfer)
	Chain string
	// CreatedTime - Deposit record creation time, unix millisecond timestamp
	CreatedTime ujson.Int64
	// UpdatedTime - Deposit record update time, unix millisecond timestamp
	UpdatedTime ujson.Int64
}

// GetFundingAssets - request for GET /api/v3/account/funding-assets (UTA)
// https://www.bitget.com/api-doc/uta/account/Get-Account-Funding-Assets
// Rate limit: 20/sec/UID, permission: UTA mgt. (read)
type GetFundingAssets struct {
	// Coin - Coin name; if omitted, all coins with assets are returned
	// (Pre-IPO tokens use mixed case, e.g. preSPAX)
	Coin string `url:",omitempty"`
}

func (o GetFundingAssets) Do(c *Client) Response[[]FundingAsset] {
	return Get(c, "account/funding-assets", o, forward[[]FundingAsset])
}

func (o *Client) GetFundingAssets(v GetFundingAssets) Response[[]FundingAsset] {
	return v.Do(o)
}

// FundingAsset - item of GET /api/v3/account/funding-assets response (UTA)
// https://www.bitget.com/api-doc/uta/account/Get-Account-Funding-Assets
type FundingAsset struct {
	// Coin - Coin name
	Coin string
	// Balance - Balance. Unit: the current asset coin
	Balance ujson.Float64
	// Available - Available. Unit: the current asset coin
	Available ujson.Float64
	// Frozen - Frozen. Unit: the current asset coin
	Frozen ujson.Float64
}

// GetAccountAssets - request for GET /api/v3/account/assets (UTA)
// https://www.bitget.com/api-doc/uta/account/Get-Account
// Rate limit: 20/sec/UID, permission: UTA mgt. (read), no request parameters
type GetAccountAssets struct{}

func (o GetAccountAssets) Do(c *Client) Response[AccountAssets] {
	return Get(c, "account/assets", o, forward[AccountAssets])
}

func (o *Client) GetAccountAssets() Response[AccountAssets] {
	return GetAccountAssets{}.Do(o)
}

// AccountAssets - response for GET /api/v3/account/assets (UTA)
// https://www.bitget.com/api-doc/uta/account/Get-Account
type AccountAssets struct {
	// AccountEquity - Account equity (USD)
	AccountEquity ujson.Float64
	// UsdtEquity - Account equity (USDT)
	UsdtEquity ujson.Float64
	// BtcEquity - Account equity (BTC)
	BtcEquity ujson.Float64
	// UnrealisedPnl - Unrealised profit and loss (USD)
	UnrealisedPnl ujson.Float64
	// UsdtUnrealisedPnl - Unrealised profit and loss (USDT)
	UsdtUnrealisedPnl ujson.Float64
	// BtcUnrealizedPnl - Unrealised profit and loss (BTC)
	// (API spells this field with "z" - btcUnrealizedPnl, unlike the USD/USDT variants with "s")
	BtcUnrealizedPnl ujson.Float64
	// EffEquity - Effective equity (USD): the net value available for margin in spot and perpetual trades under cross-margin mode, converted to fiat
	EffEquity ujson.Float64
	// Mmr - Maintenance margin (USD): the minimum margin required to maintain the position, converted to fiat
	Mmr ujson.Float64
	// Imr - Initial margin (USD): total initial margin of assets in base coin, converted to fiat
	Imr ujson.Float64
	// MgnRatio - Margin ratio
	MgnRatio ujson.Float64
	// PositionMgnRatio - Position MMR
	PositionMgnRatio ujson.Float64
	// PositionValue - Position value (USD)
	PositionValue ujson.Float64
	// Leverage - Account leverage, non-negative number
	Leverage ujson.Float64
	// Assets - Asset list (only non-zero balances are returned)
	Assets []AccountAsset
}

// AccountAsset - item of assets list in GET /api/v3/account/assets response (UTA)
// https://www.bitget.com/api-doc/uta/account/Get-Account
type AccountAsset struct {
	// Coin - Coin name
	Coin string
	// Equity - Coin equity
	Equity ujson.Float64
	// UsdValue - Coin equity (USD)
	UsdValue ujson.Float64
	// Balance - Coin balance
	Balance ujson.Float64
	// Available - Available
	Available ujson.Float64
	// Debt - Debt (applicable when placing margin orders)
	Debt ujson.Float64
	// Locked - Locked (applicable when placing spot orders)
	Locked ujson.Float64
	// Bonus - USDT bonus amount
	Bonus ujson.Float64
}
