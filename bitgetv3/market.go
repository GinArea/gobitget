package bitgetv3

import (
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
