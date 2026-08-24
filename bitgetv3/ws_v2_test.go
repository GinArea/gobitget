package bitgetv3

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/ujson"
)

// Offline unit tests of the legacy v2 WS stack (no network access)

func TestWsRequestV2Marshal(t *testing.T) {
	tests := []struct {
		name string
		req  WsRequestV2
		want string
	}{
		{
			// instType stays uppercase and the timeframe is embedded in the channel name,
			// unlike the v3 request shape
			name: "subscribe",
			req: WsRequestV2{
				Op:   "subscribe",
				Args: []SubscriptionArgsV2{wsArgsV2("candle1Dutc", UsdtFutures, "BTCUSDT")},
			},
			want: `{"op":"subscribe","args":[{"instType":"USDT-FUTURES","channel":"candle1Dutc","instId":"BTCUSDT"}]}`,
		},
		{
			name: "unsubscribe",
			req: WsRequestV2{
				Op:   "unsubscribe",
				Args: []SubscriptionArgsV2{wsArgsV2("candle2H", Spot, "BTCUSDT")},
			},
			want: `{"op":"unsubscribe","args":[{"instType":"SPOT","channel":"candle2H","instId":"BTCUSDT"}]}`,
		},
		{
			name: "subscribe batch",
			req: WsRequestV2{
				Op: "subscribe",
				Args: []SubscriptionArgsV2{
					wsArgsV2("candle2H", UsdtFutures, "BTCUSDT"),
					wsArgsV2("candle12Hutc", UsdtFutures, "ETHUSDT"),
				},
			},
			want: `{"op":"subscribe","args":[{"instType":"USDT-FUTURES","channel":"candle2H","instId":"BTCUSDT"},{"instType":"USDT-FUTURES","channel":"candle12Hutc","instId":"ETHUSDT"}]}`,
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

func TestWsResponseV2Classify(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantEvent string
		wantErr   bool
		wantParam bool
		wantArg   SubscriptionArgsV2
		wantCode  int64
	}{
		{
			name:      "subscribe ack",
			data:      `{"event":"subscribe","arg":{"instType":"USDT-FUTURES","channel":"candle2H","instId":"BTCUSDT"}}`,
			wantEvent: "subscribe",
			wantArg:   SubscriptionArgsV2{InstType: "USDT-FUTURES", Channel: "candle2H", InstId: "BTCUSDT"},
		},
		{
			name:      "unsubscribe ack",
			data:      `{"event":"unsubscribe","arg":{"instType":"SPOT","channel":"candle1Dutc","instId":"BTCUSDT"}}`,
			wantEvent: "unsubscribe",
			wantArg:   SubscriptionArgsV2{InstType: "SPOT", Channel: "candle1Dutc", InstId: "BTCUSDT"},
		},
		{
			// verbatim shape observed live: numeric code, echoed op
			name:      "param error",
			data:      `{"event":"error","code":30016,"msg":"Param error","op":"subscribe"}`,
			wantEvent: "error",
			wantErr:   true,
			wantParam: true,
			wantCode:  30016,
		},
		{
			name:      "other error",
			data:      `{"event":"error","code":30001,"msg":"channel does not exist"}`,
			wantEvent: "error",
			wantErr:   true,
			wantCode:  30001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r WsResponseV2
			if err := json.Unmarshal([]byte(tt.data), &r); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if !r.IsEvent() {
				t.Fatal("expected event message")
			}
			if r.Event != tt.wantEvent {
				t.Fatalf("expected event %s, got %s", tt.wantEvent, r.Event)
			}
			if r.IsError() != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, r.IsError())
			}
			if r.ParamError() != tt.wantParam {
				t.Fatalf("expected param error=%v, got %v", tt.wantParam, r.ParamError())
			}
			if r.Arg != tt.wantArg {
				t.Fatalf("expected arg %+v, got %+v", tt.wantArg, r.Arg)
			}
			if r.Code.Value() != tt.wantCode {
				t.Fatalf("expected code %d, got %d", tt.wantCode, r.Code.Value())
			}
		})
	}
}

