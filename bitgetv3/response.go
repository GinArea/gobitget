package bitgetv3

type Response[T any] struct {
	Time       int64
	Data       T
	Limit      RateLimit
	Error      error
	StatusCode int
	NetError   bool
}

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
