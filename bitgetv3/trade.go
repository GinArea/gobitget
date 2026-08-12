package bitgetv3

import (
	"github.com/msw-x/moon/ujson"
)

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

// Order - response for GET /api/v3/trade/order-info (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Details
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

// FeeDetail - item of feeDetail list in GET /api/v3/trade/order-info response (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Order-Details
type FeeDetail struct {
	// FeeCoin - Fee coin (empty when no fee was charged)
	FeeCoin string
	// Fee - Total fee (empty when no fee was charged)
	Fee ujson.Float64
}
