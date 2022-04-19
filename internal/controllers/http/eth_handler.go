package controllers

import (
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/entities"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/metrics"
	"github.com/valyala/fasthttp"
)

// EthHandler receive reqeust and check eth-rpc request method
// if method == eth_call - execute special logic for this case
func (s *Service) EthHandler(ctx *fasthttp.RequestCtx) {
	request, err := entities.NewRequest(ctx.Request.Body())

	if err == nil {
		metrics.TotalHTTPRequests.WithLabelValues(request.Method).Inc()
		switch request.Method {
		case "eth_call":
			s.EthCallHandler(ctx)
		default:
			s.AnyMethod(ctx)
		}

		return
	}
}
