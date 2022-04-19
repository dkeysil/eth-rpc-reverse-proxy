package controllers

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (s *Service) AnyMethod(ctx *fasthttp.RequestCtx, method string) {
	host := s.BackendResolver.GetUpstreamHost(method)
	err := s.Client.Do(ctx, host, &ctx.Response)
	if err != nil {
		zap.L().Error("error while calling any method", zap.String("method", method), zap.String("backend_host", ctx.Response.RemoteAddr().String()))
		s.BackendResolver.RemoveHost(host)
	}
}
