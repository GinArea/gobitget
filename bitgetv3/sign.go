package bitgetv3

// Sign - Bitget API credentials
// Signing scheme (implemented with the first private endpoint):
// base64(HMAC_SHA256(timestamp + method + requestPath + query/body))
// Headers: ACCESS-KEY, ACCESS-SIGN, ACCESS-TIMESTAMP, ACCESS-PASSPHRASE
// https://www.bitget.com/api-doc/common/signature
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
