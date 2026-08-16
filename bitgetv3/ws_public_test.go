package bitgetv3

import (
	"sync/atomic"
	"testing"
	"time"
)

// Live public WebSocket tests (network access required, no credentials)

func wsTickerChan(c *WsPublic, category Category, symbol string) (*Executor[[]WsTicker], chan Topic[[]WsTicker]) {
	ch := make(chan Topic[[]WsTicker], 16)
	e := c.Ticker(category, symbol)
	e.Subscribe(func(v Topic[[]WsTicker]) {
		select {
		case ch <- v:
		default:
		}
	})
	return e, ch
}

func waitTicker(t *testing.T, ch chan Topic[[]WsTicker], timeout time.Duration) Topic[[]WsTicker] {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(timeout):
		t.Fatalf("no ticker push within %v", timeout)
		return Topic[[]WsTicker]{}
	}
}

func TestWsTicker(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		symbol   string
		futures  bool
	}{
		{
			name:     "usdt futures",
			category: UsdtFutures,
			symbol:   "BTCUSDT",
			futures:  true,
		},
		{
			name:     "spot",
			category: Spot,
			symbol:   "BTCUSDT",
		},
	}

	c := NewWsPublic()
	c.Run()
	defer c.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, ch := wsTickerChan(c, tt.category, tt.symbol)
			defer e.Unsubscribe()
			v := waitTicker(t, ch, 30*time.Second)
			if want := wsArgs("ticker", tt.category, tt.symbol); v.Arg != want {
				t.Fatalf("expected arg %+v, got %+v", want, v.Arg)
			}
			if v.Action != "snapshot" {
				t.Fatalf("expected action snapshot, got %s", v.Action)
			}
			if len(v.Data) == 0 {
				t.Fatal("expected non-empty data")
			}
			d := v.Data[0]
			if d.LastPrice.Value() <= 0 {
				t.Fatalf("expected positive lastPrice, got %v", d.LastPrice.Value())
			}
			if d.Bid1Price.Value() <= 0 {
				t.Fatalf("expected positive bid1Price, got %v", d.Bid1Price.Value())
			}
			if d.Ask1Price.Value() <= 0 {
				t.Fatalf("expected positive ask1Price, got %v", d.Ask1Price.Value())
			}
			if d.Ask1Price.Value() < d.Bid1Price.Value() {
				t.Fatalf("expected ask %v >= bid %v", d.Ask1Price.Value(), d.Bid1Price.Value())
			}
			age := time.Since(time.UnixMilli(v.Ts.Value()))
			if age < -time.Minute || age > time.Minute {
				t.Fatalf("expected fresh ts, got age %v", age)
			}
			if tt.futures {
				t.Logf("markPrice: %v, indexPrice: %v, fundingRate: %v, nextFundingTime: %v",
					d.MarkPrice.Value(), d.IndexPrice.Value(), d.FundingRate.Value(), d.NextFundingTime.Value())
			}
			t.Logf("ticker: last %v, bid %v, ask %v, ts age %v",
				d.LastPrice.Value(), d.Bid1Price.Value(), d.Ask1Price.Value(), age)
		})
	}
}

func wsCandleChan(c *WsPublic, category Category, symbol string, interval Interval) (*Executor[[]WsCandle], chan Topic[[]WsCandle]) {
	ch := make(chan Topic[[]WsCandle], 16)
	e := c.Candle(category, symbol, interval)
	e.Subscribe(func(v Topic[[]WsCandle]) {
		select {
		case ch <- v:
		default:
		}
	})
	return e, ch
}

func waitCandle(t *testing.T, ch chan Topic[[]WsCandle], timeout time.Duration) Topic[[]WsCandle] {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(timeout):
		t.Fatalf("no candle push within %v", timeout)
		return Topic[[]WsCandle]{}
	}
}

