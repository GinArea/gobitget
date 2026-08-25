package bitgetv3

import "net/http"

// Response - result of an API call as returned to the caller
type Response[T any] struct {
	// Time - server request time, unix ms (requestTime of the envelope)
	Time int64
	// Data - decoded and transformed payload
	Data T
	// Limit - rate-limit info (empty: Bitget sends no rate-limit headers)
	Limit RateLimit
	// Error - transport error or *Error with the Bitget code; nil on success
	Error error
	// StatusCode - HTTP status code of the response
	StatusCode int
	// NetError - true when the request failed before reaching the exchange
	NetError bool
}

// response - raw Bitget response envelope: {code, msg, requestTime, data}
type response[T any] struct {
	Code        string
	Msg         string
	RequestTime int64
	Data        T
}

func (o *Response[T]) Ok() bool {
	return o.Error == nil
}

// Retryable - response worth repeating (possibly through another proxy):
// a network failure before any answer, exchange/cloudflare unavailability (5xx, 52x),
// or 403/429 with which the edge layer replies to a banned or rate-limited IP.
// An answer from the exchange itself carrying its own error code (HTTP 400 on Bitget,
// unlike okx/bybit which answer 200 and put the code in the body) is not cured by a
// retry - the repeat returns exactly the same code.
func (o *Response[T]) Retryable() bool {
	return o.NetError ||
		o.StatusCode >= http.StatusInternalServerError ||
		o.StatusCode == http.StatusForbidden ||
		o.StatusCode == http.StatusTooManyRequests
}

func (o *Response[T]) SetErrorIfNil(err error) {
	if o.Error == nil {
		o.Error = err
	}
}

func (o *response[T]) Error() error {
	e := Error{
		Code: o.Code,
		Text: o.Msg,
	}
	return e.Std()
}
