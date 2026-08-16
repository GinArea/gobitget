package bitgetv3

import (
	"github.com/msw-x/moon/ulog"
)

// WsPublic - public WebSocket client (UTA)
// A single endpoint serves all product types: the category is passed per subscription
// https://www.bitget.com/api-doc/uta/guide
type WsPublic struct {
	WsBase
}

func NewWsPublic() *WsPublic {
	o := new(WsPublic)
	o.init()
	o.c.WithPath(WsPublicPath)
	return o
}

// Builder methods (return *WsPublic for chaining)

func (o *WsPublic) WithLog(log *ulog.Log) *WsPublic {
	o.setLog(log)
	return o
}

func (o *WsPublic) WithProxy(proxy string) *WsPublic {
	o.setProxy(proxy)
	return o
}

func (o *WsPublic) WithLogRequest(enable bool) *WsPublic {
	o.setLogRequest(enable)
	return o
}

func (o *WsPublic) WithLogResponse(enable bool) *WsPublic {
	o.setLogResponse(enable)
	return o
}

func (o *WsPublic) WithOnDialError(f func(error) bool) *WsPublic {
	o.setOnDialError(f)
	return o
}

func (o *WsPublic) WithOnConnected(f func()) *WsPublic {
	o.setOnConnected(f)
	return o
}

func (o *WsPublic) WithOnDisconnected(f func()) *WsPublic {
	o.setOnDisconnected(f)
	return o
}

func (o *WsPublic) WithOnReady(f func()) *WsPublic {
	o.setOnReady(f)
	return o
}

func (o *WsPublic) WithOnError(f func(WsResponse)) *WsPublic {
	o.setOnError(f)
	return o
}

// Topic subscriptions

// Ticker - live price/24h-stats stream, push every 200-300ms (spot) / 300-400ms (futures)
// https://www.bitget.com/api-doc/uta/websocket/public/Tickers-Channel
func (o *WsPublic) Ticker(category Category, symbol string) *Executor[[]WsTicker] {
	return NewExecutor[[]WsTicker](wsArgs("ticker", category, symbol), o.subscriptions)
}

// There is intentionally no v3 kline factory here: the v3 WS channel lacks the required
// timeframes (no 2H, no utc-grid variants), so candle streaming goes through the legacy
// v2 client (WsPublicV2.Candle) - see BITGET_API.md
