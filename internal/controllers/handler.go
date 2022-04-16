package controllers

import (
	"github.com/valyala/fasthttp"
)

func (h *Service) AnyHandler(ctx *fasthttp.RequestCtx) {
	host := ctx.Request.Host()
	_ = host
	req := &fasthttp.Request{}
	ctx.Request.CopyTo(req)

	req.SetHost(h.BackendResolver.GetUpstreamHost())
	h.Client.Do(req, &ctx.Response)
	return
}
