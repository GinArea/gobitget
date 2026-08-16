package bitgetv3

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/ujson"
)

// Offline unit tests for the WebSocket layer (no network)

func TestWsArgs(t *testing.T) {
	tests := []struct {
		name         string
		topic        string
		category     Category
		symbol       string
		wantInstType string
	}{
		{
			name:         "spot",
			topic:        "ticker",
			category:     Spot,
			symbol:       "BTCUSDT",
			wantInstType: "spot",
		},
		{
			name:         "margin",
			topic:        "ticker",
			category:     Margin,
			symbol:       "BTCUSDT",
			wantInstType: "margin",
		},
		{
			name:         "usdt futures",
			topic:        "ticker",
			category:     UsdtFutures,
			symbol:       "BTCUSDT",
			wantInstType: "usdt-futures",
		},
		{
			name:         "coin futures",
			topic:        "books",
			category:     CoinFutures,
			symbol:       "BTCUSD",
			wantInstType: "coin-futures",
		},
		{
			name:         "usdc futures",
			topic:        "ticker",
			category:     UsdcFutures,
			symbol:       "BTCPERP",
			wantInstType: "usdc-futures",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := wsArgs(tt.topic, tt.category, tt.symbol)
			if args.InstType != tt.wantInstType {
				t.Fatalf("expected instType %s, got %s", tt.wantInstType, args.InstType)
			}
			if args.Topic != tt.topic {
				t.Fatalf("expected topic %s, got %s", tt.topic, args.Topic)
			}
			if args.Symbol != tt.symbol {
				t.Fatalf("expected symbol %s, got %s", tt.symbol, args.Symbol)
			}
		})
	}
}

