package bitgetv3

import (
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

// WsBase - common WebSocket functionality for public and private clients
type WsBase struct {
	c              *WsClient
	ready          bool
	authSent       bool
	handshake      func()
	onReady        func()
	onConnected    func()
	onDisconnected func()
	onDialError    func(error) bool
	onError        func(WsResponse)
	onLoginFailed  func()
	subscriptions  *Subscriptions
}

func (o *WsBase) init() {
	o.c = NewWsClient()
	o.subscriptions = NewSubscriptions(o)
}

func (o *WsBase) Close() {
	o.c.Close()
}

func (o *WsBase) Transport() *uws.Options {
	return o.c.Transport()
}

func (o *WsBase) Running() bool {
	return o.c.Running()
}

func (o *WsBase) Connected() bool {
	return o.c.Connected()
}

func (o *WsBase) Reconnect() {
	o.c.Reconnect()
}

func (o *WsBase) Ready() bool {
	return o.ready
}

func (o *WsBase) Run() {
	o.c.WithOnConnected(func() {
		if o.onConnected != nil {
			o.onConnected()
		}
		// the public channel has no handshake (welcome/login): ready right after connect;
		// the private client sets handshake to send the login request, and markReady()
		// is called from onResponse on a successful login ack
		if o.handshake == nil {
			o.markReady()
		} else {
			o.authSent = true
			o.handshake()
		}
	})
	o.c.WithOnDisconnected(func() {
		o.ready = false
		o.authSent = false
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
func (o *WsBase) markReady() {
	o.ready = true
	o.subscriptions.subscribeAll()
	if o.onReady != nil {
		o.onReady()
	}
}

func (o *WsBase) subscribe(args ...SubscriptionArgs) {
	o.c.Subscribe(args...)
}

func (o *WsBase) unsubscribe(args ...SubscriptionArgs) {
	o.c.Unsubscribe(args...)
}

func (o *WsBase) onResponse(r WsResponse) {
	if r.IsLogin() {
		o.authSent = false
		if r.Code.Value() == 0 {
			r.Log(o.c.Log())
			o.markReady()
		} else {
			o.loginFailed(r)
		}
		return
	}
	// an error event in the handshake window is the login rejection in its
	// event:"error" shape: no other request may be in flight before ready
	// (subscriptions are gated on Ready)
	if r.IsError() && o.authSent && !o.ready {
		o.authSent = false
		o.loginFailed(r)
		return
	}
	r.Log(o.c.Log())
	if r.IsError() && o.onError != nil {
		o.onError(r)
	}
}

// loginFailed - rejected credentials are unrecoverable: stop the reconnect loop
// instead of retrying the same bad login forever, then notify the client
func (o *WsBase) loginFailed(r WsResponse) {
	o.c.Log().Error("login:", r.Code.Value(), r.Msg)
	o.c.Cancel()
	if o.onLoginFailed != nil {
		o.onLoginFailed()
	}
}

func (o *WsBase) onTopic(data []byte) error {
	return o.subscriptions.processTopic(data)
}

func (o *WsBase) handleDialError(err error) bool {
	o.ready = false
	if o.onDialError != nil {
		return o.onDialError(err)
	}
	// continue reconnect
	return false
}

// Setters for callbacks (used by With* methods in derived types)

func (o *WsBase) setLog(log *ulog.Log) {
	o.c.WithLog(log)
}

func (o *WsBase) setProxy(proxy string) {
	o.c.WithProxy(proxy)
}

func (o *WsBase) setLogRequest(enable bool) {
	o.c.WithLogRequest(enable)
}

func (o *WsBase) setLogResponse(enable bool) {
	o.c.WithLogResponse(enable)
}

func (o *WsBase) setOnDialError(f func(error) bool) {
	o.onDialError = f
}

func (o *WsBase) setOnConnected(f func()) {
	o.onConnected = f
}

func (o *WsBase) setOnDisconnected(f func()) {
	o.onDisconnected = f
}

func (o *WsBase) setOnReady(f func()) {
	o.onReady = f
}

func (o *WsBase) setOnError(f func(WsResponse)) {
	o.onError = f
}

func (o *WsBase) setHandshake(f func()) {
	o.handshake = f
}

func (o *WsBase) setOnLoginFailed(f func()) {
	o.onLoginFailed = f
}
