package bitgetv3

import (
	"fmt"
	"strconv"

	"github.com/msw-x/moon/ujson"
)

// GetInstruments - request for GET /api/v3/market/instruments (UTA)
// https://www.bitget.com/api-doc/uta/public/Instruments
//
//	category Required string Product type: SPOT, MARGIN, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
//	symbol            string Symbol name, e.g. BTCUSDT
type GetInstruments struct {
	// Category - Product type: SPOT, MARGIN, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
	Category Category
	// Symbol - Symbol name, e.g. BTCUSDT
	Symbol string `url:",omitempty"`
}

func getInstruments[T any](o GetInstruments, c *Client) Response[[]T] {
	return GetPub(c.market(), "instruments", o, forward[[]T])
}

func (o GetInstruments) Do(c *Client) Response[[]Instrument] {
	return getInstruments[Instrument](o, c)
}

func (o GetInstruments) DoSpot(c *Client) Response[[]InstrumentSpot] {
	o.Category = Spot
	return getInstruments[InstrumentSpot](o, c)
}

func (o GetInstruments) DoMargin(c *Client) Response[[]InstrumentMargin] {
	o.Category = Margin
	return getInstruments[InstrumentMargin](o, c)
}

func (o *Client) GetInstruments(v GetInstruments) Response[[]Instrument] {
	return v.Do(o)
}

func (o *Client) GetInstrumentsSpot(v GetInstruments) Response[[]InstrumentSpot] {
	return v.DoSpot(o)
}

func (o *Client) GetInstrumentsMargin(v GetInstruments) Response[[]InstrumentMargin] {
	return v.DoMargin(o)
}

// Instrument - item in GET /api/v3/market/instruments response (Futures categories)
// https://www.bitget.com/api-doc/uta/public/Instruments
type Instrument struct {
	// Symbol - Symbol name
	Symbol string
	// Category - Product type: USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
	Category Category
	// BaseCoin - Base coin, e.g. BTC in BTCUSDT
	BaseCoin string
	// QuoteCoin - Quote coin, e.g. USDT in BTCUSDT
	QuoteCoin string
	// IsRwa - Is this an RWA Symbol: YES / NO
	IsRwa string
	// BuyLimitPriceRatio - The ratio of the buy limit price to the market price
	BuyLimitPriceRatio ujson.Float64
	// SellLimitPriceRatio - The ratio of the sell limit price to the market price
	SellLimitPriceRatio ujson.Float64
	// FeeRateUpRatio - The percentage by which the actual fee is increased relative to the base fee
	FeeRateUpRatio ujson.Float64
	// MakerFeeRate - Maker fee rate, in decimal form, e.g. 0.0002 represents 0.02%
	MakerFeeRate ujson.Float64
	// TakerFeeRate - Taker fee rate, in decimal form, e.g. 0.0002 represents 0.02%
	TakerFeeRate ujson.Float64
	// OpenCostUpRatio - The percentage by which the cost of opening a position is increased relative to the base cost
	OpenCostUpRatio ujson.Float64
	// MinOrderQty - Minimum order quantity in terms of the base coin
	MinOrderQty ujson.Float64
	// MaxOrderQty - Maximum order quantity for a single limit order in terms of the base coin, 0 indicates no limit
	MaxOrderQty ujson.Float64
	// PricePrecision - The number of decimal places allowed for the price
	PricePrecision ujson.Int64
	// QuantityPrecision - The number of decimal places allowed for the quantity
	QuantityPrecision ujson.Int64
	// QuotePrecision - The number of decimal places allowed for the price of the quote coin
	QuotePrecision ujson.Int64
	// PriceMultiplier - The order price must be a multiple of it, used along with PricePrecision
	PriceMultiplier ujson.Float64
	// QuantityMultiplier - The order quantity must be a multiple of it, used along with QuantityPrecision
	QuantityMultiplier ujson.Float64
	// Type - Futures type: perpetual, delivery
	Type string
	// MinOrderAmount - Minimum order amount in terms of the quote coin
	MinOrderAmount ujson.Float64
	// MaxSymbolOrderNum - Maximum order number in terms of the trading pair
	MaxSymbolOrderNum ujson.Int64
	// MaxProductOrderNum - Maximum order number in terms of the product line
	MaxProductOrderNum ujson.Int64
	// MaxPositionNum - Maximum position number in terms of the trading pair
	MaxPositionNum ujson.Int64
	// Status - Trading pair status: listed, online, limit_open, limit_close, offline, restrictedAPI
	Status string
	// OffTime - Trading halt time, "" if not configured
	OffTime ujson.Int64
	// LimitOpenTime - Restricted open time, "" if not configured
	LimitOpenTime ujson.Int64
	// DeliveryTime - Delivery time, available only for deliveries
	DeliveryTime ujson.Int64
	// DeliveryStartTime - Delivery start time, available only for deliveries
	DeliveryStartTime ujson.Int64
	// DeliveryPeriod - Delivery period: this_quarter, next_quarter, available only for deliveries
	DeliveryPeriod string
	// LaunchTime - Launch time, unix millisecond timestamp
	LaunchTime ujson.Int64
	// FundInterval - Funding interval in hours: 1, 8
	FundInterval ujson.Int64
	// MinLeverage - Minimum leverage
	MinLeverage ujson.Int64
	// MaxLeverage - Maximum leverage
	MaxLeverage ujson.Int64
	// MaintainTime - Maintenance time, "" if not configured
	MaintainTime string
	// SymbolType - Symbol type: crypto, metal, stock, commodity
	SymbolType string
	// MaxMarketOrderQty - Maximum order quantity for a single market order in terms of the base coin
	MaxMarketOrderQty ujson.Float64
}

