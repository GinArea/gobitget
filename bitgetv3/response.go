package bitgetv3

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
