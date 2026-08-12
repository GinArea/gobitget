package bitgetv3

import (
	"github.com/msw-x/moon/ujson"
)

// PlaceOrder - request for POST /api/v3/trade/place-order (UTA)
// https://www.bitget.com/api-doc/uta/trade/Place-Order
// Rate limit: 10/sec/UID, permission: UTA trade (read & write)
// Order limit: 400 orders across all futures pairs, 400 across all spot and margin pairs
// Note: on errors 40010/40725/45001 (timeout) query the order by ClientOid to confirm the final result
type PlaceOrder struct {
	// Category - Product type: SPOT, MARGIN, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
	Category Category
	// Symbol - Symbol name, e.g. BTCUSDT
	Symbol string
	// Qty - Order quantity; spot/margin: quote coin for market buy, base coin otherwise; USDT/USDC futures: base coin; COIN futures: quote coin
	Qty ujson.StringFloat64
	// Price - Order price; required for limit orders, not applicable for market orders
	Price ujson.StringFloat64 `json:",omitempty"`
	// Side - Order side: buy, sell
	Side string
	// OrderType - Order type: limit, market
	OrderType string
	// TimeInForce - Time in force: ioc, fok, gtc, post_only, rpi; required for limit orders, defaults to gtc
	TimeInForce string `json:",omitempty"`
	// PosSide - Position side: long, short; required in hedge-mode position; available only for futures
	PosSide string `json:",omitempty"`
	// ClientOid - Client order ID, 1-32 chars matching ^[\.A-Z\:/a-z0-9_-]{1,32}$
	ClientOid string `json:",omitempty"`
	// ReduceOnly - Reduce-only identifier: yes, no (default no)
	ReduceOnly string `json:",omitempty"`
	// StpMode - STP mode (Self Trade Prevention): none (default), cancel_taker, cancel_maker, cancel_both
	StpMode string `json:",omitempty"`
	// TpTriggerBy - Preset take-profit trigger type: market (default), mark; futures only
	TpTriggerBy string `json:",omitempty"`
	// SlTriggerBy - Preset stop-loss trigger type: market (default), mark; futures only
	SlTriggerBy string `json:",omitempty"`
	// TakeProfit - Preset take-profit trigger price
	TakeProfit ujson.StringFloat64 `json:",omitempty"`
	// StopLoss - Preset stop-loss trigger price
	StopLoss ujson.StringFloat64 `json:",omitempty"`
	// TpOrderType - Take-profit strategy order type: limit, market
	TpOrderType string `json:",omitempty"`
	// SlOrderType - Stop-loss strategy order type: limit, market
	SlOrderType string `json:",omitempty"`
	// TpLimitPrice - Take-profit strategy order execution price; valid only when TpOrderType is limit
	TpLimitPrice ujson.StringFloat64 `json:",omitempty"`
	// SlLimitPrice - Stop-loss strategy order execution price; valid only when SlOrderType is limit
	SlLimitPrice ujson.StringFloat64 `json:",omitempty"`
	// MarginMode - Margin mode: crossed (default), isolated; available only for futures
	MarginMode string `json:",omitempty"`
}

func (o PlaceOrder) Do(c *Client) Response[PlacedOrder] {
	return Post(c, "trade/place-order", o, forward[PlacedOrder])
}

func (o *Client) PlaceOrder(v PlaceOrder) Response[PlacedOrder] {
	return v.Do(o)
}

// PlacedOrder - response for POST /api/v3/trade/place-order (UTA)
// https://www.bitget.com/api-doc/uta/trade/Place-Order
type PlacedOrder struct {
	// OrderId - Order ID (null when a one-way-mode reduce-only order was auto-replaced)
	OrderId string
	// ClientOid - Client order ID
	ClientOid string
}

// GetOrder - request for GET /api/v3/trade/order-info (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Details
// Rate limit: 20/sec/UID, permission: UTA trade (read)
type GetOrder struct {
	// OrderId - Order ID; either ClientOid or OrderId must be provided; if both are present or do not match, OrderId takes priority
	OrderId string `url:",omitempty"`
	// ClientOid - Client order ID; either ClientOid or OrderId must be provided
	ClientOid string `url:",omitempty"`
}

func (o GetOrder) Do(c *Client) Response[Order] {
	return Get(c, "trade/order-info", o, forward[Order])
}

func (o *Client) GetOrder(v GetOrder) Response[Order] {
	return v.Do(o)
}

