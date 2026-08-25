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

// Interval - Candlestick interval; case matters: minutes are lowercase, hours/days are uppercase.
// Used by REST GetCandles and by the v2 WS candle channel (WsPublicV2.Candle)
// https://www.bitget.com/api-doc/uta/public/Get-Candle-Data
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
	// Interval6H - 6 hours, opens on the UTC+8 grid (04:00/10:00/16:00/22:00 UTC); use Interval6Hutc for UTC
	Interval6H Interval = "6H"
	// Interval12H - 12 hours, opens on the UTC+8 grid (04:00/16:00 UTC); use Interval12Hutc for UTC
	Interval12H Interval = "12H"
	// Interval1D - 1 day, opens 00:00 UTC+8 (16:00 UTC); use Interval1Dutc for UTC
	Interval1D Interval = "1D"

	// The values below are UNDOCUMENTED in the v3 docs (inherited from the v2 API).
	// Verified live on GET /api/v3/market/candles for SPOT, USDT-FUTURES and COIN-FUTURES
	// (grids are identical across categories): native 6H/12H/1D/3D/1W/1M candles open on the UTC+8 grid
	// (weeks on Monday, months on the 1st), the utc-suffixed variants - on the UTC grid.
	// Intervals up to 4H need no utc variant: their period divides the 8-hour offset, so both grids coincide.
	// 3H, 8H, 1Hutc and 4Hutc do not exist on REST (parameter error 40020; 8H exists on the legacy v2 WS).
	// The v3 WS kline channel accepts just the 10 documented values and rejects everything else with
	// error 30001; the legacy v2 WS supports all of them, which is why candle streaming
	// in this library goes through WsPublicV2. See BITGET_API.md for the full reference table.

	// Interval2H - 2 hours (UTC-aligned)
	Interval2H Interval = "2H"
	// Interval8H - 8 hours (UTC-aligned), exists ONLY on the legacy v2 WS (rejected by REST and v3 WS)
	Interval8H Interval = "8H"
	// Interval3D - 3 days, opens 00:00 UTC+8
	Interval3D Interval = "3D"
	// Interval1W - 1 week, opens Monday 00:00 UTC+8 (Sunday 16:00 UTC)
	Interval1W Interval = "1W"
	// Interval1M - 1 calendar month, opens on the 1st at 00:00 UTC+8
	Interval1M Interval = "1M"
	// Interval6Hutc - 6 hours, UTC grid
	Interval6Hutc Interval = "6Hutc"
	// Interval12Hutc - 12 hours, UTC grid
	Interval12Hutc Interval = "12Hutc"
	// Interval1Dutc - 1 day, opens 00:00 UTC
	Interval1Dutc Interval = "1Dutc"
	// Interval3Dutc - 3 days, opens 00:00 UTC on the UTC grid
	Interval3Dutc Interval = "3Dutc"
	// Interval1Wutc - 1 week, opens Monday 00:00 UTC
	Interval1Wutc Interval = "1Wutc"
	// Interval1Mutc - 1 calendar month, opens on the 1st at 00:00 UTC
	Interval1Mutc Interval = "1Mutc"
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

// HoldMode - Position holding mode
// https://www.bitget.com/api-doc/uta/account/Change-Position-Mode
type HoldMode string

const (
	// OneWayMode - Holding positions in a single direction, either long or short, but not both at the same time
	OneWayMode HoldMode = "one_way_mode"
	// HedgeMode - Holding both long and short positions simultaneously
	HedgeMode HoldMode = "hedge_mode"
)
