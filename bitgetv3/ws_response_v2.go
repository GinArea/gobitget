package bitgetv3

import (
	"github.com/msw-x/moon/ujson"
	"github.com/msw-x/moon/ulog"
)

// WsResponseV2 - legacy v2 WebSocket event message: operation ack or error
// Ack: {"event":"subscribe","arg":{...}}
// Error: {"event":"error","code":30016,"msg":"Param error","op":"subscribe"}
// https://www.bitget.com/api-doc/spot/websocket/public/Candlesticks-Channel
type WsResponseV2 struct {
	// Event - Event type: subscribe / unsubscribe / error
	Event string
	// Arg - Subscribed channel the event relates to
	Arg SubscriptionArgsV2
	// Code - Error code; the live server sends a JSON number (e.g. 30016), 0 on success
	Code ujson.Int64
	// Msg - Error message
	Msg string
	// Op - Echoed operation, set on errors
	Op string
}

// IsEvent - returns true for event messages (ack/error); data pushes carry no "event" field
func (o WsResponseV2) IsEvent() bool {
	return o.Event != ""
}

// IsError - returns true for error events
func (o WsResponseV2) IsError() bool {
	return o.Event == "error"
}

// IsSubscribe - returns true for subscribe acks
func (o WsResponseV2) IsSubscribe() bool {
	return o.Event == "subscribe"
}

// IsUnsubscribe - returns true for unsubscribe acks
func (o WsResponseV2) IsUnsubscribe() bool {
	return o.Event == "unsubscribe"
}

// ParamError - error 30016: invalid subscription parameter, e.g. an unsupported
// candle timeframe (verified live: candle3H returns code 30016, not 30001 like the v3 WS)
func (o WsResponseV2) ParamError() bool {
	return o.Code.Value() == 30016
}

// ServiceUpgrade - notice 30033, the upgrade warning of the v3 socket: the legacy
// endpoint lives on the same servers and announces the reset the same way
func (o WsResponseV2) ServiceUpgrade() bool {
	return o.Code.Value() == 30033
}

func (o WsResponseV2) Log(log *ulog.Log) {
	switch o.Event {
	case "subscribe", "unsubscribe":
		log.Debug(o.Event+":", o.Arg.InstType, o.Arg.Channel, o.Arg.InstId)
	case "error":
		// not a failure - see WsResponse.Log
		if o.ServiceUpgrade() {
			log.Info("service upgrade:", o.Code.Value(), o.Msg)
		} else {
			log.Error("error:", o.Code.Value(), o.Msg)
		}
	default:
		log.Warning("unknown event:", o.Event)
	}
}
