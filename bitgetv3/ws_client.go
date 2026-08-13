package bitgetv3

import (
	"bytes"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

// WsClient - protocol-level WebSocket client: framing, keepalive, message classification
type WsClient struct {
	c          *uws.Client
	onResponse func(WsResponse)
	onTopic    func([]byte) error
}

func NewWsClient() *WsClient {
	o := new(WsClient)
	o.c = uws.NewClient(MainBaseWsUrl)
	return o
}

func (o *WsClient) Close() {
	o.c.Close()
}

func (o *WsClient) Log() *ulog.Log {
	return o.c.Log()
}

func (o *WsClient) Transport() *uws.Options {
	return &o.c.Options
}

func (o *WsClient) WithLog(log *ulog.Log) *WsClient {
	o.c.WithLog(log)
	return o
}

func (o *WsClient) WithPath(path string) *WsClient {
	o.c.WithPath(path)
	return o
}

func (o *WsClient) WithProxy(proxy string) *WsClient {
	o.c.WithProxy(proxy)
	return o
}

func (o *WsClient) WithLogRequest(enable bool) *WsClient {
	o.Transport().LogSent.Size = enable
	o.Transport().LogSent.Data = enable
	return o
}

func (o *WsClient) WithLogResponse(enable bool) *WsClient {
	o.Transport().LogRecv.Size = enable
	o.Transport().LogRecv.Data = enable
	return o
}

func (o *WsClient) WithOnDialError(f func(error) bool) *WsClient {
	o.c.WithOnDialError(f)
	return o
}

func (o *WsClient) WithOnConnected(f func()) *WsClient {
	o.c.WithOnConnected(f)
	return o
}

func (o *WsClient) WithOnDisconnected(f func()) *WsClient {
	o.c.WithOnDisconnected(f)
	return o
}

func (o *WsClient) WithOnResponse(f func(WsResponse)) *WsClient {
	o.onResponse = f
	return o
}

func (o *WsClient) WithOnTopic(f func([]byte) error) *WsClient {
	o.onTopic = f
	return o
}

func (o *WsClient) Run() {
	o.c.WithOnPing(o.ping)
	o.c.WithOnMessage(o.onMessage)
	o.c.Run()
}

func (o *WsClient) Running() bool {
	return o.c.Running()
}

func (o *WsClient) Connected() bool {
	return o.c.Connected()
}

func (o *WsClient) Reconnect() {
	o.c.Reconnect()
}

func (o *WsClient) Send(r WsRequest) {
	o.c.SendJson(r)
}

// Subscribe - subscribe to channels; batching in one message saves the
// per-connection limits (10 messages/s, 240 subscription requests/hour)
func (o *WsClient) Subscribe(args ...SubscriptionArgs) {
	if len(args) > 0 {
		o.Send(WsRequest{Op: "subscribe", Args: args})
	}
}

// Unsubscribe - unsubscribe from channels, batched like Subscribe
func (o *WsClient) Unsubscribe(args ...SubscriptionArgs) {
	if len(args) > 0 {
		o.Send(WsRequest{Op: "unsubscribe", Args: args})
	}
}

// ping - Bitget keepalive is a literal text frame, not a JSON message;
// the server disconnects a client it got no ping from for 2 minutes
func (o *WsClient) ping() {
	o.c.SendText("ping")
}

func (o *WsClient) onMessage(messageType int, data []byte) {
	log := o.c.Log()
	if messageType != websocket.TextMessage {
		log.Warning("invalid message type:", uws.MessageTypeString(messageType))
		return
	}
	if bytes.Equal(data, []byte("pong")) {
		return
	}
	var r WsResponse
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
