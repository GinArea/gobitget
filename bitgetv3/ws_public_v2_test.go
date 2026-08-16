package bitgetv3

import (
	"sync/atomic"
	"testing"
	"time"
)

// Live tests of the legacy v2 public WebSocket (network access required, no credentials)

func wsCandleV2Chan(c *WsPublicV2, category Category, symbol string, interval Interval) (*ExecutorV2[[]WsCandleV2], chan TopicV2[[]WsCandleV2]) {
	ch := make(chan TopicV2[[]WsCandleV2], 16)
	e := c.Candle(category, symbol, interval)
	e.Subscribe(func(v TopicV2[[]WsCandleV2]) {
		select {
		case ch <- v:
		default:
		}
	})
	return e, ch
}

func waitCandleV2(t *testing.T, ch chan TopicV2[[]WsCandleV2], timeout time.Duration) TopicV2[[]WsCandleV2] {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(timeout):
		t.Fatalf("no candle push within %v", timeout)
		return TopicV2[[]WsCandleV2]{}
	}
}

func TestWsCandleV2(t *testing.T) {
	const hourMs = int64(3_600_000)

	tests := []struct {
		name     string
		category Category
		symbol   string
		interval Interval
		// intervalMs - expected candle period for the UTC-grid alignment check
		intervalMs int64
		// maxAge - how old the current candle start may be
		maxAge time.Duration
		// second - expect a follow-up push: pushes come once per second while trades occur,
		// but quiet stretches happen (especially on spot), hence the generous wait
		second bool
	}{
		{
			name:       "spot 1m",
			category:   Spot,
			symbol:     "BTCUSDT",
			interval:   Interval1m,
			intervalMs: 60_000,
			maxAge:     2 * time.Minute,
			second:     true,
		},
		{
			// the timeframe the v3 WS lacks entirely
			name:       "usdt futures 2H",
			category:   UsdtFutures,
			symbol:     "BTCUSDT",
			interval:   Interval2H,
			intervalMs: 2 * hourMs,
			maxAge:     2*time.Hour + time.Minute,
		},
		{
			// exists only on the v2 WS
			name:       "usdt futures 8H",
			category:   UsdtFutures,
			symbol:     "BTCUSDT",
			interval:   Interval8H,
			intervalMs: 8 * hourMs,
			maxAge:     8*time.Hour + time.Minute,
		},
		{
			// UTC-grid daily candles the v3 WS cannot provide (its 1D opens 16:00 UTC)
			name:       "usdt futures 1Dutc",
			category:   UsdtFutures,
			symbol:     "BTCUSDT",
			interval:   Interval1Dutc,
			intervalMs: 24 * hourMs,
			maxAge:     24*time.Hour + time.Minute,
		},
		{
			name:       "spot 12Hutc",
			category:   Spot,
			symbol:     "BTCUSDT",
			interval:   Interval12Hutc,
			intervalMs: 12 * hourMs,
			maxAge:     12*time.Hour + time.Minute,
		},
	}

	c := NewWsPublicV2()
	c.Run()
	defer c.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, ch := wsCandleV2Chan(c, tt.category, tt.symbol, tt.interval)
			defer e.Unsubscribe()
			v := waitCandleV2(t, ch, 30*time.Second)
			if want := wsArgsV2("candle"+string(tt.interval), tt.category, tt.symbol); v.Arg != want {
				t.Fatalf("expected arg %+v, got %+v", want, v.Arg)
			}
			if v.Action != "snapshot" && v.Action != "update" {
				t.Fatalf("expected action snapshot or update, got %s", v.Action)
			}
			if len(v.Data) == 0 {
				t.Fatal("expected non-empty data")
			}
			// the snapshot carries candle history sorted oldest-first: the current candle is last
			d := v.Data[len(v.Data)-1]
			if first := v.Data[0]; first.Start.Value() > d.Start.Value() {
				t.Fatalf("expected candles sorted oldest-first, got first start %v > last start %v",
					first.Start.Value(), d.Start.Value())
			}
			if d.Open.Value() <= 0 || d.High.Value() <= 0 || d.Low.Value() <= 0 || d.Close.Value() <= 0 {
				t.Fatalf("expected positive ohlc, got o=%v h=%v l=%v c=%v",
					d.Open.Value(), d.High.Value(), d.Low.Value(), d.Close.Value())
			}
			if d.High.Value() < d.Low.Value() {
				t.Fatalf("expected high %v >= low %v", d.High.Value(), d.Low.Value())
			}
			if d.High.Value() < d.Open.Value() || d.High.Value() < d.Close.Value() {
				t.Fatalf("expected high %v >= open %v and close %v", d.High.Value(), d.Open.Value(), d.Close.Value())
			}
			if d.Low.Value() > d.Open.Value() || d.Low.Value() > d.Close.Value() {
				t.Fatalf("expected low %v <= open %v and close %v", d.Low.Value(), d.Open.Value(), d.Close.Value())
			}
			if d.QuoteVolume.Value() < 0 || d.UsdtVolume.Value() < 0 {
				t.Fatalf("expected non-negative volumes, got quote=%v usdt=%v",
					d.QuoteVolume.Value(), d.UsdtVolume.Value())
			}
			// all the streamed intervals must sit on the UTC grid - the point of using the v2 WS
			for i, x := range v.Data {
				if x.Start.Value()%tt.intervalMs != 0 {
					t.Fatalf("candle %d: expected start %v UTC-aligned to %v ms", i, x.Start.Value(), tt.intervalMs)
				}
			}
			startAge := time.Since(time.UnixMilli(d.Start.Value()))
			if startAge < -5*time.Second || startAge > tt.maxAge {
				t.Fatalf("expected the current candle start within %v, got age %v", tt.maxAge, startAge)
			}
			age := time.Since(time.UnixMilli(v.Ts.Value()))
			if age < -time.Minute || age > time.Minute {
				t.Fatalf("expected fresh ts, got age %v", age)
			}
			if tt.second {
				waitCandleV2(t, ch, 30*time.Second)
			}
			t.Logf("snapshot of %d candles; last: start %v, o %v, h %v, l %v, c %v, vol %v, ts age %v",
				len(v.Data), time.UnixMilli(d.Start.Value()).UTC().Format("2006-01-02 15:04:05"),
				d.Open.Value(), d.High.Value(), d.Low.Value(), d.Close.Value(), d.Volume.Value(), age)
		})
	}
}