func TestWsRequestMarshal(t *testing.T) {
	tests := []struct {
		name string
		req  WsRequest
		want string
	}{
		{
			name: "subscribe",
			req: WsRequest{
				Op:   "subscribe",
				Args: []SubscriptionArgs{wsArgs("ticker", UsdtFutures, "BTCUSDT")},
			},
			want: `{"op":"subscribe","args":[{"instType":"usdt-futures","topic":"ticker","symbol":"BTCUSDT"}]}`,
		},
		{
			name: "unsubscribe",
			req: WsRequest{
				Op:   "unsubscribe",
				Args: []SubscriptionArgs{wsArgs("ticker", Spot, "BTCUSDT")},
			},
			want: `{"op":"unsubscribe","args":[{"instType":"spot","topic":"ticker","symbol":"BTCUSDT"}]}`,
		},
		{
			name: "subscribe batch",
			req: WsRequest{
				Op: "subscribe",
				Args: []SubscriptionArgs{
					wsArgs("ticker", UsdtFutures, "BTCUSDT"),
					wsArgs("ticker", UsdtFutures, "ETHUSDT"),
				},
			},
			want: `{"op":"subscribe","args":[{"instType":"usdt-futures","topic":"ticker","symbol":"BTCUSDT"},{"instType":"usdt-futures","topic":"ticker","symbol":"ETHUSDT"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// uws.SendJson serializes via ujson.MarshalLowerCase
			b, err := ujson.MarshalLowerCase(tt.req)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if string(b) != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, string(b))
			}
		})
	}
}

func TestWsResponseClassify(t *testing.T) {
	tests := []struct {
		name            string
		data            string
		wantEvent       bool
		wantError       bool
		wantSubscribe   bool
		wantUnsubscribe bool
		wantCode        int64
		wantArg         SubscriptionArgs
	}{
		{
			// verbatim from https://www.bitget.com/api-doc/uta/websocket/public/Tickers-Channel
			name:          "subscribe ack",
			data:          `{"event":"subscribe","arg":{"instType":"spot","topic":"ticker","symbol":"BTCUSDT"},"connId":"xxxxxxxxxx"}`,
			wantEvent:     true,
			wantSubscribe: true,
			wantArg:       SubscriptionArgs{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"},
		},
		{
			name:            "unsubscribe ack",
			data:            `{"event":"unsubscribe","arg":{"instType":"usdt-futures","topic":"ticker","symbol":"BTCUSDT"},"connId":"xxxxxxxxxx"}`,
			wantEvent:       true,
			wantUnsubscribe: true,
			wantArg:         SubscriptionArgs{InstType: "usdt-futures", Topic: "ticker", Symbol: "BTCUSDT"},
		},
		{
			// verbatim error example from https://www.bitget.com/api-doc/uta/guide (code as a string)
			name:      "error event",
			data:      `{"event":"error","code":"30005","msg":"error"}`,
			wantEvent: true,
			wantError: true,
			wantCode:  30005,
		},
		{
			// verbatim live error: the real server sends code as a JSON number, unlike the docs
			name:      "error event numeric code",
			data:      `{"event":"error","code":30001,"msg":"{\"instType\":\"usdt-futures\",\"symbol\":\"NOSUCHSYMBOL\",\"topic\":\"ticker\"} doesn't exist","connId":"0aa7eefffe582cd1"}`,
			wantEvent: true,
			wantError: true,
			wantCode:  30001,
		},
		{
			// ticker push has no "event" field and must not be treated as one
			name:    "ticker push",
			data:    `{"data":[{"lastPrice":"100000"}],"arg":{"instType":"spot","symbol":"BTCUSDT","topic":"ticker"},"action":"snapshot","ts":1736371332162}`,
			wantArg: SubscriptionArgs{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r WsResponse
			if err := json.Unmarshal([]byte(tt.data), &r); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if r.IsEvent() != tt.wantEvent {
				t.Fatalf("expected IsEvent %v, got %v", tt.wantEvent, r.IsEvent())
			}
			if r.IsError() != tt.wantError {
				t.Fatalf("expected IsError %v, got %v", tt.wantError, r.IsError())
			}
			if r.IsSubscribe() != tt.wantSubscribe {
				t.Fatalf("expected IsSubscribe %v, got %v", tt.wantSubscribe, r.IsSubscribe())
			}
			if r.IsUnsubscribe() != tt.wantUnsubscribe {
				t.Fatalf("expected IsUnsubscribe %v, got %v", tt.wantUnsubscribe, r.IsUnsubscribe())
			}
			if r.Code.Value() != tt.wantCode {
				t.Fatalf("expected code %d, got %v", tt.wantCode, r.Code.Value())
			}
			if (tt.wantCode == 30001) != r.TopicNotExists() {
				t.Fatalf("expected TopicNotExists %v, got %v", tt.wantCode == 30001, r.TopicNotExists())
			}
			if r.Arg != tt.wantArg {
				t.Fatalf("expected arg %+v, got %+v", tt.wantArg, r.Arg)
			}
		})
	}
}

// testShot - neutral payload for infrastructure tests of the topic envelope
type testShot struct {
	LastPrice ujson.Float64
}

func TestUnmarshalRawTopic(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
		wantArg SubscriptionArgs
		wantTs  int64
		wantLen int
	}{
		{
			name:    "snapshot",
			data:    `{"data":[{"lastPrice":"100000"}],"arg":{"instType":"spot","symbol":"BTCUSDT","topic":"ticker"},"action":"snapshot","ts":1736371332162}`,
			wantArg: SubscriptionArgs{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"},
			wantTs:  1736371332162,
			wantLen: 1,
		},
		{
			name:    "empty data",
			data:    `{"data":[],"arg":{"instType":"spot","symbol":"BTCUSDT","topic":"ticker"},"action":"snapshot","ts":1736371332162}`,
			wantArg: SubscriptionArgs{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"},
			wantTs:  1736371332162,
			wantLen: 0,
		},
		{
			name:    "invalid payload type",
			data:    `{"data":123,"arg":{"instType":"spot","symbol":"BTCUSDT","topic":"ticker"},"action":"snapshot","ts":1736371332162}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var raw RawTopic
			if err := json.Unmarshal([]byte(tt.data), &raw); err != nil {
				t.Fatalf("unmarshal raw topic failed: %v", err)
			}
			topic, err := UnmarshalRawTopic[[]testShot](raw)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got success")
				}
				return
			}
			if err != nil {
				t.Fatalf("unmarshal topic failed: %v", err)
			}
			if topic.Arg != tt.wantArg {
				t.Fatalf("expected arg %+v, got %+v", tt.wantArg, topic.Arg)
			}
			if topic.Action != "snapshot" {
				t.Fatalf("expected action snapshot, got %s", topic.Action)
			}
			if topic.Ts.Value() != tt.wantTs {
				t.Fatalf("expected ts %d, got %v", tt.wantTs, topic.Ts.Value())
			}
			if len(topic.Data) != tt.wantLen {
				t.Fatalf("expected %d items, got %d", tt.wantLen, len(topic.Data))
			}
			if tt.wantLen > 0 && topic.Data[0].LastPrice.Value() != 100000 {
				t.Fatalf("expected lastPrice 100000, got %v", topic.Data[0].LastPrice.Value())
			}
		})
	}
}

