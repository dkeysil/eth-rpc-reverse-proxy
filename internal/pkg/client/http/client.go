package client

import (
	"bytes"
	"fmt"
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

func (c *reverseProxyClient) Do(ctx *fasthttp.RequestCtx, host string, res *fasthttp.Response) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	ctx.Request.CopyTo(req)

	// pass forwarded headers to letting backend knows the initiator of the request
	req.Header.Add("X-Forwarded-For", ctx.RemoteIP().String())
	req.Header.AddBytesKV([]byte("X-Forwarded-Proto"), ctx.Request.Header.Protocol())
	req.Header.AddBytesKV([]byte("X-Forwarded-Host"), ctx.Request.Header.Host())

	uri := &fasthttp.URI{}
	uri.Parse([]byte{}, []byte(host))
	path := path.Join(string(uri.Path()), string(ctx.Path()))
	uri.SetPath(path)
	req.SetURI(uri)

	if bytes.Equal(req.Header.Method(), []byte("GET")) {
		fmt.Println(req)
	}

	return c.client.Do(req, res)
}
