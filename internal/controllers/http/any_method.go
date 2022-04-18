package controllers

import (
	"github.com/valyala/fasthttp"
)

func (s *Service) AnyMethod(ctx *fasthttp.RequestCtx, method string) {
	host := s.BackendResolver.GetUpstreamHost(method)
	s.Client.Do(ctx, host, &ctx.Response)
}
