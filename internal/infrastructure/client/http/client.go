package client

import (
	"path"

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

// Do copying users request, add X-Forwarded headers, changing URI to one of the upstreams
func (c *reverseProxyClient) Do(ctx *fasthttp.RequestCtx, host string, res *fasthttp.Response) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	ctx.Request.CopyTo(req)

	// pass forwarded headers to letting backend knows the initiator of the request
	req.Header.Add("X-Forwarded-For", ctx.RemoteIP().String())
	req.Header.AddBytesKV([]byte("X-Forwarded-Proto"), ctx.Request.Header.Protocol())
	req.Header.AddBytesKV([]byte("X-Forwarded-Host"), ctx.Request.Header.Host())

	// extend upstreams path with user path
	upstreamURI := &fasthttp.URI{}
	upstreamURI.Parse([]byte{}, []byte(host))
	path := path.Join(string(upstreamURI.Path()), string(ctx.Path()))
	upstreamURI.SetPath(path)
	req.SetURI(upstreamURI)

	return c.client.Do(req, res)
}
