package bitgetv3

import (
	"net/http"
	"time"

	"github.com/msw-x/moon/uhttp"
)

type Client struct {
	c                *uhttp.Client
	s                *Sign
	onTransportError OnTransportError
	apiPath          string // API path for signing (e.g., "api/v3")
}

func NewClient() *Client {
	o := new(Client)
	o.c = uhttp.NewClient()
	o.apiPath = ApiVersion
	o.WithBaseUrl(MainBaseUrl)
	o.WithPath(ApiVersion)
	return o
}

func (o *Client) WithTimeout(timeout time.Duration) *Client {
	o.c.WithTimeout(timeout)
	return o
}

func (o *Client) WithTrace(trace func(uhttp.Response)) *Client {
	o.c.WithTrace(trace)
	return o
}

func (o *Client) WithTransport(tranport *http.Transport) *Client {
	o.c.WithTransport(tranport)
	return o
}

func (o *Client) WithProxy(proxy string) *Client {
	o.c.WithProxy(proxy)
	return o
}

func (o *Client) Copy() *Client {
	r := new(Client)
	r.c = o.c.Copy()
	r.s = o.s
	r.onTransportError = o.onTransportError
	r.apiPath = o.apiPath
	return r
}

func (o *Client) WithBaseUrl(url string) *Client {
	o.c.WithBase(url)
	return o
}

func (o *Client) WithPath(path string) *Client {
	o.c.WithPath(path)
	o.apiPath = path
	return o
}

func (o *Client) WithAppendPath(path string) *Client {
	o.c.WithAppendPath(path)
	return o
}

func (o *Client) WithAuth(key, secret, password string) *Client {
	o.s = NewSign(key, secret, password)
	return o
}

func (o *Client) WithOnReadBodyError(f uhttp.OnError) *Client {
	o.c.WithOnReadBodyError(f)
	return o
}

func (o *Client) WithOnTransportError(f OnTransportError) *Client {
	o.onTransportError = f
	return o
}

// Path helpers for public endpoints only
// Note: Private endpoints pass path directly to Get/Post for correct signing

func (o *Client) market() *Client {
	return o.Copy().WithAppendPath("market")
}

// v2 - client copy targeting the classic v2 API (affiliate/broker endpoints)
// The signing scheme is identical to UTA v3, only the api path differs
func (o *Client) v2() *Client {
	return o.Copy().WithPath(ApiVersion2)
}

type OnTransportError func(err error, method string, statusCode int, attempt int) bool
