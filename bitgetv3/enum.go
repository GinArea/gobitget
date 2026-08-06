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