func TestWsCandle(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		symbol   string
		interval Interval
	}{
		{
			name:     "usdt futures",
			category: UsdtFutures,
			symbol:   "BTCUSDT",
			interval: Interval1m,
		},
		{
			name:     "spot",
			category: Spot,
			symbol:   "BTCUSDT",
			interval: Interval1m,
		},
	}

	c := NewWsPublic()
	c.Run()
	defer c.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, ch := wsCandleChan(c, tt.category, tt.symbol, tt.interval)
			defer e.Unsubscribe()
			v := waitCandle(t, ch, 30*time.Second)
			want := wsArgs("kline", tt.category, tt.symbol)
			want.Interval = string(tt.interval)
			if v.Arg != want {
				t.Fatalf("expected arg %+v, got %+v", want, v.Arg)
			}
			if v.Action != "snapshot" && v.Action != "update" {
				t.Fatalf("expected action snapshot or update, got %s", v.Action)
			}
			if len(v.Data) == 0 {
				t.Fatal("expected non-empty data")
			}
			// the initial snapshot carries recent candle history sorted oldest-first:
			// the current candle is the last item
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
			startAge := time.Since(time.UnixMilli(d.Start.Value()))
			if startAge < -5*time.Second || startAge > 2*time.Minute {
				t.Fatalf("expected start of the current 1m candle, got age %v", startAge)
			}
			age := time.Since(time.UnixMilli(v.Ts.Value()))
			if age < -time.Minute || age > time.Minute {
				t.Fatalf("expected fresh ts, got age %v", age)
			}
			// BTCUSDT trades continuously: candles are pushed once per second
			waitCandle(t, ch, 10*time.Second)
			t.Logf("snapshot of %d candles; last: start %v, o %v, h %v, l %v, c %v, vol %v, ts age %v",
				len(v.Data), time.UnixMilli(d.Start.Value()).UTC().Format("15:04:05"),
				d.Open.Value(), d.High.Value(), d.Low.Value(), d.Close.Value(), d.Volume.Value(), age)
		})
	}
}

func TestWsTickerUnsubscribe(t *testing.T) {
	c := NewWsPublic()
	c.Run()
	defer c.Close()

	var count atomic.Int64
	ch := make(chan struct{}, 1)
	e := c.Ticker(UsdtFutures, "BTCUSDT")
	e.Subscribe(func(Topic[[]WsTicker]) {
		count.Add(1)
		select {
		case ch <- struct{}{}:
		default:
		}
	})
	select {
	case <-ch:
	case <-time.After(30 * time.Second):
		t.Fatal("no ticker push within 30s")
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

func TestWsTickerReconnect(t *testing.T) {
	disconnected := make(chan struct{}, 1)
	c := NewWsPublic().WithOnDisconnected(func() {
		select {
		case disconnected <- struct{}{}:
		default:
		}
	})
	c.Run()
	defer c.Close()

	e, ch := wsTickerChan(c, UsdtFutures, "BTCUSDT")
	defer e.Unsubscribe()
	waitTicker(t, ch, 30*time.Second)

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
	v := waitTicker(t, ch, 30*time.Second)
	if want := wsArgs("ticker", UsdtFutures, "BTCUSDT"); v.Arg != want {
		t.Fatalf("expected arg %+v after reconnect, got %+v", want, v.Arg)
	}
	if !c.Ready() {
		t.Fatal("expected ready after reconnect")
	}
	t.Log("pushes resumed after reconnect")
}

func TestWsTickerError(t *testing.T) {
	errCh := make(chan WsResponse, 1)
	c := NewWsPublic().WithOnError(func(r WsResponse) {
		select {
		case errCh <- r:
		default:
		}
	})
	c.Run()
	defer c.Close()

	e := c.Ticker(UsdtFutures, "NOSUCHSYMBOL")
	e.Subscribe(func(Topic[[]WsTicker]) {})
	defer e.Unsubscribe()

	select {
	case r := <-errCh:
		if r.Code.Value() == 0 {
			t.Fatalf("expected non-zero error code, got %+v", r)
		}
		if !r.TopicNotExists() {
			t.Fatalf("expected topic not exists error (30001), got %v", r.Code.Value())
		}
		t.Logf("error event: code %v, msg %s", r.Code.Value(), r.Msg)
	case <-time.After(15 * time.Second):
		t.Fatal("no error event within 15s for unknown symbol")
	}
}

func TestWsKeepalive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping keepalive test in short mode")
	}
	c := NewWsPublic()
	c.Run()
	defer c.Close()

	var lastPush atomic.Int64
	e := c.Ticker(UsdtFutures, "BTCUSDT")
	e.Subscribe(func(Topic[[]WsTicker]) {
		lastPush.Store(time.Now().UnixMilli())
	})
	defer e.Unsubscribe()

	// the server drops a connection it got no ping from for 2 minutes:
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
