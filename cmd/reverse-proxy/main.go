package main

import (
	"encoding/json"
	"fmt"

	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/config"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Body struct {
	Method string `json:"method"`
}

/*
TODO:
1. Websockets
2. Support of removing dead backends
3. Retries
4. More logs
5. Prometheus + Grafana
6. Docker
*/

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
		var body Body
		err := json.Unmarshal(ctx.Request.Body(), &body)
		if err == nil && len(body.Method) > 0 {
			switch body.Method {
			case "eth_call":
				service.EthCallHandler(ctx)
			default:
				service.AnyMethod(ctx, body.Method)
			}

			return
		}

		service.Handler(ctx)
	})
}