// GetOpenOrders - request for GET /api/v3/trade/unfilled-orders (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Pending
// Rate limit: 20/sec/UID, permission: UTA trade (read)
// Order limit: 400 orders across all futures pairs, 400 across all spot and margin pairs
type GetOpenOrders struct {
	// Category - Product type: SPOT, MARGIN, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES; if omitted, all categories are returned
	Category Category `url:",omitempty"`
	// Symbol - Symbol name, e.g. BTCUSDT
	Symbol string `url:",omitempty"`
	// StartTime - Start timestamp, unix millisecond timestamp
	StartTime int64 `url:",omitempty"`
	// EndTime - End timestamp, unix millisecond timestamp
	EndTime int64 `url:",omitempty"`
	// Limit - Limit per page, default 100, max 100
	Limit int `url:",omitempty"`
	// Cursor - Pagination cursor: omitted in the first query, taken from the previous query for subsequent pages
	Cursor string `url:",omitempty"`
}

func (o GetOpenOrders) Do(c *Client) Response[OpenOrders] {
	return Get(c, "trade/unfilled-orders", o, forward[OpenOrders])
}

func (o *Client) GetOpenOrders(v GetOpenOrders) Response[OpenOrders] {
	return v.Do(o)
}

// OpenOrders - response for GET /api/v3/trade/unfilled-orders (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Pending
type OpenOrders struct {
	// List - Order list (empty array when there are no open orders; tradeSide, cancelReason and execType are not returned by this endpoint)
	List []Order
	// Cursor - Cursor for the next page: the smallest orderId in the current page (null when the list is empty)
	Cursor string
}

// Order - response for GET /api/v3/trade/order-info, item of list in GET /api/v3/trade/unfilled-orders (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Details
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Pending
type Order struct {
	// OrderId - Order ID
	OrderId string
	// ClientOid - Client order ID
	ClientOid string
	// Category - Product type: SPOT, MARGIN, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
	Category Category
	// Symbol - Symbol name
	Symbol string
	// OrderType - Order type: limit, market
	OrderType string
	// Side - Order side: buy, sell
	Side string
	// Price - Order price (empty for market orders)
	Price ujson.Float64
	// Qty - Order quantity, the unit is base coin
	Qty ujson.Float64
	// Amount - Order amount, the unit is quote coin
	Amount ujson.Float64
	// CumExecQty - Cumulative executed quantity, the unit is base coin
	CumExecQty ujson.Float64
	// CumExecValue - Cumulative executed value, the unit is quote coin
	CumExecValue ujson.Float64
	// AvgPrice - Average executed price
	AvgPrice ujson.Float64
	// TimeInForce - Time in force: ioc, fok, gtc, post_only, rpi
	TimeInForce string
	// OrderStatus - Order status: live, new, partially_filled, filled, cancelled
	OrderStatus string
	// PosSide - Position side: long, short
	PosSide string
	// HoldMode - Position mode: one_way_mode, hedge_mode
	HoldMode string
	// TradeSide - Trade side: open/close (live API returns detailed variants, e.g. open_long, close_long)
	TradeSide string
	// ReduceOnly - Reduce-only identifier: YES, NO; available only for futures
	ReduceOnly string
	// MarginMode - Margin mode: crossed, isolated; available only for futures
	MarginMode string
	// StpMode - STP mode (Self Trade Prevention): none, cancel_taker, cancel_maker, cancel_both
	StpMode string
	// TakeProfit - Take-profit trigger price (API returns null when not set)
	TakeProfit ujson.Float64
	// StopLoss - Stop-loss trigger price (API returns null when not set)
	StopLoss ujson.Float64
	// TpTriggerBy - Take-profit trigger type: market (market price), mark (mark price)
	TpTriggerBy string
	// SlTriggerBy - Stop-loss trigger type: market (market price), mark (mark price)
	SlTriggerBy string
	// TpOrderType - Take-profit order type: limit, market
	TpOrderType string
	// SlOrderType - Stop-loss order type: limit, market
	SlOrderType string
	// TpLimitPrice - Take-profit limit order execution price
	TpLimitPrice ujson.Float64
	// SlLimitPrice - Stop-loss limit order execution price
	SlLimitPrice ujson.Float64
	// FeeDetail - Fee detail
	FeeDetail []FeeDetail
	// DelegateType - Delegate type, e.g. normal (normal limit), market (best-price normal), liquidation; full enumeration in the docs
	DelegateType string
	// CancelReason - Cancel reason, e.g. normal_cancel; empty when the order is not cancelled
	CancelReason string
	// ExecType - Execution type: normal, offset (netting of hedged positions), reduce (forced reduction), liquidation, delivery
	ExecType string
	// CreatedTime - Created timestamp, unix millisecond timestamp
	CreatedTime ujson.Int64
	// UpdatedTime - Updated timestamp, unix millisecond timestamp
	UpdatedTime ujson.Int64
}

// FeeDetail - item of feeDetail list in GET /api/v3/trade/order-info and GET /api/v3/trade/unfilled-orders responses (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Details
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Pending
type FeeDetail struct {
	// FeeCoin - Fee coin (empty when no fee was charged)
	FeeCoin string
	// Fee - Total fee (empty when no fee was charged)
	Fee ujson.Float64
}
