package bitgetv3

import (
	"testing"
	"time"
)

func TestGetInstruments(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		symbol   string
	}{
		{
			name:     "usdt futures",
			category: UsdtFutures,
		},
		{
			name:     "usdt futures BTCUSDT",
			category: UsdtFutures,
			symbol:   "BTCUSDT",
		},
		{
			name:     "coin futures",
			category: CoinFutures,
		},
	}

	c := NewClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := c.GetInstruments(GetInstruments{Category: tt.category, Symbol: tt.symbol})
			if !r.Ok() {
				t.Fatalf("GetInstruments failed: %v", r.Error)
			}
			if len(r.Data) == 0 {
				t.Fatal("expected non-empty instruments list")
			}
			if tt.symbol != "" {
				if len(r.Data) != 1 {
					t.Fatalf("expected 1 instrument, got %d", len(r.Data))
				}
				if r.Data[0].Symbol != tt.symbol {
					t.Fatalf("expected symbol %s, got %s", tt.symbol, r.Data[0].Symbol)
				}
			}
			for _, v := range r.Data {
				if v.Category != tt.category {
					t.Fatalf("expected category %s, got %s", tt.category, v.Category)
				}
			}
			t.Logf("total instruments: %d", len(r.Data))
		})
	}
}

func TestGetInstrumentsSpot(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
	}{
		{
			name: "all",
		},
		{
			name:   "BTCUSDT",
			symbol: "BTCUSDT",
		},
	}

	c := NewClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := c.GetInstrumentsSpot(GetInstruments{Symbol: tt.symbol})
			if !r.Ok() {
				t.Fatalf("GetInstrumentsSpot failed: %v", r.Error)
			}
			if len(r.Data) == 0 {
				t.Fatal("expected non-empty instruments list")
			}
			if tt.symbol != "" {
				if len(r.Data) != 1 {
					t.Fatalf("expected 1 instrument, got %d", len(r.Data))
				}
				if r.Data[0].Symbol != tt.symbol {
					t.Fatalf("expected symbol %s, got %s", tt.symbol, r.Data[0].Symbol)
				}
			}
			t.Logf("total instruments: %d", len(r.Data))
		})
	}
}

func TestGetInstrumentsMargin(t *testing.T) {
	c := NewClient()
	r := c.GetInstrumentsMargin(GetInstruments{})
	if !r.Ok() {
		t.Fatalf("GetInstrumentsMargin failed: %v", r.Error)
	}
	if len(r.Data) == 0 {
		t.Fatal("expected non-empty instruments list")
	}
	t.Logf("total instruments: %d", len(r.Data))
}

func TestGetTickers(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		symbol   string
	}{
		{
			name:     "usdt futures",
			category: UsdtFutures,
		},
		{
			name:     "usdt futures BTCUSDT",
			category: UsdtFutures,
			symbol:   "BTCUSDT",
		},
		{
			name:     "coin futures",
			category: CoinFutures,
		},
	}

	c := NewClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := c.GetTickers(GetTickers{Category: tt.category, Symbol: tt.symbol})
			if !r.Ok() {
				t.Fatalf("GetTickers failed: %v", r.Error)
			}
			if len(r.Data) == 0 {
				t.Fatal("expected non-empty tickers list")
			}
			if tt.symbol != "" {
				if len(r.Data) != 1 {
					t.Fatalf("expected 1 ticker, got %d", len(r.Data))
				}
				if r.Data[0].Symbol != tt.symbol {
					t.Fatalf("expected symbol %s, got %s", tt.symbol, r.Data[0].Symbol)
				}
				if r.Data[0].LastPrice.Value() <= 0 {
					t.Fatalf("expected positive last price, got %v", r.Data[0].LastPrice)
				}
				if r.Data[0].Ts.Value() <= 0 {
					t.Fatalf("expected positive ts, got %v", r.Data[0].Ts)
				}
			}
			for _, v := range r.Data {
				if v.Category != tt.category {
					t.Fatalf("expected category %s, got %s", tt.category, v.Category)
				}
			}
			t.Logf("total tickers: %d", len(r.Data))
		})
	}
}

func TestGetTickersSpot(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
	}{
		{
			name: "all",
		},
		{
			name:   "BTCUSDT",
			symbol: "BTCUSDT",
		},
	}

	c := NewClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := c.GetTickersSpot(GetTickers{Symbol: tt.symbol})
			if !r.Ok() {
				t.Fatalf("GetTickersSpot failed: %v", r.Error)
			}
			if len(r.Data) == 0 {
				t.Fatal("expected non-empty tickers list")
			}
			if tt.symbol != "" {
				if len(r.Data) != 1 {
					t.Fatalf("expected 1 ticker, got %d", len(r.Data))
				}
				if r.Data[0].Symbol != tt.symbol {
					t.Fatalf("expected symbol %s, got %s", tt.symbol, r.Data[0].Symbol)
				}
				if r.Data[0].LastPrice.Value() <= 0 {
					t.Fatalf("expected positive last price, got %v", r.Data[0].LastPrice)
				}
				if r.Data[0].Ts.Value() <= 0 {
					t.Fatalf("expected positive ts, got %v", r.Data[0].Ts)
				}
			}
			t.Logf("total tickers: %d", len(r.Data))
		})
	}
}

