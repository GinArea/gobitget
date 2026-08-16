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
		{
			// interval is added for kline and keeps its case (1D, not 1d),
			// verbatim request example from the Candlesticks Channel docs
			name: "subscribe kline",
			req: WsRequest{
				Op:   "subscribe",
				Args: []SubscriptionArgs{{InstType: "spot", Topic: "kline", Symbol: "BTCUSDT", Interval: "1D"}},
			},
			want: `{"op":"subscribe","args":[{"instType":"spot","topic":"kline","symbol":"BTCUSDT","interval":"1D"}]}`,
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
			// verbatim from https://www.bitget.com/api-doc/uta/websocket/public/Candlesticks-Channel
			// the kline ack echoes the interval in arg
			name:          "kline subscribe ack",
			data:          `{"event":"subscribe","arg":{"instType":"spot","topic":"kline","symbol":"BTCUSDT","interval":"1D"},"connId":"xxxxxxxxxx"}`,
			wantEvent:     true,
			wantSubscribe: true,
			wantArg:       SubscriptionArgs{InstType: "spot", Topic: "kline", Symbol: "BTCUSDT", Interval: "1D"},
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

func TestUnmarshalRawTopicTicker(t *testing.T) {
	// verbatim spot push example from https://www.bitget.com/api-doc/uta/websocket/public/Tickers-Channel
	spotPush := `{
		"data": [
			{
				"bid1Price": "99999",
				"lowPrice24h": "98200",
				"ask1Size": "188.312553",
				"volume24h": "37.722858",
				"price24hPcnt": "0.01833",
				"highPrice24h": "100000",
				"turnover24h": "3750302.979626",
				"bid1Size": "186.183209",
				"ask1Price": "100000",
				"openPrice24h": "0",
				"lastPrice": "100000",
				"platformTurnover24h": "677732572.225658"
			}
		],
		"arg": {"instType": "spot","symbol": "BTCUSDT","topic": "ticker"},
		"action": "snapshot",
		"ts": 1736371332162
	}`
	// futures push carries extra fields (the docs give no futures sample, fields per the parameter table)
	futuresPush := `{
		"data": [
			{
				"lastPrice": "99756.7",
				"bid1Price": "99756.6",
				"ask1Price": "99756.7",
				"indexPrice": "99750.1",
				"markPrice": "99755.0",
				"fundingRate": "0.0001",
				"nextFundingTime": "1746698732562",
				"openInterest": "12345.678",
				"deliveryStatus": "normal"
			}
		],
		"arg": {"instType": "usdt-futures","symbol": "BTCUSDT","topic": "ticker"},
		"action": "snapshot",
		"ts": 1746698732563
	}`

	tests := []struct {
		name    string
		data    string
		wantErr bool
		wantArg SubscriptionArgs
		wantTs  int64
		wantLen int
		check   func(t *testing.T, v WsTicker)
	}{
		{
			name:    "spot snapshot",
			data:    spotPush,
			wantArg: SubscriptionArgs{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"},
			wantTs:  1736371332162,
			wantLen: 1,
			check: func(t *testing.T, v WsTicker) {
				if v.LastPrice.Value() != 100000 {
					t.Fatalf("expected lastPrice 100000, got %v", v.LastPrice.Value())
				}
				if v.Bid1Price.Value() != 99999 {
					t.Fatalf("expected bid1Price 99999, got %v", v.Bid1Price.Value())
				}
				if v.Ask1Price.Value() != 100000 {
					t.Fatalf("expected ask1Price 100000, got %v", v.Ask1Price.Value())
				}
				if v.Bid1Size.Value() != 186.183209 {
					t.Fatalf("expected bid1Size 186.183209, got %v", v.Bid1Size.Value())
				}
				if v.Volume24h.Value() != 37.722858 {
					t.Fatalf("expected volume24h 37.722858, got %v", v.Volume24h.Value())
				}
				if v.Price24hPcnt.Value() != 0.01833 {
					t.Fatalf("expected price24hPcnt 0.01833, got %v", v.Price24hPcnt.Value())
				}
				if v.OpenPrice24h.Value() != 0 {
					t.Fatalf("expected openPrice24h 0, got %v", v.OpenPrice24h.Value())
				}
				if v.MarkPrice.Value() != 0 {
					t.Fatalf("expected empty markPrice on spot, got %v", v.MarkPrice.Value())
				}
			},
		},
		{
			name:    "futures snapshot",
			data:    futuresPush,
			wantArg: SubscriptionArgs{InstType: "usdt-futures", Topic: "ticker", Symbol: "BTCUSDT"},
			wantTs:  1746698732563,
			wantLen: 1,
			check: func(t *testing.T, v WsTicker) {
				if v.MarkPrice.Value() != 99755.0 {
					t.Fatalf("expected markPrice 99755.0, got %v", v.MarkPrice.Value())
				}
				if v.IndexPrice.Value() != 99750.1 {
					t.Fatalf("expected indexPrice 99750.1, got %v", v.IndexPrice.Value())
				}
				if v.FundingRate.Value() != 0.0001 {
					t.Fatalf("expected fundingRate 0.0001, got %v", v.FundingRate.Value())
				}
				if v.NextFundingTime.Value() != 1746698732562 {
					t.Fatalf("expected nextFundingTime 1746698732562, got %v", v.NextFundingTime.Value())
				}
				if v.OpenInterest.Value() != 12345.678 {
					t.Fatalf("expected openInterest 12345.678, got %v", v.OpenInterest.Value())
				}
				if v.DeliveryStatus != "normal" {
					t.Fatalf("expected deliveryStatus normal, got %s", v.DeliveryStatus)
				}
			},
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
			topic, err := UnmarshalRawTopic[[]WsTicker](raw)
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
			if tt.check != nil {
				tt.check(t, topic.Data[0])
			}
		})
	}
}

func TestUnmarshalRawTopicCandle(t *testing.T) {
	// verbatim push example from https://www.bitget.com/api-doc/uta/websocket/public/Candlesticks-Channel
	klinePush := `{
		"data": [
			{
				"volume": "0.423",
				"high": "400005",
				"low": "276670",
				"start": "1710518400000",
				"close": "400005",
				"turnover": "148190.38375",
				"open": "276670"
			}
		],
		"arg": {"instType": "spot","symbol": "BTCUSDT","topic": "kline","interval": "1D"},
		"action": "snapshot",
		"ts": 1736370735556
	}`

	tests := []struct {
		name    string
		data    string
		wantErr bool
		wantArg SubscriptionArgs
		wantTs  int64
		wantLen int
		check   func(t *testing.T, v WsCandle)
	}{
		{
			name:    "spot snapshot",
			data:    klinePush,
			wantArg: SubscriptionArgs{InstType: "spot", Topic: "kline", Symbol: "BTCUSDT", Interval: "1D"},
			wantTs:  1736370735556,
			wantLen: 1,
			check: func(t *testing.T, v WsCandle) {
				if v.Start.Value() != 1710518400000 {
					t.Fatalf("expected start 1710518400000, got %v", v.Start.Value())
				}
				if v.Open.Value() != 276670 {
					t.Fatalf("expected open 276670, got %v", v.Open.Value())
				}
				if v.High.Value() != 400005 {
					t.Fatalf("expected high 400005, got %v", v.High.Value())
				}
				if v.Low.Value() != 276670 {
					t.Fatalf("expected low 276670, got %v", v.Low.Value())
				}
				if v.Close.Value() != 400005 {
					t.Fatalf("expected close 400005, got %v", v.Close.Value())
				}
				if v.Volume.Value() != 0.423 {
					t.Fatalf("expected volume 0.423, got %v", v.Volume.Value())
				}
				if v.Turnover.Value() != 148190.38375 {
					t.Fatalf("expected turnover 148190.38375, got %v", v.Turnover.Value())
				}
			},
		},
		{
			name:    "empty data",
			data:    `{"data":[],"arg":{"instType":"spot","symbol":"BTCUSDT","topic":"kline","interval":"1m"},"action":"snapshot","ts":1736370735556}`,
			wantArg: SubscriptionArgs{InstType: "spot", Topic: "kline", Symbol: "BTCUSDT", Interval: "1m"},
			wantTs:  1736370735556,
			wantLen: 0,
		},
		{
			name:    "invalid payload type",
			data:    `{"data":123,"arg":{"instType":"spot","symbol":"BTCUSDT","topic":"kline","interval":"1m"},"action":"snapshot","ts":1736370735556}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var raw RawTopic
			if err := json.Unmarshal([]byte(tt.data), &raw); err != nil {
				t.Fatalf("unmarshal raw topic failed: %v", err)
			}
			topic, err := UnmarshalRawTopic[[]WsCandle](raw)
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
			if tt.check != nil {
				tt.check(t, topic.Data[0])
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

func klinePush(instType, symbol, interval string) []byte {
	return []byte(`{"data":[{"open":"100000","close":"100001"}],"arg":{"instType":"` + instType +
		`","symbol":"` + symbol + `","topic":"kline","interval":"` + interval +
		`"},"action":"snapshot","ts":1736371332162}`)
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
		NewExecutor[[]WsTicker](argsBtc, s).Subscribe(func(Topic[[]WsTicker]) { got++ })
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

	t.Run("route by interval", func(t *testing.T) {
		args1m := wsArgs("kline", UsdtFutures, "BTCUSDT")
		args1m.Interval = "1m"
		args1d := wsArgs("kline", UsdtFutures, "BTCUSDT")
		args1d.Interval = "1D"
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		var got1m, got1d int
		s.subscribe(args1m, func(RawTopic) error { got1m++; return nil })
		s.subscribe(args1d, func(RawTopic) error { got1d++; return nil })
		if err := s.processTopic(klinePush("usdt-futures", "BTCUSDT", "1m")); err != nil {
			t.Fatalf("expected 1m topic routed, got error: %v", err)
		}
		// the interval echo matches case-insensitively, like the other arg fields
		if err := s.processTopic(klinePush("usdt-futures", "BTCUSDT", "1d")); err != nil {
			t.Fatalf("expected 1d topic routed, got error: %v", err)
		}
		if got1m != 1 || got1d != 1 {
			t.Fatalf("expected each handler called once, got 1m=%d 1d=%d", got1m, got1d)
		}
	})

	t.Run("interval mismatch", func(t *testing.T) {
		args1m := wsArgs("kline", UsdtFutures, "BTCUSDT")
		args1m.Interval = "1m"
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		s.subscribe(args1m, noop)
		if err := s.processTopic(klinePush("usdt-futures", "BTCUSDT", "5m")); err == nil {
			t.Fatal("expected not found error for unsubscribed interval")
		}
	})

	t.Run("kline and ticker independent", func(t *testing.T) {
		argsKline := wsArgs("kline", UsdtFutures, "BTCUSDT")
		argsKline.Interval = "1m"
		c := &fakeSubscriptionClient{ready: true}
		s := NewSubscriptions(c)
		var gotTicker, gotKline int
		s.subscribe(argsBtc, func(RawTopic) error { gotTicker++; return nil })
		s.subscribe(argsKline, func(RawTopic) error { gotKline++; return nil })
		if err := s.processTopic(tickerPush("usdt-futures", "BTCUSDT")); err != nil {
			t.Fatalf("expected ticker topic routed, got error: %v", err)
		}
		if err := s.processTopic(klinePush("usdt-futures", "BTCUSDT", "1m")); err != nil {
			t.Fatalf("expected kline topic routed, got error: %v", err)
		}
		if gotTicker != 1 || gotKline != 1 {
			t.Fatalf("expected each handler called once, got ticker=%d kline=%d", gotTicker, gotKline)
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