func TestWsCandleV2Error(t *testing.T) {
	errCh := make(chan WsResponseV2, 1)
	c := NewWsPublicV2().WithOnError(func(r WsResponseV2) {
		select {
		case errCh <- r:
		default:
		}
	})
	c.Run()
	defer c.Close()

	// 3H is the one timeframe the v2 WS rejects
	e := c.Candle(UsdtFutures, "BTCUSDT", Interval("3H"))
	e.Subscribe(func(TopicV2[[]WsCandleV2]) {})
	defer e.Unsubscribe()

	select {
	case r := <-errCh:
		if r.Code.Value() == 0 {
			t.Fatalf("expected non-zero error code, got %+v", r)
		}
		if !r.ParamError() {
			t.Fatalf("expected param error (30016), got %v", r.Code.Value())
		}
		t.Logf("error event: code %v, msg %s", r.Code.Value(), r.Msg)
	case <-time.After(15 * time.Second):
		t.Fatal("no error event within 15s for an unsupported timeframe")
	}
}

func TestWsCandleV2Unsubscribe(t *testing.T) {
	c := NewWsPublicV2()
	c.Run()
	defer c.Close()

	var count atomic.Int64
	ch := make(chan struct{}, 1)
	e := c.Candle(UsdtFutures, "BTCUSDT", Interval1m)
	e.Subscribe(func(TopicV2[[]WsCandleV2]) {
		count.Add(1)
		select {
		case ch <- struct{}{}:
		default:
		}
	})
	select {
	case <-ch:
	case <-time.After(30 * time.Second):
		t.Fatal("no candle push within 30s")
	}

	e.Unsubscribe()
	// grace period: the server may push a few more messages until it processes the request
	time.Sleep(time.Second)
	before := count.Load()
	time.Sleep(3 * time.Second)
	after := count.Load()
	if after != before {
		t.Fatalf("expected no pushes after unsubscribe, got %d", after-before)
	}
	t.Logf("pushes before unsubscribe: %d", before)
}

func TestWsCandleV2Reconnect(t *testing.T) {
	disconnected := make(chan struct{}, 1)
	c := NewWsPublicV2().WithOnDisconnected(func() {
		select {
		case disconnected <- struct{}{}:
		default:
		}
	})
	c.Run()
	defer c.Close()

	e, ch := wsCandleV2Chan(c, UsdtFutures, "BTCUSDT", Interval1m)
	defer e.Unsubscribe()
	waitCandleV2(t, ch, 30*time.Second)

	c.Reconnect()
	select {
	case <-disconnected:
	case <-time.After(10 * time.Second):
		t.Fatal("no disconnect within 10s after Reconnect")
	}
	if c.Ready() {
		t.Fatal("expected not ready right after disconnect")
	}

	// drain pushes received before the disconnect
	for {
		select {
		case <-ch:
			continue
		default:
		}
		break
	}
	// pushes must resume without an explicit resubscribe (subscribeAll on reconnect)
	v := waitCandleV2(t, ch, 30*time.Second)
	if want := wsArgsV2("candle1m", UsdtFutures, "BTCUSDT"); v.Arg != want {
		t.Fatalf("expected arg %+v after reconnect, got %+v", want, v.Arg)
	}
	if !c.Ready() {
		t.Fatal("expected ready after reconnect")
	}
	t.Log("pushes resumed after reconnect")
}

func TestWsV2Keepalive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping keepalive test in short mode")
	}
	c := NewWsPublicV2()
	c.Run()
	defer c.Close()

	var lastPush atomic.Int64
	e := c.Candle(UsdtFutures, "BTCUSDT", Interval1m)
	e.Subscribe(func(TopicV2[[]WsCandleV2]) {
		lastPush.Store(time.Now().UnixMilli())
	})
	defer e.Unsubscribe()

	// like v3, the server drops a connection it got no ping from:
	// a broken keepalive shows up as a push gap within a 2.5 minute window
	deadline := time.Now().Add(150 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(10 * time.Second)
		if !c.Connected() {
			t.Fatal("connection lost during keepalive window")
		}
		last := lastPush.Load()
		if last == 0 {
			t.Fatal("no pushes received")
		}
		if gap := time.Since(time.UnixMilli(last)); gap > 30*time.Second {
			t.Fatalf("push gap %v exceeds 30s", gap)
		}
	}
	t.Log("connection alive for 2.5 minutes, pushes uninterrupted")
}
