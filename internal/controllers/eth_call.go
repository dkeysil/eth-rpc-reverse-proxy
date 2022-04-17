package controllers

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (s *Service) EthCallHandler(ctx *fasthttp.RequestCtx) {
	resChan := make(chan *fasthttp.Response, 2)

	go s.asyncRequest(ctx, s.BackendResolver.GetUpstreamHost(string(ctx.Path())), resChan)
	go s.asyncRequest(ctx, s.BackendResolver.GetUpstreamHost(""), resChan)

	for i := 0; i < 2; i++ {
		select {
		case res := <-resChan:
			defer fasthttp.ReleaseResponse(res)
			if !res.ConnectionClose() || res.StatusCode() >= 400 {
				continue
			}

			zap.L().Info("response eth_call", zap.String("backend_host", res.RemoteAddr().String()))

			res.CopyTo(&ctx.Response)
			return
		case <-ctx.Done():
			break
		}
	}

	zap.L().Error("no one backend upstream works")
	ctx.Response = fasthttp.Response{
		Header: fasthttp.ResponseHeader{},
	}

	ctx.Response.Header.SetStatusCode(500)
}

func (s *Service) asyncRequest(ctx *fasthttp.RequestCtx, host string, resChan chan<- *fasthttp.Response) {
	zap.L().Info("requesting eth_call", zap.String("backend_host", host))
	res := fasthttp.AcquireResponse()

	s.Client.Do(ctx, host, res)
	resChan <- res
}
