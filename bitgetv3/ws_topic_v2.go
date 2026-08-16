package bitgetv3

import (
	"encoding/json"
	"fmt"

	"github.com/msw-x/moon/ujson"
)

// TopicV2 - legacy v2 push data envelope: {"action":"snapshot","arg":{...},"data":[...],"ts":123}
// https://www.bitget.com/api-doc/spot/websocket/public/Candlesticks-Channel
type TopicV2[T any] struct {
	// Arg - subscribed channel; the symbol comes only here, not in Data
	Arg SubscriptionArgsV2
	// Action - data push action: snapshot (full push) / update (incremental push)
	Action string
	// Ts - data push timestamp, unix ms
	Ts ujson.Int64
	// Data - subscribed data
	Data T
}

// RawTopicV2 - v2 push envelope with an unparsed payload, routed by Arg
type RawTopicV2 TopicV2[json.RawMessage]

// UnmarshalRawTopicV2 - parse the raw payload into a typed topic
func UnmarshalRawTopicV2[T any](raw RawTopicV2) (ret TopicV2[T], err error) {
	ret.Arg = raw.Arg
	ret.Action = raw.Action
	ret.Ts = raw.Ts
	err = json.Unmarshal(raw.Data, &ret.Data)
	return
}

// WsCandleV2 - item of the v2 Candlestick Channel push data; the wire format is an array
// of 8 strings: [ts, open, high, low, close, baseVolume, quoteVolume, usdtVolume]
// (field order verified live on a BTC-quoted pair; for USDT-quoted symbols the last two match).
// The initial snapshot carries recent candle history (500 items or the full available depth)
// sorted oldest-first, so the current candle is the last item - same behavior as the v3 WS
// https://www.bitget.com/api-doc/spot/websocket/public/Candlesticks-Channel
type WsCandleV2 struct {
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
	// Volume - Trade volume, base coin
	Volume ujson.Float64
	// QuoteVolume - Trade volume, quote coin
	QuoteVolume ujson.Float64
	// UsdtVolume - Trade volume in USDT
	UsdtVolume ujson.Float64
}

// UnmarshalJSON - decode the wire array form into the named fields
func (o *WsCandleV2) UnmarshalJSON(data []byte) error {
	var row []json.RawMessage
	if err := json.Unmarshal(data, &row); err != nil {
		return err
	}
	if len(row) != 8 {
		return fmt.Errorf("candle row length is %d, expected 8", len(row))
	}
	fields := []struct {
		name string
		u    json.Unmarshaler
	}{
		{"start", &o.Start},
		{"open", &o.Open},
		{"high", &o.High},
		{"low", &o.Low},
		{"close", &o.Close},
		{"volume", &o.Volume},
		{"quoteVolume", &o.QuoteVolume},
		{"usdtVolume", &o.UsdtVolume},
	}
	for i, f := range fields {
		if err := f.u.UnmarshalJSON(row[i]); err != nil {
			return fmt.Errorf("candle %s: %w", f.name, err)
		}
	}
	return nil
}