// fakeSubscriptionClient - records subscribe/unsubscribe calls for Subscriptions tests
type fakeSubscriptionClient struct {
	ready        bool
	subscribes   [][]SubscriptionArgs
	unsubscribes [][]SubscriptionArgs
}

func (o *fakeSubscriptionClient) Ready() bool {
	return o.ready
}

func (o *fakeSubscriptionClient) subscribe(args ...SubscriptionArgs) {
	o.subscribes = append(o.subscribes, args)
}

func (o *fakeSubscriptionClient) unsubscribe(args ...SubscriptionArgs) {
	o.unsubscribes = append(o.unsubscribes, args)
}

func tickerPush(instType, symbol string) []byte {
	return []byte(`{"data":[{"lastPrice":"100000"}],"arg":{"instType":"` + instType +
		`","symbol":"` + symbol + `","topic":"ticker"},"action":"snapshot","ts":1736371332162}`)
}

func TestSubscriptions(t *testing.T) {
	argsBtc := wsArgs("ticker", UsdtFutures, "BTCUSDT")
	argsEth := wsArgs("ticker", UsdtFutures, "ETHUSDT")
	noop := func(RawTopic) error { return nil }

	t.Run("deferred until ready", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: false}
		s := NewSubscriptions(c)
		s.subscribe(argsBtc, noop)
		if len(c.subscribes) != 0 {
			t.Fatalf("expected no sends while not ready, got %d", len(c.subscribes))
		}
		s.subscribeAll()
		if len(c.subscribes) != 1 || len(c.subscribes[0]) != 1 || c.subscribes[0][0] != argsBtc {
			t.Fatalf("expected one batch with one arg, got %+v", c.subscribes)
		}
	})

	t.Run("immediate when ready", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		s.subscribe(argsBtc, noop)
		if len(c.subscribes) != 1 || c.subscribes[0][0] != argsBtc {
			t.Fatalf("expected immediate send, got %+v", c.subscribes)
		}
	})

	t.Run("subscribe all batches", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: false}
		s := NewSubscriptions(c)
		s.subscribe(argsBtc, noop)
		s.subscribe(argsEth, noop)
		s.subscribeAll()
		if len(c.subscribes) != 1 {
			t.Fatalf("expected a single batch request, got %d", len(c.subscribes))
		}
		batch := c.subscribes[0]
		if len(batch) != 2 {
			t.Fatalf("expected 2 args in batch, got %d", len(batch))
		}
		found := map[SubscriptionArgs]bool{}
		for _, a := range batch {
			found[a] = true
		}
		if !found[argsBtc] || !found[argsEth] {
			t.Fatalf("expected both args in batch, got %+v", batch)
		}
	})

	t.Run("subscribe all empty", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		s.subscribeAll()
		if len(c.subscribes) != 0 {
			t.Fatalf("expected no requests with no subscriptions, got %d", len(c.subscribes))
		}
	})

	t.Run("unsubscribe removes handler", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		s.subscribe(argsBtc, noop)
		if err := s.processTopic(tickerPush("usdt-futures", "BTCUSDT")); err != nil {
			t.Fatalf("expected topic routed, got error: %v", err)
		}
		s.unsubscribe(argsBtc)
		if len(c.unsubscribes) != 1 || c.unsubscribes[0][0] != argsBtc {
			t.Fatalf("expected unsubscribe sent, got %+v", c.unsubscribes)
		}
		if err := s.processTopic(tickerPush("usdt-futures", "BTCUSDT")); err == nil {
			t.Fatal("expected not found error after unsubscribe")
		}
	})

	t.Run("route case-insensitive", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		var got int
		s.subscribe(argsBtc, func(RawTopic) error { got++; return nil })
		if err := s.processTopic(tickerPush("USDT-FUTURES", "btcusdt")); err != nil {
			t.Fatalf("expected case-insensitive match, got error: %v", err)
		}
		if got != 1 {
			t.Fatalf("expected handler called once, got %d", got)
		}
	})

	t.Run("unknown arg", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		s.subscribe(argsBtc, noop)
		if err := s.processTopic(tickerPush("usdt-futures", "XRPUSDT")); err == nil {
			t.Fatal("expected not found error for unknown symbol")
		}
	})

	t.Run("handler error propagates", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		sentinel := errors.New("handler failed")
		s.subscribe(argsBtc, func(RawTopic) error { return sentinel })
		if err := s.processTopic(tickerPush("usdt-futures", "BTCUSDT")); !errors.Is(err, sentinel) {
			t.Fatalf("expected sentinel error, got %v", err)
		}
	})

	t.Run("payload unmarshal error propagates", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		var got int
		NewExecutor[[]testShot](argsBtc, s).Subscribe(func(Topic[[]testShot]) { got++ })
		bad := []byte(`{"data":123,"arg":{"instType":"usdt-futures","symbol":"BTCUSDT","topic":"ticker"},"action":"snapshot","ts":1}`)
		if err := s.processTopic(bad); err == nil {
			t.Fatal("expected unmarshal error")
		}
		if got != 0 {
			t.Fatalf("expected handler not called on bad payload, got %d", got)
		}
	})

	t.Run("invalid topic json", func(t *testing.T) {
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		if err := s.processTopic([]byte(`not json`)); err == nil {
			t.Fatal("expected error on invalid json")
		}
	})

}

