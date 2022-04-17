package controllers

import (
	"github.com/valyala/fasthttp"
)

func (s *Service) AnyHandler(ctx *fasthttp.RequestCtx) {
	host := s.BackendResolver.GetUpstreamHost(string(ctx.Path()))
	s.Client.Do(ctx, host, &ctx.Response)
}