// InstrumentSpot - item in GET /api/v3/market/instruments response (SPOT category)
// https://www.bitget.com/api-doc/uta/public/Instruments
type InstrumentSpot struct {
	// Symbol - Symbol name
	Symbol string
	// Category - Product type: SPOT
	Category Category
	// BaseCoin - Base coin, e.g. BTC in BTCUSDT
	BaseCoin string
	// QuoteCoin - Quote coin, e.g. USDT in BTCUSDT
	QuoteCoin string
	// IsRwa - Is this an RWA Symbol: YES / NO
	IsRwa string
	// IsReality - Reality identifier: yes (reality stock token), no
	IsReality string
	// BuyLimitPriceRatio - The ratio of the buy limit price to the market price
	BuyLimitPriceRatio ujson.Float64
	// SellLimitPriceRatio - The ratio of the sell limit price to the market price
	SellLimitPriceRatio ujson.Float64
	// MinOrderQty - Minimum order quantity in terms of the base coin
	MinOrderQty ujson.Float64
	// MaxOrderQty - Maximum order quantity for a single limit order in terms of the base coin, 0 indicates no limit
	MaxOrderQty ujson.Float64
	// PricePrecision - The number of decimal places allowed for the price
	PricePrecision ujson.Int64
	// QuantityPrecision - The number of decimal places allowed for the quantity
	QuantityPrecision ujson.Int64
	// QuotePrecision - The number of decimal places allowed for the price of the quote coin
	QuotePrecision ujson.Int64
	// MinOrderAmount - Minimum order amount in terms of the quote coin
	MinOrderAmount ujson.Float64
	// MaxSymbolOrderNum - Maximum order number in terms of the trading pair
	MaxSymbolOrderNum ujson.Int64
	// MaxProductOrderNum - Maximum order number in terms of the product line
	MaxProductOrderNum ujson.Int64
	// MaxPositionNum - Maximum position number in terms of the trading pair
	MaxPositionNum ujson.Int64
	// Status - Trading pair status: listed, online, limit_open, limit_close, offline, restrictedAPI
	Status string
	// MaintainTime - Maintenance time, "" if not configured
	MaintainTime string
	// AreaSymbol - Area symbol: YES / NO, returned only for pairs where the value is YES
	AreaSymbol string
	// SymbolType - Symbol type: crypto, metal, stock, commodity
	SymbolType string
	// LaunchTime - Launch time, unix millisecond timestamp
	LaunchTime ujson.Int64
}

