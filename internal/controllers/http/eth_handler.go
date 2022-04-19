package controllers

import (
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/entities"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/metrics"
	"github.com/valyala/fasthttp"
)

func (s *Service) EthHandler(ctx *fasthttp.RequestCtx, method string) {
	request, err := entities.NewRequest(ctx.Request.Body())

	if err == nil && len(request.Method) > 0 {
		metrics.TotalHTTPRequests.WithLabelValues(request.Method).Inc()
		switch request.Method {
		case "eth_call":
			s.EthCallHandler(ctx)
		default:
			s.AnyMethod(ctx, request.Method)
		}

		return
	}
}