func TestWsClientV2OnMessage(t *testing.T) {
	newClient := func() (*WsClientV2, *[]WsResponseV2, *[][]byte) {
		c := NewWsClientV2()
		responses := new([]WsResponseV2)
		topics := new([][]byte)
		c.WithOnResponse(func(r WsResponseV2) { *responses = append(*responses, r) })
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
			data:          `{"event":"subscribe","arg":{"instType":"USDT-FUTURES","channel":"candle2H","instId":"BTCUSDT"}}`,
			wantResponses: 1,
		},
		{
			name:          "error routed to response",
			messageType:   websocket.TextMessage,
			data:          `{"event":"error","code":30016,"msg":"Param error","op":"subscribe"}`,
			wantResponses: 1,
		},
		{
			name:        "push routed to topic",
			messageType: websocket.TextMessage,
			data:        candleV2Push("USDT-FUTURES", "candle1m", "BTCUSDT"),
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

func TestUnmarshalRawTopicCandleV2(t *testing.T) {
	t.Run("live push", func(t *testing.T) {
		// verbatim row captured live from candle1m USDT-FUTURES BTCUSDT
		data := `{"action":"snapshot","arg":{"instType":"USDT-FUTURES","channel":"candle1m","instId":"BTCUSDT"},"data":[["1786886340000","63039","63039","63038.9","63039","1.4653","92371.02404","92371.02404"]],"ts":1786886340123}`
		var raw RawTopicV2
		if err := json.Unmarshal([]byte(data), &raw); err != nil {
			t.Fatalf("unmarshal raw failed: %v", err)
		}
		topic, err := UnmarshalRawTopicV2[[]WsCandleV2](raw)
		if err != nil {
			t.Fatalf("unmarshal topic failed: %v", err)
		}
		want := SubscriptionArgsV2{InstType: "USDT-FUTURES", Channel: "candle1m", InstId: "BTCUSDT"}
		if topic.Arg != want {
			t.Fatalf("expected arg %+v, got %+v", want, topic.Arg)
		}
		if topic.Action != "snapshot" {
			t.Fatalf("expected action snapshot, got %s", topic.Action)
		}
		if topic.Ts.Value() != 1786886340123 {
			t.Fatalf("expected ts 1786886340123, got %v", topic.Ts.Value())
		}
		if len(topic.Data) != 1 {
			t.Fatalf("expected 1 candle, got %d", len(topic.Data))
		}
		d := topic.Data[0]
		if d.Start.Value() != 1786886340000 {
			t.Fatalf("expected start 1786886340000, got %v", d.Start.Value())
		}
		if d.Open.Value() != 63039 || d.High.Value() != 63039 || d.Low.Value() != 63038.9 || d.Close.Value() != 63039 {
			t.Fatalf("unexpected ohlc: %+v", d)
		}
		if d.Volume.Value() != 1.4653 {
			t.Fatalf("expected volume 1.4653, got %v", d.Volume.Value())
		}
		if d.QuoteVolume.Value() != 92371.02404 || d.UsdtVolume.Value() != 92371.02404 {
			t.Fatalf("unexpected volumes: %+v", d)
		}
	})

	t.Run("distinct quote and usdt volumes", func(t *testing.T) {
		// verbatim row captured live from candle1H SPOT ETHBTC (BTC-quoted pair)
		row := `["1786878000000","0.02985","0.02987","0.02984","0.02984","0.111","0.003313935","208.73045398515"]`
		var d WsCandleV2
		if err := json.Unmarshal([]byte(row), &d); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if d.Volume.Value() != 0.111 {
			t.Fatalf("expected base volume 0.111, got %v", d.Volume.Value())
		}
		if d.QuoteVolume.Value() != 0.003313935 {
			t.Fatalf("expected quote volume 0.003313935, got %v", d.QuoteVolume.Value())
		}
		if d.UsdtVolume.Value() != 208.73045398515 {
			t.Fatalf("expected usdt volume 208.73045398515, got %v", d.UsdtVolume.Value())
		}
	})

	t.Run("empty data", func(t *testing.T) {
		raw := RawTopicV2{Data: json.RawMessage(`[]`)}
		topic, err := UnmarshalRawTopicV2[[]WsCandleV2](raw)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(topic.Data) != 0 {
			t.Fatalf("expected empty data, got %d", len(topic.Data))
		}
	})

	t.Run("invalid rows", func(t *testing.T) {
		tests := []struct {
			name string
			row  string
		}{
			{
				name: "short row",
				row:  `["1786886340000","63039"]`,
			},
			{
				name: "long row",
				row:  `["1786886340000","1","2","3","4","5","6","7","8"]`,
			},
			{
				name: "invalid number",
				row:  `["1786886340000","oops","2","3","4","5","6","7"]`,
			},
			{
				name: "not an array",
				row:  `{"start":"1786886340000"}`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var d WsCandleV2
				if err := json.Unmarshal([]byte(tt.row), &d); err == nil {
					t.Fatalf("expected error, got %+v", d)
				}
			})
		}
	})
}

// candleV2Push - build a minimal valid v2 candle push message
func candleV2Push(instType, channel, instId string) string {
	return fmt.Sprintf(`{"action":"update","arg":{"instType":"%s","channel":"%s","instId":"%s"},"data":[["1786886340000","63039","63039","63038.9","63039","1.4653","92371.02404","92371.02404"]],"ts":1786886340123}`,
		instType, channel, instId)
}

type fakeSubscriptionClientV2 struct {
	ready        bool
	subscribes   [][]SubscriptionArgsV2
	unsubscribes [][]SubscriptionArgsV2
}

func (o *fakeSubscriptionClientV2) Ready() bool { return o.ready }
func (o *fakeSubscriptionClientV2) subscribe(args ...SubscriptionArgsV2) {
	o.subscribes = append(o.subscribes, args)
}
func (o *fakeSubscriptionClientV2) unsubscribe(args ...SubscriptionArgsV2) {
	o.unsubscribes = append(o.unsubscribes, args)
}

func TestSubscriptionsV2(t *testing.T) {
	args2H := wsArgsV2("candle2H", UsdtFutures, "BTCUSDT")
	args1D := wsArgsV2("candle1Dutc", UsdtFutures, "BTCUSDT")
	noop := func(RawTopicV2) error { return nil }

	t.Run("deferred until ready", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: false}
		s := NewSubscriptionsV2(c)
		s.subscribe(args2H, noop)
		if len(c.subscribes) != 0 {
			t.Fatalf("expected no sends while not ready, got %d", len(c.subscribes))
		}
		s.subscribeAll()
		if len(c.subscribes) != 1 || len(c.subscribes[0]) != 1 || c.subscribes[0][0] != args2H {
			t.Fatalf("expected one batch with one arg, got %+v", c.subscribes)
		}
	})

	t.Run("immediate when ready", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: true}
		s := NewSubscriptionsV2(c)
		s.subscribe(args2H, noop)
		if len(c.subscribes) != 1 || c.subscribes[0][0] != args2H {
			t.Fatalf("expected immediate send, got %+v", c.subscribes)
		}
	})

	t.Run("subscribe all batches", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: false}
		s := NewSubscriptionsV2(c)
		s.subscribe(args2H, noop)
		s.subscribe(args1D, noop)
		s.subscribeAll()
		if len(c.subscribes) != 1 {
			t.Fatalf("expected a single batch request, got %d", len(c.subscribes))
		}
		if len(c.subscribes[0]) != 2 {
			t.Fatalf("expected 2 args in batch, got %d", len(c.subscribes[0]))
		}
	})

	t.Run("unsubscribe removes handler", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: true}
		s := NewSubscriptionsV2(c)
		s.subscribe(args2H, noop)
		if err := s.processTopic([]byte(candleV2Push("USDT-FUTURES", "candle2H", "BTCUSDT"))); err != nil {
			t.Fatalf("expected topic routed, got error: %v", err)
		}
		s.unsubscribe(args2H)
		if len(c.unsubscribes) != 1 || c.unsubscribes[0][0] != args2H {
			t.Fatalf("expected unsubscribe sent, got %+v", c.unsubscribes)
		}
		// the handler is gone, but the server has not acked yet: pushes still in
		// flight are expected and must not be reported as an unknown channel
		if err := s.processTopic([]byte(candleV2Push("USDT-FUTURES", "candle2H", "BTCUSDT"))); err != nil {
			t.Fatalf("expected in-flight push dropped silently, got error: %v", err)
		}
		s.unsubscribeConfirmed(args2H)
		if err := s.processTopic([]byte(candleV2Push("USDT-FUTURES", "candle2H", "BTCUSDT"))); err == nil {
			t.Fatal("expected not found error after the unsubscribe ack")
		}
	})

	t.Run("route case-insensitive instType and instId", func(t *testing.T) {
		// the channel name is matched verbatim: its case carries the timeframe
		// (candle1m is a minute, candle1M a month)
		c := &fakeSubscriptionClientV2{ready: true}
		s := NewSubscriptionsV2(c)
		var got int
		s.subscribe(args2H, func(RawTopicV2) error { got++; return nil })
		if err := s.processTopic([]byte(candleV2Push("usdt-futures", "candle2H", "btcusdt"))); err != nil {
			t.Fatalf("expected case-insensitive routing, got error: %v", err)
		}
		if got != 1 {
			t.Fatalf("expected 1 call, got %d", got)
		}
	})

	t.Run("route by channel", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: true}
		s := NewSubscriptionsV2(c)
		var got2H, got1D int
		s.subscribe(args2H, func(RawTopicV2) error { got2H++; return nil })
		s.subscribe(args1D, func(RawTopicV2) error { got1D++; return nil })
		if err := s.processTopic([]byte(candleV2Push("USDT-FUTURES", "candle1Dutc", "BTCUSDT"))); err != nil {
			t.Fatalf("expected topic routed, got error: %v", err)
		}
		if got2H != 0 || got1D != 1 {
			t.Fatalf("expected only the 1Dutc handler called, got 2H=%d 1Dutc=%d", got2H, got1D)
		}
	})

	t.Run("channel mismatch", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: true}
		s := NewSubscriptionsV2(c)
		s.subscribe(args2H, noop)
		if err := s.processTopic([]byte(candleV2Push("USDT-FUTURES", "candle4H", "BTCUSDT"))); err == nil {
			t.Fatal("expected not found error for an unregistered channel")
		}
	})

	t.Run("instId mismatch", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: true}
		s := NewSubscriptionsV2(c)
		s.subscribe(args2H, noop)
		if err := s.processTopic([]byte(candleV2Push("USDT-FUTURES", "candle2H", "ETHUSDT"))); err == nil {
			t.Fatal("expected not found error for an unregistered symbol")
		}
	})

	t.Run("handler error propagates", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: true}
		s := NewSubscriptionsV2(c)
		wantErr := fmt.Errorf("handler failed")
		s.subscribe(args2H, func(RawTopicV2) error { return wantErr })
		if err := s.processTopic([]byte(candleV2Push("USDT-FUTURES", "candle2H", "BTCUSDT"))); err != wantErr {
			t.Fatalf("expected handler error, got %v", err)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		c := &fakeSubscriptionClientV2{ready: true}
		s := NewSubscriptionsV2(c)
		if err := s.processTopic([]byte("not json")); err == nil {
			t.Fatal("expected error for invalid payload")
		}
	})
}
