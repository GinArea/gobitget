package bitgetv3

// Category - Product type
// https://www.bitget.com/api-doc/uta/enum#category
type Category string

const (
	// Spot - Spot trading
	Spot Category = "SPOT"
	// Margin - Margin trading
	Margin Category = "MARGIN"
	// UsdtFutures - USDT futures
	UsdtFutures Category = "USDT-FUTURES"
	// CoinFutures - Coin-M futures
	CoinFutures Category = "COIN-FUTURES"
	// UsdcFutures - USDC futures
	UsdcFutures Category = "USDC-FUTURES"
)

// Interval - Candlestick interval; case matters: minutes are lowercase, hours/days are uppercase
// https://www.bitget.com/api-doc/uta/websocket/public/Candlesticks-Channel
type Interval string

const (
	// Interval1m - 1 minute
	Interval1m Interval = "1m"
	// Interval3m - 3 minutes
	Interval3m Interval = "3m"
	// Interval5m - 5 minutes
	Interval5m Interval = "5m"
	// Interval15m - 15 minutes
	Interval15m Interval = "15m"
	// Interval30m - 30 minutes
	Interval30m Interval = "30m"
	// Interval1H - 1 hour
	Interval1H Interval = "1H"
	// Interval4H - 4 hours
	Interval4H Interval = "4H"
	// Interval6H - 6 hours
	Interval6H Interval = "6H"
	// Interval12H - 12 hours
	Interval12H Interval = "12H"
	// Interval1D - 1 day
	Interval1D Interval = "1D"
)

// CandleType - Candlestick price source
// For rtoken symbols only market is supported: mark/index/premium silently fall back to market
// https://www.bitget.com/api-doc/uta/public/Get-Candle-Data
type CandleType string

const (
	// CandleMarket - Market price candles (default)
	CandleMarket CandleType = "market"
	// CandleMark - Mark price candles (futures only; volume/turnover are 0, verified live)
	CandleMark CandleType = "mark"
	// CandleIndex - Index price candles (futures only)
	CandleIndex CandleType = "index"
	// CandlePremium - Premium index price candles (futures only)
	CandlePremium CandleType = "premium"
)
