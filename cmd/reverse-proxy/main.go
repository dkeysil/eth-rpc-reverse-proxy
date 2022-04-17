package main

import (
	"fmt"

	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	zap.ReplaceGlobals(log)

	client := client.NewReverseProxyClient(&fasthttp.Client{})

	backendResolver := backendresolver.NewResolver(map[string][]string{"*": {"localhost:8081"}, "/eth_call": {"localhost:8082"}})

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
