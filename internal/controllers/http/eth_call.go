package controllers

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (s *Service) EthCallHandler(ctx *fasthttp.RequestCtx) {
	resChan := make(chan *fasthttp.Response, 2)
	errChan := make(chan error, 2)
	defer close(resChan)
	defer close(errChan)

	go s.asyncRequest(ctx, s.BackendResolver.GetUpstreamHost("eth_call"), resChan, errChan)
	go s.asyncRequest(ctx, s.BackendResolver.GetUpstreamHost(""), resChan, errChan)

	for i := 0; i < 2; i++ {
		select {
		case res := <-resChan:
			defer fasthttp.ReleaseResponse(res)
			if res.StatusCode() >= 400 {
				continue
			}

			zap.L().Info("response eth_call", zap.String("backend_host", res.RemoteAddr().String()))

			res.CopyTo(&ctx.Response)
			return
		case err := <-errChan:
			zap.L().Error("backend is down", zap.Error(err))
		case <-ctx.Done():
			break
		}
	}

	zap.L().Error("no one backend upstream works")
	ctx.Response = fasthttp.Response{
		Header: fasthttp.ResponseHeader{},
	}
	ctx.Response.Header.SetStatusCode(500)
	ctx.Response.SetBodyString("all backends is down")
}

func (s *Service) asyncRequest(ctx *fasthttp.RequestCtx, host string, resChan chan<- *fasthttp.Response, errChan chan<- error) {
	zap.L().Info("requesting eth_call", zap.String("backend_host", host))
	res := fasthttp.AcquireResponse()

	defer func() {
		if recover() != nil {
			fasthttp.ReleaseResponse(res)
		}
	}()

	if err := s.Client.Do(ctx, host, res); err != nil {
		errChan <- err
		fasthttp.ReleaseResponse(res)
	} else {
		resChan <- res
	}
}
