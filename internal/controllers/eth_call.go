package controllers

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (s *Service) EthCallHandler(ctx *fasthttp.RequestCtx) {
	reqUpstream := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(reqUpstream)
	ctx.Request.CopyTo(reqUpstream)
	reqUpstream.SetHost(s.BackendResolver.GetUpstreamHost())

	reqEthCallUpstream := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(reqEthCallUpstream)
	ctx.Request.CopyTo(reqEthCallUpstream)
	reqEthCallUpstream.SetHost(s.BackendResolver.GetEthCallUpstreamHost())

	resChan := make(chan *fasthttp.Response, 2)

	go s.requestEthCall(reqUpstream, resChan)
	go s.requestEthCall(reqEthCallUpstream, resChan)

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

func (s *Service) requestEthCall(req *fasthttp.Request, resChan chan<- *fasthttp.Response) {
	// todo: reuse response with sync.Pool
	zap.L().Info("requesting eth_call", zap.ByteString("backend_host", req.Host()))
	res := fasthttp.AcquireResponse()

	s.Client.Do(req, res)
	resChan <- res
}