func TestUnmarshalCandle(t *testing.T) {
	tests := []struct {
		name    string
		row     []string
		want    Candle
		wantErr bool
	}{
		{
			// verbatim row from the docs response example
			name: "docs example",
			row:  []string{"1687708800000", "27176.93", "27177.43", "27166.93", "27177.43", "2990.08", "81246917.3294"},
			want: Candle{
				Ts:       1687708800000,
				Open:     27176.93,
				High:     27177.43,
				Low:      27166.93,
				Close:    27177.43,
				Volume:   2990.08,
				Turnover: 81246917.3294,
			},
		},
		{
			name:    "short row",
			row:     []string{"1687708800000", "27176.93"},
			wantErr: true,
		},
		{
			name:    "long row",
			row:     []string{"1687708800000", "1", "2", "3", "4", "5", "6", "7"},
			wantErr: true,
		},
		{
			name:    "invalid ts",
			row:     []string{"not-a-number", "1", "2", "3", "4", "5", "6"},
			wantErr: true,
		},
		{
			name:    "invalid price",
			row:     []string{"1687708800000", "1", "2", "oops", "4", "5", "6"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := unmarshalCandle(tt.row)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", r)
				}
				return
			}
			if err != nil {
				t.Fatalf("unmarshalCandle failed: %v", err)
			}
			if r != tt.want {
				t.Fatalf("expected %+v, got %+v", tt.want, r)
			}
		})
	}
}

func TestGetCandles(t *testing.T) {
	const (
		minuteMs = int64(60_000)
		dayMs    = int64(86_400_000)
	)
	now := time.Now().UnixMilli()

	tests := []struct {
		name       string
		req        GetCandles
		intervalMs int64
		// alignOffsetMs - candle start alignment offset: daily candles open at 00:00 UTC+8
		alignOffsetMs int64
		wantLen       int
		fresh         bool
	}{
		{
			name:       "usdt futures 1m",
			req:        GetCandles{Category: UsdtFutures, Symbol: "BTCUSDT", Interval: Interval1m},
			intervalMs: minuteMs,
			fresh:      true,
		},
		{
			name:       "spot 1m",
			req:        GetCandles{Category: Spot, Symbol: "BTCUSDT", Interval: Interval1m},
			intervalMs: minuteMs,
			fresh:      true,
		},
		{
			name:       "limit 10",
			req:        GetCandles{Category: UsdtFutures, Symbol: "BTCUSDT", Interval: Interval1m, Limit: 10},
			intervalMs: minuteMs,
			wantLen:    10,
			fresh:      true,
		},
		{
			name:          "1D",
			req:           GetCandles{Category: UsdtFutures, Symbol: "BTCUSDT", Interval: Interval1D},
			intervalMs:    dayMs,
			alignOffsetMs: 8 * 3_600_000,
		},
		{
			name:       "mark price",
			req:        GetCandles{Category: UsdtFutures, Symbol: "BTCUSDT", Interval: Interval1m, Type: CandleMark, Limit: 10},
			intervalMs: minuteMs,
			wantLen:    10,
			fresh:      true,
		},
		{
			// half-hour window one hour ago: startTime is exclusive, endTime is inclusive,
			// so 30 candle starts fall within (start, end]
			name: "time window",
			req: GetCandles{
				Category:  UsdtFutures,
				Symbol:    "BTCUSDT",
				Interval:  Interval1m,
				StartTime: (now - 90*minuteMs) / minuteMs * minuteMs,
				EndTime:   (now - 60*minuteMs) / minuteMs * minuteMs,
			},
			intervalMs: minuteMs,
			wantLen:    30,
		},
	}

	c := NewClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := c.GetCandles(tt.req)
			if !r.Ok() {
				t.Fatalf("GetCandles failed: %v", r.Error)
			}
			if len(r.Data) == 0 {
				t.Fatal("expected non-empty candles list")
			}
			if tt.wantLen != 0 && len(r.Data) != tt.wantLen {
				t.Fatalf("expected %d candles, got %d", tt.wantLen, len(r.Data))
			}
			for i, v := range r.Data {
				if v.Open <= 0 || v.High <= 0 || v.Low <= 0 || v.Close <= 0 {
					t.Fatalf("candle %d: expected positive ohlc, got o=%v h=%v l=%v c=%v", i, v.Open, v.High, v.Low, v.Close)
				}
				if v.High < v.Low {
					t.Fatalf("candle %d: expected high %v >= low %v", i, v.High, v.Low)
				}
				if v.High < v.Open || v.High < v.Close {
					t.Fatalf("candle %d: expected high %v >= open %v and close %v", i, v.High, v.Open, v.Close)
				}
				if v.Low > v.Open || v.Low > v.Close {
					t.Fatalf("candle %d: expected low %v <= open %v and close %v", i, v.Low, v.Open, v.Close)
				}
				if (v.Ts+tt.alignOffsetMs)%tt.intervalMs != 0 {
					t.Fatalf("candle %d: expected ts %v aligned to interval %v ms", i, v.Ts, tt.intervalMs)
				}
				if i > 0 && v.Ts <= r.Data[i-1].Ts {
					t.Fatalf("candle %d: expected ascending ts, got %v after %v", i, v.Ts, r.Data[i-1].Ts)
				}
				if tt.req.StartTime != 0 && (v.Ts <= tt.req.StartTime || v.Ts > tt.req.EndTime) {
					t.Fatalf("candle %d: expected ts %v within (%v, %v]", i, v.Ts, tt.req.StartTime, tt.req.EndTime)
				}
			}
			if tt.fresh {
				last := r.Data[len(r.Data)-1]
				if age := time.Since(time.UnixMilli(last.Ts)); age > 2*time.Minute {
					t.Fatalf("expected fresh last candle, got age %v", age)
				}
			}
			last := r.Data[len(r.Data)-1]
			t.Logf("candles: %d, last: start %v, o %v, h %v, l %v, c %v, vol %v",
				len(r.Data), time.UnixMilli(last.Ts).UTC().Format("2006-01-02 15:04:05"),
				last.Open, last.High, last.Low, last.Close, last.Volume)
		})
	}
}
