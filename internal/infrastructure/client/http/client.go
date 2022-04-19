package client

import (
	"path"

	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/backend_resolver"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type ReverseProxyClient interface {
	AsyncDo(ctx *fasthttp.RequestCtx, host string, resChan chan<- *fasthttp.Response)
}

type reverseProxyClient struct {
	client          *fasthttp.Client
	backendResolver backendresolver.BackendResolver
}

func NewReverseProxyClient(client *fasthttp.Client, backendResolver backendresolver.BackendResolver) ReverseProxyClient {
	return &reverseProxyClient{
		client:          client,
		backendResolver: backendResolver,
	}
}

// Do copying users request, add X-Forwarded headers, changing URI to one of the upstreams
func (c *reverseProxyClient) AsyncDo(ctx *fasthttp.RequestCtx, host string, resChan chan<- *fasthttp.Response) {
	zap.L().Debug("async request to backend", zap.String("host", host))

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

	res := fasthttp.AcquireResponse()
	defer func() {
		// first request can close channel faster than second got response
		if recover() != nil {
			fasthttp.ReleaseResponse(res)
		}
	}()

	err := c.client.Do(req, res)
	if err != nil {
		// here we can inject additional error info into the response, like number of retries and backend urls if needed, but
		// Requirement: Mirroring should be made transparently for clients.
		res.SetStatusCode(503)
		zap.L().Error("backend is down", zap.Error(err), zap.String("host", host))
		c.backendResolver.RemoveHost(host)
	}

	resChan <- res
}
