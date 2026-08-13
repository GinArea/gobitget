package bitgetv3

import (
	"github.com/msw-x/moon/ujson"
	"github.com/msw-x/moon/ulog"
)

// WsResponse - WebSocket event message: operation ack or error
// Ack: {"event":"subscribe","arg":{...},"connId":"..."}
// Error: {"event":"error","code":30001,"msg":"..."}
// https://www.bitget.com/api-doc/uta/websocket/public/Tickers-Channel
type WsResponse struct {
	// Event - Event type: subscribe / unsubscribe / error
	Event string
	// Arg - Subscribed channel the event relates to
	Arg SubscriptionArgs
	// Code - Error code; documented as a string, but the live server sends a JSON number
	// (e.g. 30001), so ujson.Int64 accepts both; 0 on success
	Code ujson.Int64
	// Msg - Error message
	Msg string
	// ConnId - Connection ID
	ConnId string
}

// IsEvent - returns true for event messages (ack/error); data pushes carry no "event" field
func (o WsResponse) IsEvent() bool {
	return o.Event != ""
}

// IsError - returns true for error events
func (o WsResponse) IsError() bool {
	return o.Event == "error"
}

// IsSubscribe - returns true for subscribe acks
func (o WsResponse) IsSubscribe() bool {
	return o.Event == "subscribe"
}

// IsUnsubscribe - returns true for unsubscribe acks
func (o WsResponse) IsUnsubscribe() bool {
	return o.Event == "unsubscribe"
}

// TopicNotExists - error 30001: the subscribed channel/symbol doesn't exist
// (verified live: subscribing a non-existent symbol returns code 30001)
func (o WsResponse) TopicNotExists() bool {
	return o.Code.Value() == 30001
}

func (o WsResponse) Log(log *ulog.Log) {
	switch o.Event {
	case "subscribe", "unsubscribe":
		log.Debug(o.Event+":", o.Arg.InstType, o.Arg.Topic, o.Arg.Symbol)
	case "error":
		log.Error("error:", o.Code.Value(), o.Msg)
	default:
		log.Warning("unknown event:", o.Event)
	}
}
