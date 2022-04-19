package controllers

import (
	"github.com/valyala/fasthttp"
)

// EthCallHandler has to be duplicated to another list of backends
// response - is the fastest response from backends
func (s *Service) EthCallHandler(ctx *fasthttp.RequestCtx) {
	resChan := make(chan *fasthttp.Response)

	go s.Client.AsyncDo(ctx, s.BackendResolver.GetUpstreamHost("eth_call"), resChan)
	go s.Client.AsyncDo(ctx, s.BackendResolver.GetUpstreamHost("*"), resChan)

	res := <-resChan
	defer fasthttp.ReleaseResponse(res)

	if res.StatusCode() >= 400 {
		res = <-resChan
		defer fasthttp.ReleaseResponse(res)
	}

	res.CopyTo(&ctx.Response)
}
