package bitgetv3

import (
	"bytes"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

// WsClientV2 - protocol-level client of the legacy v2 WebSocket: framing, keepalive,
// message classification; the framing and keepalive match the v3 WsClient, only the
// arg/event shapes differ
type WsClientV2 struct {
	c          *uws.Client
	onResponse func(WsResponseV2)
	onTopic    func([]byte) error
}

func NewWsClientV2() *WsClientV2 {
	o := new(WsClientV2)
	o.c = uws.NewClient(MainBaseWsUrl)
	return o
}

func (o *WsClientV2) Close() {
	o.c.Close()
}

func (o *WsClientV2) Log() *ulog.Log {
	return o.c.Log()
}

func (o *WsClientV2) Transport() *uws.Options {
	return &o.c.Options
}

func (o *WsClientV2) WithLog(log *ulog.Log) *WsClientV2 {
	o.c.WithLog(log)
	return o
}

func (o *WsClientV2) WithPath(path string) *WsClientV2 {
	o.c.WithPath(path)
	return o
}

func (o *WsClientV2) WithProxy(proxy string) *WsClientV2 {
	o.c.WithProxy(proxy)
	return o
}

func (o *WsClientV2) WithLogRequest(enable bool) *WsClientV2 {
	o.Transport().LogSent.Size = enable
	o.Transport().LogSent.Data = enable
	return o
}

func (o *WsClientV2) WithLogResponse(enable bool) *WsClientV2 {
	o.Transport().LogRecv.Size = enable
	o.Transport().LogRecv.Data = enable
	return o
}

func (o *WsClientV2) WithOnDialError(f func(error) bool) *WsClientV2 {
	o.c.WithOnDialError(f)
	return o
}

func (o *WsClientV2) WithOnConnected(f func()) *WsClientV2 {
	o.c.WithOnConnected(f)
	return o
}

func (o *WsClientV2) WithOnDisconnected(f func()) *WsClientV2 {
	o.c.WithOnDisconnected(f)
	return o
}

func (o *WsClientV2) WithOnResponse(f func(WsResponseV2)) *WsClientV2 {
	o.onResponse = f
	return o
}

func (o *WsClientV2) WithOnTopic(f func([]byte) error) *WsClientV2 {
	o.onTopic = f
	return o
}

func (o *WsClientV2) Run() {
	o.c.WithOnPing(o.ping)
	o.c.WithOnMessage(o.onMessage)
	o.c.Run()
}

func (o *WsClientV2) Running() bool {
	return o.c.Running()
}

func (o *WsClientV2) Connected() bool {
	return o.c.Connected()
}

func (o *WsClientV2) Reconnect() {
	o.c.Reconnect()
}

func (o *WsClientV2) Send(r WsRequestV2) {
	o.c.SendJson(r)
}

// Subscribe - subscribe to channels, batched like the v3 client
func (o *WsClientV2) Subscribe(args ...SubscriptionArgsV2) {
	if len(args) > 0 {
		o.Send(WsRequestV2{Op: "subscribe", Args: args})
	}
}

// Unsubscribe - unsubscribe from channels, batched like Subscribe
func (o *WsClientV2) Unsubscribe(args ...SubscriptionArgsV2) {
	if len(args) > 0 {
		o.Send(WsRequestV2{Op: "unsubscribe", Args: args})
	}
}

// ping - the v2 keepalive is the same literal text frame as v3 (verified live)
func (o *WsClientV2) ping() {
	o.c.SendText("ping")
}

func (o *WsClientV2) onMessage(messageType int, data []byte) {
	log := o.c.Log()
	if messageType != websocket.TextMessage {
		log.Warning("invalid message type:", uws.MessageTypeString(messageType))
		return
	}
	if bytes.Equal(data, []byte("pong")) {
		return
	}
	var r WsResponseV2
	err := json.Unmarshal(data, &r)
	if err == nil {
		if r.IsEvent() {
			if o.onResponse != nil {
				o.onResponse(r)
			}
		} else {
			if o.onTopic != nil {
				err = o.onTopic(data)
			}
		}
	}
	if err != nil {
		log.Error(err)
	}
}
