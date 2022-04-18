package controllers

import (
	"github.com/valyala/fasthttp"
)

func (s *Service) Handler(ctx *fasthttp.RequestCtx) {
	host := s.BackendResolver.GetUpstreamHost("")
	s.Client.Do(ctx, host, &ctx.Response)
}
