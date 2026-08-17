package bitgetv3

import (
	"encoding/json"

	"github.com/msw-x/moon/ujson"
)

// Topic - push data envelope: {"data":[...],"arg":{...},"action":"snapshot","ts":123}
// https://www.bitget.com/api-doc/uta/websocket/public/Tickers-Channel
type Topic[T any] struct {
	// Arg - subscribed channel; the symbol comes only here, not in Data
	Arg SubscriptionArgs
	// Action - data push action: snapshot (full push) / update (incremental push)
	Action string
	// Ts - data push timestamp, unix ms
	Ts ujson.Int64
	// Data - subscribed data
	Data T
}

// RawTopic - push envelope with an unparsed payload, routed by Arg
type RawTopic Topic[json.RawMessage]

// UnmarshalRawTopic - parse the raw payload into a typed topic
func UnmarshalRawTopic[T any](raw RawTopic) (ret Topic[T], err error) {
	ret.Arg = raw.Arg
	ret.Action = raw.Action
	ret.Ts = raw.Ts
	err = json.Unmarshal(raw.Data, &ret.Data)
	return
}

// WsPosition - item of Positions Channel push data; every field arrives as a JSON string,
// numeric fields may be empty strings (parsed as 0)
// The field set is verified live against a real position push: the push carries no
// undocumented extras, and of the documented fields only marginRate never arrives
// https://www.bitget.com/api-doc/uta/websocket/private/Positions-Channel
type WsPosition struct {
	// Symbol - Symbol name
	Symbol string
	// MarginCoin - Margin coin
	MarginCoin string
	// MarginSize - Margin size
	MarginSize ujson.Float64
	// MarginMode - Margin mode: crossed, isolated
	MarginMode string
	// PosSide - Position side: long, short
	PosSide string
	// HoldMode - Holding mode: one_way_mode, hedge_mode
	HoldMode string
	// PositionStatus - Position status: opening (ongoing), ended (completed)
	PositionStatus string
	// Size - Position size: size = available + frozen
	Size ujson.Float64
	// Available - Available position size
	Available ujson.Float64
	// Frozen - Frozen position size
	Frozen ujson.Float64
	// AvgPrice - Average open price
	AvgPrice ujson.Float64
	// Leverage - Leverage multiple
	Leverage ujson.Float64
	// CurRealisedPnl - Realised PnL
	CurRealisedPnl ujson.Float64
	// UnrealisedPnl - Unrealised PnL
	UnrealisedPnl ujson.Float64
	// LiqPrice - Estimated liquidation price
	LiqPrice ujson.Float64
	// Mmr - Maintain margin rate
	Mmr ujson.Float64
	// MarginRate - Margin rate; documented, but ABSENT from the live push (verified live) -
	// always parses as 0
	MarginRate ujson.Float64
	// MarkPrice - Mark price
	MarkPrice ujson.Float64
	// OpenFeeTotal - Total opening fee
	OpenFeeTotal ujson.Float64
	// CloseFeeTotal - Total closing fee
	CloseFeeTotal ujson.Float64
	// BreakEvenPrice - Break-even price
	BreakEvenPrice ujson.Float64
	// ProfitRate - Profit rate = unrealized PnL / initial margin, where
	// initial margin = average open price * position size / leverage / margin coin index price
	ProfitRate ujson.Float64
	// TotalFundingFee - Total funding fee over the position's lifetime; 0 indicates
	// that no funding fee has been charged yet
	TotalFundingFee ujson.Float64
	// CashDividend - Cash dividend, unit: USDT
	CashDividend ujson.Float64
	// CreatedTime - Position creation time, unix ms
	CreatedTime ujson.Int64
	// UpdatedTime - Latest position update time, unix ms
	UpdatedTime ujson.Int64
}

