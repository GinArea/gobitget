package bitgetv3

import (
	"github.com/msw-x/moon/ujson"
)

// GetPositions - request for GET /api/v3/position/current-position (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Position
// Rate limit: 20/sec/UID, permission: UTA trade (read)
type GetPositions struct {
	// Category - Product type: USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
	Category Category
	// Symbol - Symbol name, e.g. BTCUSDT; if omitted, all positions in the category are returned
	Symbol string `url:",omitempty"`
	// PosSide - Position side: long, short; if provided, only positions of this side are returned
	PosSide string `url:",omitempty"`
}

// positionList - wire object of GET /api/v3/position/current-position response data
type positionList struct {
	// List - Position list (null when there are no open positions)
	List []Position
}

func (o GetPositions) Do(c *Client) Response[[]Position] {
	return Get(c, "position/current-position", o, func(v positionList) ([]Position, error) {
		return v.List, nil
	})
}

func (o *Client) GetPositions(v GetPositions) Response[[]Position] {
	return v.Do(o)
}

// Position - item of list in GET /api/v3/position/current-position response (UTA)
// https://www.bitget.com/api-doc/uta/trade/Get-Position
type Position struct {
	// Category - Product type: USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
	Category Category
	// Symbol - Symbol name
	Symbol string
	// MarginCoin - Margin coin
	MarginCoin string
	// HoldMode - Holding mode: one_way_mode, hedge_mode
	HoldMode string
	// PosSide - Position side: long, short
	PosSide string
	// MarginMode - Margin mode: crossed, isolated
	MarginMode string
	// PositionBalance - Position balance (margin amount), unit: margin coin; in isolated margin mode reflects the isolated margin amount for this position
	PositionBalance ujson.Float64
	// Available - Available position
	Available ujson.Float64
	// Frozen - Frozen position
	Frozen ujson.Float64
	// Total - Total position (available + frozen)
	Total ujson.Float64
	// Leverage - Leverage multiple
	Leverage ujson.Float64
	// CurRealisedPnl - Current realised profit and loss
	CurRealisedPnl ujson.Float64
	// AvgPrice - Average entry price
	AvgPrice ujson.Float64
	// PositionStatus - Position status: normal
	PositionStatus string
	// UnrealisedPnl - Unrealised profit and loss; in isolated margin mode reflects the unrealised PnL for this isolated position
	UnrealisedPnl ujson.Float64
	// LiquidationPrice - Estimated liquidation price; less than or equal to 0 means liquidation will not occur
	LiquidationPrice ujson.Float64
	// Mmr - Maintenance margin rate
	Mmr ujson.Float64
	// ProfitRate - Profit rate
	ProfitRate ujson.Float64
	// MarkPrice - Mark price
	MarkPrice ujson.Float64
	// BreakEvenPrice - Break-even price
	BreakEvenPrice ujson.Float64
	// TotalFunding - Total funding: the accumulated fund fee during the position's duration; zero indicates no fees have been charged
	TotalFunding ujson.Float64
	// OpenFeeTotal - Fees deducted on position opening during the position's lifetime
	OpenFeeTotal ujson.Float64
	// CloseFeeTotal - Fees deducted on position closing during the position's lifetime
	CloseFeeTotal ujson.Float64
	// CashDividend - Cash dividend, unit: USDT
	CashDividend ujson.Float64
	// CreatedTime - Created timestamp, unix millisecond timestamp
	CreatedTime ujson.Int64
	// UpdatedTime - Updated timestamp, unix millisecond timestamp
	UpdatedTime ujson.Int64
}
