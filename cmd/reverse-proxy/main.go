package main

import (
	"fmt"

	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	zap.ReplaceGlobals(log)

	client := fasthttp.Client{}
	backendResolver := backendresolver.NewBackends([]string{"localhost:8081"}, []string{"localhost:8082"})

	service := &controllers.Service{
		Client:          client,
		BackendResolver: backendResolver,
	}

	port := 8080

	log.Info("starting listening", zap.Int("port", port))

	fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/eth_call":
			service.EthCallHandler(ctx)
		default:
			service.AnyHandler(ctx)
		}
	})
}
