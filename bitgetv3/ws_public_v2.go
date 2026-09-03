package bitgetv3

import (
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

// WsPublicV2 - public client of the legacy v2 WebSocket
// Exists only to stream kline timeframes the v3 WS rejects: 2H, 8H and the utc-suffixed
// variants (6Hutc/12Hutc/1Dutc/3Dutc/1Wutc/1Mutc); a single endpoint serves all product
// types, the category is passed per subscription. See BITGET_API.md
// https://www.bitget.com/api-doc/spot/websocket/public/Candlesticks-Channel
type WsPublicV2 struct {
	c              *WsClientV2
	ready          bool
	onReady        func()
	onConnected    func()
	onDisconnected func()
	onDialError    func(error) bool
	onError        func(WsResponseV2)
	subscriptions  *SubscriptionsV2
}

func NewWsPublicV2() *WsPublicV2 {
	o := new(WsPublicV2)
	o.c = NewWsClientV2()
	o.c.WithPath(WsPublicV2Path)
	o.subscriptions = NewSubscriptionsV2(o)
	return o
}

// Builder methods (return *WsPublicV2 for chaining)

func (o *WsPublicV2) WithLog(log *ulog.Log) *WsPublicV2 {
	o.c.WithLog(log)
	return o
}

func (o *WsPublicV2) WithProxy(proxy string) *WsPublicV2 {
	o.c.WithProxy(proxy)
	return o
}

func (o *WsPublicV2) WithLogRequest(enable bool) *WsPublicV2 {
	o.c.WithLogRequest(enable)
	return o
}

func (o *WsPublicV2) WithLogResponse(enable bool) *WsPublicV2 {
	o.c.WithLogResponse(enable)
	return o
}

func (o *WsPublicV2) WithOnDialError(f func(error) bool) *WsPublicV2 {
	o.onDialError = f
	return o
}

func (o *WsPublicV2) WithOnConnected(f func()) *WsPublicV2 {
	o.onConnected = f
	return o
}

func (o *WsPublicV2) WithOnDisconnected(f func()) *WsPublicV2 {
	o.onDisconnected = f
	return o
}

func (o *WsPublicV2) WithOnReady(f func()) *WsPublicV2 {
	o.onReady = f
	return o
}

func (o *WsPublicV2) WithOnError(f func(WsResponseV2)) *WsPublicV2 {
	o.onError = f
	return o
}

func (o *WsPublicV2) Close() {
	o.c.Close()
}

func (o *WsPublicV2) Transport() *uws.Options {
	return o.c.Transport()
}

func (o *WsPublicV2) Running() bool {
	return o.c.Running()
}

func (o *WsPublicV2) Connected() bool {
	return o.c.Connected()
}

func (o *WsPublicV2) Reconnect() {
	o.c.Reconnect()
}

func (o *WsPublicV2) Ready() bool {
	return o.ready
}

func (o *WsPublicV2) Run() {
	o.c.WithOnConnected(func() {
		if o.onConnected != nil {
			o.onConnected()
		}
		o.markReady()
	})
	o.c.WithOnDisconnected(func() {
		o.ready = false
		// no ack will arrive for a pending unsubscribe, and nothing is in flight
		// any more: the server-side subscription state died with the socket
		o.subscriptions.resetUnsubscribing()
		if o.onDisconnected != nil {
			o.onDisconnected()
		}
	})
	o.c.WithOnDialError(o.handleDialError)
	o.c.WithOnResponse(o.onResponse)
	o.c.WithOnTopic(o.onTopic)
	o.c.Run()
}

// markReady - open the gate for subscriptions and resubscribe to everything registered;
// ready is set first so a concurrent subscribe is sent immediately instead of being
// silently skipped between subscribeAll and the flag flip (a duplicate subscribe is harmless)
func (o *WsPublicV2) markReady() {
	o.ready = true
	o.subscriptions.subscribeAll()
	if o.onReady != nil {
		o.onReady()
	}
}

func (o *WsPublicV2) subscribe(args ...SubscriptionArgsV2) {
	o.c.Subscribe(args...)
}

func (o *WsPublicV2) unsubscribe(args ...SubscriptionArgsV2) {
	o.c.Unsubscribe(args...)
}

func (o *WsPublicV2) onResponse(r WsResponseV2) {
	// a routine notice, not an error - see WsBase.onResponse
	if r.IsError() && r.ServiceUpgrade() {
		r.Log(o.c.Log())
		return
	}
	// the ack closes the window in which the server was still pushing the channel
	// we had already dropped locally: a push after it is a real routing error again
	if r.IsUnsubscribe() {
		o.subscriptions.unsubscribeConfirmed(r.Arg)
	}
	r.Log(o.c.Log())
	if r.IsError() && o.onError != nil {
		o.onError(r)
	}
}

func (o *WsPublicV2) onTopic(data []byte) error {
	return o.subscriptions.processTopic(data)
}

func (o *WsPublicV2) handleDialError(err error) bool {
	o.ready = false
	if o.onDialError != nil {
		return o.onDialError(err)
	}
	// continue reconnect
	return false
}

// Topic subscriptions

// Candle - candlestick stream of the legacy v2 WS (channel candle<interval>):
// pushed about once per second while trades occur, else once per interval.
// The first push is a snapshot with recent candle history sorted oldest-first
// (the current candle is last).
// Supported intervals: 1m, 3m, 5m, 15m, 30m, 1H, 2H, 4H, 6H, 8H, 12H, 1D, 3D, 1W, 1M
// and the UTC-grid variants 6Hutc, 12Hutc, 1Dutc, 3Dutc, 1Wutc, 1Mutc; 3H does not exist
// (error 30016). Native 6H/12H/1D/3D/1W/1M open on the UTC+8 grid - see BITGET_API.md
// https://www.bitget.com/api-doc/spot/websocket/public/Candlesticks-Channel
func (o *WsPublicV2) Candle(category Category, symbol string, interval Interval) *ExecutorV2[[]WsCandleV2] {
	return NewExecutorV2[[]WsCandleV2](wsArgsV2("candle"+string(interval), category, symbol), o.subscriptions)
}
