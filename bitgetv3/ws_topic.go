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
