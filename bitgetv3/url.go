package bitgetv3

const (
	// REST API (UTA)
	MainBaseUrl = "https://api.bitget.com"
	ApiVersion  = "api/v3"
	// REST API (classic v2, used by affiliate/broker endpoints)
	ApiVersion2 = "api/v2"
	// WebSocket API (UTA)
	// https://www.bitget.com/api-doc/uta/guide
	MainBaseWsUrl = "wss://ws.bitget.com"
	WsPublicPath  = "v3/ws/public"
	// WebSocket API (legacy v2): used for kline timeframes the v3 WS lacks
	// (2H, 8H and the utc-suffixed variants) - see BITGET_API.md
	WsPublicV2Path = "v2/ws/public"
)