// InstrumentMargin - item in GET /api/v3/market/instruments response (MARGIN category)
// https://www.bitget.com/api-doc/uta/public/Instruments
type InstrumentMargin struct {
	// Symbol - Symbol name
	Symbol string
	// Category - Product type: MARGIN
	Category Category
	// BaseCoin - Base coin, e.g. BTC in BTCUSDT
	BaseCoin string
	// QuoteCoin - Quote coin, e.g. USDT in BTCUSDT
	QuoteCoin string
	// BuyLimitPriceRatio - The ratio of the buy limit price to the market price
	BuyLimitPriceRatio ujson.Float64
	// SellLimitPriceRatio - The ratio of the sell limit price to the market price
	SellLimitPriceRatio ujson.Float64
	// MinOrderQty - Minimum order quantity in terms of the base coin
	MinOrderQty ujson.Float64
	// MaxOrderQty - Maximum order quantity for a single limit order in terms of the base coin, 0 indicates no limit
	MaxOrderQty ujson.Float64
	// PricePrecision - The number of decimal places allowed for the price
	PricePrecision ujson.Int64
	// QuantityPrecision - The number of decimal places allowed for the quantity
	QuantityPrecision ujson.Int64
	// QuotePrecision - The number of decimal places allowed for the price of the quote coin
	QuotePrecision ujson.Int64
	// MinOrderAmount - Minimum order amount in terms of the quote coin
	MinOrderAmount ujson.Float64
	// MaxSymbolOrderNum - Maximum order number in terms of the trading pair
	MaxSymbolOrderNum ujson.Int64
	// MaxProductOrderNum - Maximum order number in terms of the product line
	MaxProductOrderNum ujson.Int64
	// MaxPositionNum - Maximum position number in terms of the trading pair
	MaxPositionNum ujson.Int64
	// Status - Trading pair status: listed, online, limit_open, limit_close, offline, restrictedAPI
	Status string
	// MaintainTime - Maintenance time, "" if not configured
	MaintainTime string
	// IsIsolatedBaseBorrowable - Base coin borrowable status: YES / NO
	IsIsolatedBaseBorrowable string
	// IsIsolatedQuotedBorrowable - Quote coin borrowable status: YES / NO
	IsIsolatedQuotedBorrowable string
	// WarningRiskRatio - Warning risk ratio
	WarningRiskRatio ujson.Float64
	// LiquidationRiskRatio - Liquidation risk ratio
	LiquidationRiskRatio ujson.Float64
	// MaxCrossedLeverage - Maximum leverage for cross margin
	MaxCrossedLeverage ujson.Int64
	// MaxIsolatedLeverage - Maximum leverage for isolated margin
	MaxIsolatedLeverage ujson.Int64
	// UserMinBorrow - Minimum borrowable amount
	UserMinBorrow ujson.Float64
	// AreaSymbol - Area symbol: yes / no
	AreaSymbol string
	// MaxLeverage - Maximum leverage
	MaxLeverage ujson.Int64
	// SymbolType - Symbol type: crypto, metal, stock, commodity
	SymbolType string
	// LaunchTime - Launch time, unix millisecond timestamp, may be null
	LaunchTime ujson.Int64
}

// GetTickers - request for GET /api/v3/market/tickers (UTA)
// https://www.bitget.com/api-doc/uta/public/Tickers
//
//	category Required string Product type: SPOT, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
//	symbol            string Symbol name, e.g. BTCUSDT
type GetTickers struct {
	// Category - Product type: SPOT, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES (MARGIN is not supported)
	Category Category
	// Symbol - Symbol name, e.g. BTCUSDT
	Symbol string `url:",omitempty"`
}

func getTickers[T any](o GetTickers, c *Client) Response[[]T] {
	return GetPub(c.market(), "tickers", o, forward[[]T])
}

