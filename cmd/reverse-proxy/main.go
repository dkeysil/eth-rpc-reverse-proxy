package main

import (
	"fmt"

	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/config"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	zap.ReplaceGlobals(log)

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	client := client.NewReverseProxyClient(&fasthttp.Client{})

	backendResolver := backendresolver.NewResolver(config.Upstreams)

	service := &controllers.Service{
		Client:          client,
		BackendResolver: backendResolver,
	}

	log.Info("starting listening", zap.String("host", config.Server.Host), zap.String("port", config.Server.Port))
	fasthttp.ListenAndServe(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port), func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/eth_call":
			service.EthCallHandler(ctx)
		default:
			service.AnyHandler(ctx)
		}
	})
}
