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
	// Symbol - Symbol name, e.g. BTCUSDT; private account-wide topics (e.g. position) have no symbol
	Symbol string `json:"symbol,omitempty"`
}

// wsArgs - build subscription args from a REST Category (uppercase -> lowercase instType)
func wsArgs(topic string, category Category, symbol string) SubscriptionArgs {
	return SubscriptionArgs{
		InstType: strings.ToLower(string(category)),
		Topic:    topic,
		Symbol:   symbol,
	}
}

// wsPrivateArgs - build private subscription args: instType is the literal "UTA"
// (uppercase, unlike the lowercase instType of public channels) and there is no symbol
// https://www.bitget.com/api-doc/uta/websocket/private/Positions-Channel
func wsPrivateArgs(topic string) SubscriptionArgs {
	return SubscriptionArgs{
		InstType: "UTA",
		Topic:    topic,
	}
}

// WsLoginRequest - private WebSocket login operation
// https://www.bitget.com/api-doc/uta/websocket/private/WebSocket-Private
type WsLoginRequest struct {
	// Op - Operation: login
	Op string `json:"op"`
	// Args - Login credentials
	Args []LoginArgs `json:"args"`
}

// LoginArgs - private WebSocket login credentials
// https://www.bitget.com/api-doc/uta/websocket/private/WebSocket-Private
type LoginArgs struct {
	// ApiKey - API key
	ApiKey string `json:"apiKey"`
	// Passphrase - API key passphrase, plain text
	Passphrase string `json:"passphrase"`
	// Timestamp - Unix timestamp in seconds (unlike the REST millisecond timestamp)
	Timestamp string `json:"timestamp"`
	// Sign - base64(HMAC_SHA256(secret, timestamp + "GET" + "/user/verify"))
	Sign string `json:"sign"`
}
