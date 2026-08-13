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

// WsCandle - item of Candlestick Channel push data; every field arrives as a JSON string
// The symbol and interval are not repeated here: they come only in the topic arg
// Pushed once per second while trades occur, else once per interval
// The initial snapshot carries recent candle history (~500 items) sorted oldest-first,
// so the current candle is the last item (verified live; the docs show a single item only)
// https://www.bitget.com/api-doc/uta/websocket/public/Candlesticks-Channel
type WsCandle struct {
	// Start - Candle start timestamp, unix ms
	Start ujson.Int64
	// Open - Open price
	Open ujson.Float64
	// High - Highest price
	High ujson.Float64
	// Low - Lowest price
	Low ujson.Float64
	// Close - Close price
	Close ujson.Float64
	// Volume - Trade volume, base coin (absent for rtoken 1m -> 0)
	Volume ujson.Float64
	// Turnover - Turnover, quote coin (absent for rtoken 1m -> 0)
	Turnover ujson.Float64
}
