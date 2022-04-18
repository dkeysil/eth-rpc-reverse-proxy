package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dgrr/websocket"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/config"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Body struct {
	Method string `json:"method"`
}

/*
TODO:
1. Websockets [done]
2. Support of removing dead backends
3. Retries
4. Remove superfluous type conversion (string -> []byte)
5. More logs
6. Prometheus + Grafana
7. Docker
*/

func main() {
	log, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(log)

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	resolvers := NewResolvers(config)

	clients := NewClients(resolvers)

	services := NewServices(clients, resolvers)

	ws := websocket.Server{}
	ws.HandleData(services.websocketService.OnMessage)

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
					services.httpService.EthCallHandler(ctx)
				default:
					services.httpService.AnyMethod(ctx, body.Method)
				}

				return
			}

			services.httpService.Handler(ctx)
		}
	})
}