func TestWsClientOnMessage(t *testing.T) {
	newClient := func() (*WsClient, *[]WsResponse, *[][]byte) {
		c := NewWsClient()
		responses := new([]WsResponse)
		topics := new([][]byte)
		c.WithOnResponse(func(r WsResponse) { *responses = append(*responses, r) })
		c.WithOnTopic(func(d []byte) error { *topics = append(*topics, d); return nil })
		return c, responses, topics
	}

	tests := []struct {
		name          string
		messageType   int
		data          string
		wantResponses int
		wantTopics    int
	}{
		{
			name:        "pong swallowed",
			messageType: websocket.TextMessage,
			data:        "pong",
		},
		{
			name:        "binary ignored",
			messageType: websocket.BinaryMessage,
			data:        `{"event":"subscribe"}`,
		},
		{
			name:        "invalid json ignored",
			messageType: websocket.TextMessage,
			data:        "not json",
		},
		{
			name:          "event routed to response",
			messageType:   websocket.TextMessage,
			data:          `{"event":"subscribe","arg":{"instType":"spot","topic":"ticker","symbol":"BTCUSDT"},"connId":"x"}`,
			wantResponses: 1,
		},
		{
			name:          "error routed to response",
			messageType:   websocket.TextMessage,
			data:          `{"event":"error","code":"30005","msg":"error"}`,
			wantResponses: 1,
		},
		{
			name:        "push routed to topic",
			messageType: websocket.TextMessage,
			data:        `{"data":[{"lastPrice":"100000"}],"arg":{"instType":"spot","symbol":"BTCUSDT","topic":"ticker"},"action":"snapshot","ts":1736371332162}`,
			wantTopics:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, responses, topics := newClient()
			c.onMessage(tt.messageType, []byte(tt.data))
			if len(*responses) != tt.wantResponses {
				t.Fatalf("expected %d responses, got %d", tt.wantResponses, len(*responses))
			}
			if len(*topics) != tt.wantTopics {
				t.Fatalf("expected %d topics, got %d", tt.wantTopics, len(*topics))
			}
		})
	}
}
