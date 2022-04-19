package controllers

import (
	"github.com/valyala/fasthttp"
)

func (s *Service) AnyMethod(ctx *fasthttp.RequestCtx) {
	resChan := make(chan *fasthttp.Response)

	go s.Client.AsyncDo(ctx, s.BackendResolver.GetUpstreamHost("*"), resChan)

	res := <-resChan
	defer fasthttp.ReleaseResponse(res)

	res.CopyTo(&ctx.Response)
}