func (o GetTickers) Do(c *Client) Response[[]Ticker] {
	return getTickers[Ticker](o, c)
}

func (o GetTickers) DoSpot(c *Client) Response[[]TickerSpot] {
	o.Category = Spot
	return getTickers[TickerSpot](o, c)
}

func (o *Client) GetTickers(v GetTickers) Response[[]Ticker] {
	return v.Do(o)
}

func (o *Client) GetTickersSpot(v GetTickers) Response[[]TickerSpot] {
	return v.DoSpot(o)
}

// Ticker - item in GET /api/v3/market/tickers response (Futures categories)
// https://www.bitget.com/api-doc/uta/public/Tickers
type Ticker struct {
	// Symbol - Symbol name
	Symbol string
	// Category - Product type: USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
	Category Category
	// LastPrice - Latest price
	LastPrice ujson.Float64
	// OpenPrice24h - Market price 24 hours ago
	OpenPrice24h ujson.Float64
	// HighPrice24h - Highest price in the last 24 hours
	HighPrice24h ujson.Float64
	// LowPrice24h - Lowest price in the last 24 hours
	LowPrice24h ujson.Float64
	// Ask1Price - Best ask price
	Ask1Price ujson.Float64
	// Bid1Price - Best bid price
	Bid1Price ujson.Float64
	// Bid1Size - Best bid quantity
	Bid1Size ujson.Float64
	// Ask1Size - Best ask quantity
	Ask1Size ujson.Float64
	// Price24hPcnt - 24-hour price change percentage
	Price24hPcnt ujson.Float64
	// Volume24h - 24-hour volume
	Volume24h ujson.Float64
	// Turnover24h - 24-hour turnover
	Turnover24h ujson.Float64
	// IndexPrice - Index price
	IndexPrice ujson.Float64
	// MarkPrice - Mark price
	MarkPrice ujson.Float64
	// FundingRate - Funding rate
	FundingRate ujson.Float64
	// OpenInterest - Open interest
	OpenInterest ujson.Float64
	// DeliveryStartTime - Delivery start time, available only for deliveries, "" for perpetuals
	DeliveryStartTime ujson.Int64
	// DeliveryTime - Delivery time, available only for deliveries, "" for perpetuals
	DeliveryTime ujson.Int64
	// DeliveryStatus - Delivery status: delivery_config_period, delivery_normal, delivery_before, delivery_period, available only for deliveries
	DeliveryStatus string
	// Ts - The timestamp that the system generated the data, unix millisecond timestamp
	Ts ujson.Int64
}

// TickerSpot - item in GET /api/v3/market/tickers response (SPOT category)
// https://www.bitget.com/api-doc/uta/public/Tickers
type TickerSpot struct {
	// Symbol - Symbol name
	Symbol string
	// Category - Product type: SPOT
	Category Category
	// LastPrice - Latest price
	LastPrice ujson.Float64
	// OpenPrice24h - Market price 24 hours ago
	OpenPrice24h ujson.Float64
	// HighPrice24h - Highest price in the last 24 hours
	HighPrice24h ujson.Float64
	// LowPrice24h - Lowest price in the last 24 hours
	LowPrice24h ujson.Float64
	// Ask1Price - Best ask price
	Ask1Price ujson.Float64
	// Bid1Price - Best bid price
	Bid1Price ujson.Float64
	// Bid1Size - Best bid quantity
	Bid1Size ujson.Float64
	// Ask1Size - Best ask quantity
	Ask1Size ujson.Float64
	// Price24hPcnt - 24-hour price change percentage
	Price24hPcnt ujson.Float64
	// Volume24h - 24-hour volume
	Volume24h ujson.Float64
	// Turnover24h - 24-hour turnover
	Turnover24h ujson.Float64
	// PlatformTurnover24h - 24-hour platform turnover, only available for rtoken
	PlatformTurnover24h ujson.Float64
	// Ts - The timestamp that the system generated the data, unix millisecond timestamp
	Ts ujson.Int64
}

