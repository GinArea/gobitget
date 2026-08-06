package bitgetv3

import (
	"testing"
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
