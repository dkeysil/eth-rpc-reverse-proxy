package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dgrr/websocket"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/config"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/metrics"
	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
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
7. Docker [done]
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

	r := router.New()
	r.POST("/", func(ctx *fasthttp.RequestCtx) {
		var body Body
		err := json.Unmarshal(ctx.Request.Body(), &body)
		if err == nil && len(body.Method) > 0 {
			metrics.TotalHTTPRequests.WithLabelValues(body.Method).Inc()
			switch body.Method {
			case "eth_call":
				services.httpService.EthCallHandler(ctx)
			default:
				services.httpService.AnyMethod(ctx, body.Method)
			}

			return
		}
	})

	r.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

	r.ANY("/*", services.httpService.Handler)

	log.Info("starting listening", zap.String("host", config.Server.Host), zap.String("port", config.Server.Port))

	fasthttp.ListenAndServe(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port), func(ctx *fasthttp.RequestCtx) {
		zap.L().Debug("got new request")
		if bytes.Compare(ctx.Request.Header.Peek("Upgrade"), []byte("websocket")) == 0 {
			ws.Upgrade(ctx)
		} else {
			r.Handler(ctx)
		}
	})
}
