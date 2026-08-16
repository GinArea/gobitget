package bitgetv3

// WsRequestV2 - legacy v2 WebSocket operation request (subscribe/unsubscribe)
// https://www.bitget.com/api-doc/spot/websocket/public/Candlesticks-Channel
type WsRequestV2 struct {
	// Op - Operation: subscribe / unsubscribe
	Op string `json:"op"`
	// Args - Subscribed channels (batching multiple channels in one request is allowed)
	Args []SubscriptionArgsV2 `json:"args"`
}

// SubscriptionArgsV2 - subscribed v2 channel; comparable, used as the subscription map key
// Unlike the v3 args there is no separate topic/interval: the timeframe is embedded
// in the channel name (candle1m ... candle1Mutc)
// https://www.bitget.com/api-doc/spot/websocket/public/Candlesticks-Channel
type SubscriptionArgsV2 struct {
	// InstType - Product type: SPOT / USDT-FUTURES / COIN-FUTURES / USDC-FUTURES (uppercase, unlike the v3 WS)
	InstType string `json:"instType"`
	// Channel - Channel name with the embedded timeframe, e.g. candle1Dutc
	Channel string `json:"channel"`
	// InstId - Symbol name, e.g. BTCUSDT
	InstId string `json:"instId"`
}

// wsArgsV2 - build v2 subscription args from a REST Category (used as is: the v2 WS wants uppercase)
func wsArgsV2(channel string, category Category, symbol string) SubscriptionArgsV2 {
	return SubscriptionArgsV2{
		InstType: string(category),
		Channel:  channel,
		InstId:   symbol,
	}
}
