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
	Key      string
	Secret   string
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

	// Pre-sign: timestamp + method + /api/v3/path + "?" + query (GET) | body (POST)
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

func signHmac(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	io.WriteString(h, message)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
