package bitgetv3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Sign - Bitget API credentials
// ACCESS-SIGN = base64(HMAC_SHA256(secret, timestamp + method.toUpperCase() + requestPath + "?" + queryString + body))
// "?" + queryString is omitted when the query is empty; query keys sorted in ascending order; body is empty for GET
// Headers: ACCESS-KEY, ACCESS-SIGN, ACCESS-TIMESTAMP (unix ms), ACCESS-PASSPHRASE (plain), Content-Type, locale
// The request expires 30 seconds after the timestamp
// https://www.bitget.com/api-doc/uta/guide
type Sign struct {
	// Key - API key (ACCESS-KEY header)
	Key string
	// Secret - API secret used as the HMAC-SHA256 key
	Secret string
	// Password - API passphrase (ACCESS-PASSPHRASE header, sent as plain text)
	Password string
}

func NewSign(key, secret, password string) *Sign {
	o := new(Sign)
	o.Key = key
	o.Secret = secret
	o.Password = password
	return o
}

func (o *Sign) timestamp() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

func (o *Sign) HeaderGet(h http.Header, v url.Values, path string, apiPath string) {
	o.header(h, v.Encode(), path, http.MethodGet, apiPath)
}

func (o *Sign) HeaderPost(h http.Header, body []byte, path string, apiPath string) {
	o.header(h, string(body), path, http.MethodPost, apiPath)
}

func (o *Sign) header(h http.Header, data string, path string, method string, apiPath string) {
	ts := o.timestamp()

	// Pre-sign: timestamp + method + /{apiPath}/path + "?" + query (GET) | body (POST)
	preSign := ts + method + "/" + apiPath + "/" + path
	if data != "" {
		if method == http.MethodGet {
			preSign += "?" + data
		} else {
			preSign += data
		}
	}

	h.Set("ACCESS-KEY", o.Key)
	h.Set("ACCESS-SIGN", signHmac(preSign, o.Secret))
	h.Set("ACCESS-TIMESTAMP", ts)
	h.Set("ACCESS-PASSPHRASE", o.Password)
	h.Set("Content-Type", "application/json")
	h.Set("locale", "en-US")
}

// WebSocket - build the login args for the private WebSocket
// sign = base64(HMAC_SHA256(secret, timestamp + "GET" + "/user/verify")), timestamp in unix seconds
// (unlike the REST millisecond timestamp)
// https://www.bitget.com/api-doc/uta/websocket/private/WebSocket-Private
func (o *Sign) WebSocket() LoginArgs {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	return LoginArgs{
		ApiKey:     o.Key,
		Passphrase: o.Password,
		Timestamp:  ts,
		Sign:       wsSign(ts, o.Secret),
	}
}

// wsSign - sign the WebSocket login pre-sign string for the given timestamp
func wsSign(ts, secret string) string {
	return signHmac(ts+"GET"+"/user/verify", secret)
}

func signHmac(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	io.WriteString(h, message)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
