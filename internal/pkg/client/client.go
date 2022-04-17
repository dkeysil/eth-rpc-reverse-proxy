package client

import (
	"github.com/valyala/fasthttp"
)

type ReverseProxyClient interface {
	Do(ctx *fasthttp.RequestCtx, host string, res *fasthttp.Response) error
}

type reverseProxyClient struct {
	client *fasthttp.Client
}

func NewReverseProxyClient(client *fasthttp.Client) ReverseProxyClient {
	return &reverseProxyClient{
		client: client,
	}
}

func (c *reverseProxyClient) Do(ctx *fasthttp.RequestCtx, host string, res *fasthttp.Response) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// pass forwarded headers to letting backend knows the initiator of the request
	req.Header.Add("X-Forwarded-For", ctx.RemoteAddr().String())
	req.Header.AddBytesKV([]byte("X-Forwarded-Proto"), req.Header.Protocol())
	req.Header.AddBytesKV([]byte("X-Forwarded-Host"), req.Header.Host())

	req.SetHost(host)

	return c.client.Do(req, res)
}
