package bitgetv3

import (
	"github.com/msw-x/moon/ulog"
)

// WsPrivate - private WebSocket client (UTA)
// Authenticates with a login request right after connect; subscriptions are
// deferred until the login ack and replayed automatically on reconnect
// https://www.bitget.com/api-doc/uta/websocket/private/WebSocket-Private
type WsPrivate struct {
	WsBase
	s *Sign
}

func NewWsPrivate(key, secret, password string) *WsPrivate {
	o := new(WsPrivate)
	o.init()
	o.s = NewSign(key, secret, password)
	o.c.WithPath(WsPrivatePath)
	o.setHandshake(func() {
		o.c.Login(o.s.WebSocket())
	})
	return o
}

// Builder methods (return *WsPrivate for chaining)

func (o *WsPrivate) WithLog(log *ulog.Log) *WsPrivate {
	o.setLog(log)
	return o
}

func (o *WsPrivate) WithProxy(proxy string) *WsPrivate {
	o.setProxy(proxy)
	return o
}

func (o *WsPrivate) WithLogRequest(enable bool) *WsPrivate {
	o.setLogRequest(enable)
	return o
}

func (o *WsPrivate) WithLogResponse(enable bool) *WsPrivate {
	o.setLogResponse(enable)
	return o
}

func (o *WsPrivate) WithOnDialError(f func(error) bool) *WsPrivate {
	o.setOnDialError(f)
	return o
}

func (o *WsPrivate) WithOnConnected(f func()) *WsPrivate {
	o.setOnConnected(f)
	return o
}

func (o *WsPrivate) WithOnDisconnected(f func()) *WsPrivate {
	o.setOnDisconnected(f)
	return o
}

func (o *WsPrivate) WithOnReady(f func()) *WsPrivate {
	o.setOnReady(f)
	return o
}

func (o *WsPrivate) WithOnError(f func(WsResponse)) *WsPrivate {
	o.setOnError(f)
	return o
}

// WithOnLoginFailed - called when the server rejects the login credentials;
// the reconnect loop is stopped before the callback (see WsBase.loginFailed)
func (o *WsPrivate) WithOnLoginFailed(f func()) *WsPrivate {
	o.setOnLoginFailed(f)
	return o
}

// Topic subscriptions

// Position - futures position stream: a snapshot on subscription (action
// "snapshot", empty when flat), then event pushes on position opens/closes and
// close-order events with action "update" (verified live; the docs show
// "snapshot" only, and the order channel uses "snapshot" for events)
// https://www.bitget.com/api-doc/uta/websocket/private/Positions-Channel
func (o *WsPrivate) Position() *Executor[[]WsPosition] {
	return NewExecutor[[]WsPosition](wsPrivateArgs("position"), o.subscriptions)
}

// Order - unified-account order stream: pushes on order place/fill/cancel across
// spot, margin and futures; unlike the position channel there is NO snapshot on
// subscription (verified live). A fast-filling order may arrive as a single
// "filled" push or as a "new" push followed by a "filled" one
// https://www.bitget.com/api-doc/uta/websocket/private/Order-Channel
func (o *WsPrivate) Order() *Executor[[]WsOrder] {
	return NewExecutor[[]WsOrder](wsPrivateArgs("order"), o.subscriptions)
}

// Account - unified-account balance stream: a snapshot on subscription (action
// "snapshot"), then pushes on order fills, fund settlement and balance changes
// (transfers, airdrops, loans, etc.) with action "update" (verified live; like
// the position channel and unlike the order channel, which uses "snapshot" for
// events). Values are absolute balances, not deltas
// https://www.bitget.com/api-doc/uta/websocket/private/Account-Channel
func (o *WsPrivate) Account() *Executor[[]WsAccount] {
	return NewExecutor[[]WsAccount](wsPrivateArgs("account"), o.subscriptions)
}