// GetCandles - request for GET /api/v3/market/candles (UTA)
// https://www.bitget.com/api-doc/uta/public/Get-Candle-Data
//
//	category  Required string Product type: SPOT, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
//	symbol    Required string Symbol name, e.g. BTCUSDT
//	interval  Required string Granularity: 1m, 3m, 5m, 15m, 30m, 1H, 4H, 6H, 12H, 1D
//	startTime          string Start timestamp, unix ms
//	endTime            string End timestamp, unix ms
//	type               string Candlestick type: market (default), mark, index, premium
//	limit              string Limit per page
//
// rtoken symbols support only the market type (mark/index/premium silently fall back to market)
// and only the 1m, 5m, 15m, 1H, 4H, 1D intervals (others return a parameter error)
//
// Undocumented intervals also work (inherited from the v2 API, verified live for all categories):
// 2H, 3D, 1W, 1M and 6Hutc, 12Hutc, 1Dutc, 3Dutc, 1Wutc, 1Mutc.
// Native 6H/12H/1D/3D/1W/1M candles open on the UTC+8 grid (weeks on Monday, months on the 1st),
// the utc-suffixed variants - on the UTC grid; intervals up to 4H are UTC-aligned as is.
// 3H, 8H, 1Hutc and 4Hutc do not exist (parameter error 40020).
// See BITGET_API.md for the full interval/alignment reference table
//
// History depth is shallow and varies by interval: 1m rows exist only ~31 days back,
// older windows return empty data without an error (verified live).
// For deeper history use GetHistoryCandles
type GetCandles struct {
	// Category - Product type: SPOT, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES (MARGIN is not supported)
	Category Category
	// Symbol - Symbol name, e.g. BTCUSDT
	Symbol string
	// Interval - Granularity: 1m, 3m, 5m, 15m, 30m, 1H, 4H, 6H, 12H, 1D
	// (+ undocumented 2H, 3D, 1W, 1M and 6Hutc, 12Hutc, 1Dutc, 3Dutc, 1Wutc, 1Mutc - see the note above)
	Interval Interval
	// StartTime - Start timestamp, unix ms, exclusive: candles with ts > startTime (verified live)
	StartTime int64 `url:",omitempty"`
	// EndTime - End timestamp, unix ms, inclusive: candles with ts <= endTime (verified live)
	EndTime int64 `url:",omitempty"`
	// Type - Candlestick type: market (default), mark, index, premium
	Type CandleType `url:",omitempty"`
	// Limit - Limit per page: default 100, maximum 1000, above 1000 -> parameter error 40020
	// (verified live; the docs table swaps default and maximum)
	Limit int `url:",omitempty"`
}

func (o GetCandles) Do(c *Client) Response[[]Candle] {
	return GetPub(c.market(), "candles", o, func(l [][]string) ([]Candle, error) {
		return transformList(l, unmarshalCandle)
	})
}

func (o *Client) GetCandles(v GetCandles) Response[[]Candle] {
	return v.Do(o)
}

// Candle - item in GET /api/v3/market/candles and history-candles responses:
// [ts, open, high, low, close, volume, turnover]
// Every element arrives as a JSON string; the row shape is the same for all categories
// Rows are sorted oldest-first (verified live)
// https://www.bitget.com/api-doc/uta/public/Get-Candle-Data
type Candle struct {
	// Ts - Candle start timestamp, unix ms; daily candles open at 00:00 UTC+8 (16:00 UTC, verified live)
	Ts int64
	// Open - Open price
	Open float64
	// High - Highest price
	High float64
	// Low - Lowest price
	Low float64
	// Close - Close price
	Close float64
	// Volume - Trade volume, base coin
	Volume float64
	// Turnover - Turnover, quote coin
	Turnover float64
}

