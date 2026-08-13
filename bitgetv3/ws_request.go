package bitgetv3

import (
	"strings"
)

// WsRequest - WebSocket operation request (subscribe/unsubscribe)
// https://www.bitget.com/api-doc/uta/websocket/public/Tickers-Channel
type WsRequest struct {
	// Op - Operation: subscribe / unsubscribe
	Op string `json:"op"`
	// Args - Subscribed channels (batching multiple channels in one request is allowed)
	Args []SubscriptionArgs `json:"args"`
}

// SubscriptionArgs - subscribed channel; comparable, used as the subscription map key
// https://www.bitget.com/api-doc/uta/websocket/public/Tickers-Channel
type SubscriptionArgs struct {
	// InstType - Product type: spot / usdt-futures / coin-futures / usdc-futures (lowercase, unlike REST Category)
	InstType string `json:"instType"`
	// Topic - Topic name: ticker / books / kline ...
	Topic string `json:"topic"`
	// Symbol - Symbol name, e.g. BTCUSDT
	Symbol string `json:"symbol"`
}

// wsArgs - build subscription args from a REST Category (uppercase -> lowercase instType)
func wsArgs(topic string, category Category, symbol string) SubscriptionArgs {
	return SubscriptionArgs{
		InstType: strings.ToLower(string(category)),
		Topic:    topic,
		Symbol:   symbol,
	}
}
