package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dgrr/websocket"
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/config"
	controllers "github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers/http"
	wsControllers "github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers/ws"
	client "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client/http"
	wsClient "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client/ws"
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
4. Remove superfluous type conversion (string -> []byte)
5. More logs
6. Prometheus + Grafana
7. Docker
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
	wsService := &wsControllers.Service{
		Client:          wsClient.NewWSReverseProxyClient(backendResolver.GetAllUpstreams("ws:*")),
		BackendResolver: backendResolver,
	}

	ws := websocket.Server{}
	ws.HandleData(wsService.OnMessage)

	log.Info("starting listening", zap.String("host", config.Server.Host), zap.String("port", config.Server.Port))
	fasthttp.ListenAndServe(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port), func(ctx *fasthttp.RequestCtx) {
		if bytes.Compare(ctx.Request.Header.Peek("Upgrade"), []byte("websocket")) == 0 {
			ws.Upgrade(ctx)
		} else {
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
		}
	})
}