// WsOrder - item of Order Channel push data; every field arrives as a JSON string,
// numeric fields may be empty strings (parsed as 0)
// Field names follow the WS docs and differ from the REST Order type
// (HoldSide/TotalProfit/Amount here vs PosSide/no equivalent there; no ExecType or TP/SL fields)
// https://www.bitget.com/api-doc/uta/websocket/private/Order-Channel
type WsOrder struct {
	// Category - Business line: spot, margin, usdt-futures, coin-futures, usdc-futures;
	// the push carries lowercase values, unlike the uppercase REST Category constants
	// (verified live)
	Category Category
	// Symbol - Symbol name
	Symbol string
	// OrderId - Order ID
	OrderId string
	// ClientOid - Client order ID
	ClientOid string
	// RequestId - Custom request ID passed in when modifying the order; only returned
	// if it was provided in the modify order request
	RequestId string
	// Price - Order price (empty for market orders)
	Price ujson.Float64
	// Qty - Order quantity, the unit is base coin
	Qty ujson.Float64
	// Amount - Order amount, the unit is quote coin
	Amount ujson.Float64
	// HoldMode - Holding mode: one_way_mode, hedge_mode
	HoldMode string
	// HoldSide - Position side: long, short
	HoldSide string
	// DelegateType - Delegate type, e.g. normal (normal limit), market, liquidation;
	// full enumeration in the docs
	DelegateType string
	// TradeSide - Trade side: open/close per the docs, but the live push sends detailed
	// variants like the REST Order (e.g. buy_single/sell_single in one-way mode)
	TradeSide string
	// OrderType - Order type: limit, market
	OrderType string
	// TimeInForce - Time in force: ioc, fok, gtc, post_only, rpi
	TimeInForce string
	// Side - Order side: buy, sell
	Side string
	// MarginMode - Margin mode: crossed, isolated
	MarginMode string
	// MarginCoin - Margin coin
	MarginCoin string
	// ReduceOnly - Reduce-only identifier: yes, no
	ReduceOnly string
	// CumExecQty - Cumulative executed quantity
	CumExecQty ujson.Float64
	// CumExecValue - Cumulative executed value
	CumExecValue ujson.Float64
	// AvgPrice - Average execution price; 0 if not executed
	AvgPrice ujson.Float64
	// TotalProfit - Total profit
	TotalProfit ujson.Float64
	// OrderStatus - Order status: new (order matching), partially_filled, filled, cancelled
	OrderStatus string
	// CancelReason - Reason for order cancellation
	CancelReason string
	// Leverage - Leverage multiple
	Leverage ujson.Float64
	// FeeDetail - Fee detail list
	FeeDetail []FeeDetail
	// CreatedTime - Created time, unix ms
	CreatedTime ujson.Int64
	// UpdatedTime - Updated time, unix ms
	UpdatedTime ujson.Int64
	// StpMode - STP mode (Self Trade Prevention): none, cancel_taker, cancel_maker, cancel_both
	StpMode string

	// The fields below are absent from the WS docs but present in the live push (verified live)

	// TpTriggerBy - Take-profit trigger type: market, mark (undocumented in the WS docs)
	TpTriggerBy string
	// SlTriggerBy - Stop-loss trigger type: market, mark (undocumented in the WS docs)
	SlTriggerBy string
	// TakeProfit - Take-profit trigger price; the push key is all-lowercase "takeprofit",
	// matched case-insensitively (undocumented in the WS docs)
	TakeProfit ujson.Float64
	// StopLoss - Stop-loss trigger price; the push key is all-lowercase "stoploss",
	// matched case-insensitively (undocumented in the WS docs)
	StopLoss ujson.Float64
	// TpOrderType - Take-profit order type: limit, market (undocumented in the WS docs)
	TpOrderType string
	// SlOrderType - Stop-loss order type: limit, market (undocumented in the WS docs)
	SlOrderType string
	// TpLimitPrice - Take-profit limit order execution price (undocumented in the WS docs)
	TpLimitPrice ujson.Float64
	// SlLimitPrice - Stop-loss limit order execution price (undocumented in the WS docs)
	SlLimitPrice ujson.Float64
	// MatchType - Match type, e.g. 0 (undocumented in the WS docs, meaning unknown)
	MatchType ujson.Int64
}

// WsTicker - item of Tickers Channel push data; every field arrives as a JSON string
// The symbol is not repeated here: it comes only in the topic arg
// https://www.bitget.com/api-doc/uta/websocket/public/Tickers-Channel
type WsTicker struct {
	// LastPrice - Latest traded price
	LastPrice ujson.Float64
	// Bid1Price - Best bid price
	Bid1Price ujson.Float64
	// Bid1Size - Best bid size
	Bid1Size ujson.Float64
	// Ask1Price - Best ask price
	Ask1Price ujson.Float64
	// Ask1Size - Best ask size
	Ask1Size ujson.Float64
	// HighPrice24h - Highest price in the last 24 hours
	HighPrice24h ujson.Float64
	// LowPrice24h - Lowest price in the last 24 hours
	LowPrice24h ujson.Float64
	// OpenPrice24h - Open price 24 hours ago
	OpenPrice24h ujson.Float64
	// Volume24h - 24-hour trading volume, base coin
	Volume24h ujson.Float64
	// Turnover24h - 24-hour trading volume, quote coin
	Turnover24h ujson.Float64
	// Price24hPcnt - 24-hour price change as a fraction (e.g. 0.01833 = 1.833%)
	Price24hPcnt ujson.Float64
	// PlatformTurnover24h - 24-hour platform trading volume (rtoken only)
	PlatformTurnover24h ujson.Float64
	// IndexPrice - Index price (futures only)
	IndexPrice ujson.Float64
	// MarkPrice - Mark price (futures only)
	MarkPrice ujson.Float64
	// FundingRate - Current funding rate (futures only)
	FundingRate ujson.Float64
	// NextFundingTime - Next funding settlement time, unix ms (futures only)
	NextFundingTime ujson.Int64
	// OpenInterest - Open interest (futures only)
	OpenInterest ujson.Float64
	// DeliveryTime - Delivery time, unix ms (delivery futures only)
	DeliveryTime ujson.Int64
	// DeliveryStartTime - Delivery start time, unix ms (delivery futures only)
	DeliveryStartTime ujson.Int64
	// DeliveryStatus - Delivery status: configuration / normal / before / period (delivery futures only)
	DeliveryStatus string
}