// unmarshalCandle - parse one candle row of 7 strings
func unmarshalCandle(s []string) (r Candle, err error) {
	if len(s) != 7 {
		err = fmt.Errorf("candle row length is %d, expected 7", len(s))
		return
	}
	r.Ts, err = strconv.ParseInt(s[0], 10, 64)
	if err != nil {
		err = fmt.Errorf("candle ts: %w", err)
		return
	}
	fields := []struct {
		name string
		p    *float64
		s    string
	}{
		{"open", &r.Open, s[1]},
		{"high", &r.High, s[2]},
		{"low", &r.Low, s[3]},
		{"close", &r.Close, s[4]},
		{"volume", &r.Volume, s[5]},
		{"turnover", &r.Turnover, s[6]},
	}
	for _, f := range fields {
		*f.p, err = strconv.ParseFloat(f.s, 64)
		if err != nil {
			err = fmt.Errorf("candle %s: %w", f.name, err)
			return
		}
	}
	return
}

// GetHistoryCandles - request for GET /api/v3/market/history-candles (UTA)
// https://www.bitget.com/api-doc/uta/public/Get-History-Candle-Data
//
//	category  Required string Product type: SPOT, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
//	symbol    Required string Symbol name, e.g. BTCUSDT
//	interval  Required string Granularity: 1m, 3m, 5m, 15m, 30m, 1H, 4H, 6H, 12H, 1D
//	startTime          string Start timestamp, unix ms
//	endTime            string End timestamp, unix ms
//	type               string Candlestick type: market (default), mark, index, premium
//	limit              string Limit per page
//
// Unlike candles, whose history depth is shallow (1m rows exist only ~31 days back),
// history-candles serves deep history: BTCUSDT 1m rows verified live back to 2022-08
// for both SPOT and USDT-FUTURES. The response row shape and ordering are identical
// to candles.
//
// The time window works differently from GetCandles: startTime is inclusive and
// endTime is exclusive - [startTime, endTime) - and the window must not exceed
// 90 days, otherwise error code 00001 (not the usual parameter error 40020).
// Without startTime the same 90-day window applies IMPLICITLY before endTime
// (before "now" if endTime is also omitted): the response holds
// min(limit, 90 days / interval) rows, silently, with no 00001 error - for
// intervals longer than 21.6h a page is always shorter than limit (1Dutc:
// exactly 90 rows), so only an empty page marks the end of history (verified live).
//
// Only closed candles are returned: the current unclosed bar is never included
// (unlike candles, whose last row is the live unclosed bar; verified live).
//
// Undocumented intervals work the same as in GetCandles (2H, 3D, 1W, 1M and the
// utc-suffixed variants verified live; 8H is a parameter error) - see the GetCandles
// note and BITGET_API.md for the interval/alignment reference
type GetHistoryCandles struct {
	// Category - Product type: SPOT, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES (MARGIN is not supported)
	Category Category
	// Symbol - Symbol name, e.g. BTCUSDT
	Symbol string
	// Interval - Granularity: 1m, 3m, 5m, 15m, 30m, 1H, 4H, 6H, 12H, 1D
	// (+ the undocumented intervals of GetCandles - see the note above)
	Interval Interval
	// StartTime - Start timestamp, unix ms, inclusive: candles with ts >= startTime
	// (opposite of GetCandles; verified live)
	StartTime int64 `url:",omitempty"`
	// EndTime - End timestamp, unix ms, exclusive: candles with ts < endTime
	// (opposite of GetCandles; verified live); endTime-startTime must not exceed 90 days
	EndTime int64 `url:",omitempty"`
	// Type - Candlestick type: market (default), mark, index, premium
	Type CandleType `url:",omitempty"`
	// Limit - Limit per page: default 100, maximum 100 (unlike candles), above 100 -> parameter error 40020
	Limit int `url:",omitempty"`
}

func (o GetHistoryCandles) Do(c *Client) Response[[]Candle] {
	return GetPub(c.market(), "history-candles", o, func(l [][]string) ([]Candle, error) {
		return transformList(l, unmarshalCandle)
	})
}

func (o *Client) GetHistoryCandles(v GetHistoryCandles) Response[[]Candle] {
	return v.Do(o)
}
