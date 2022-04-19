package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/dgrr/websocket"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/config"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/metrics"
	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.uber.org/zap"
)

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
	r.POST("/", services.httpService.EthHandler)

	// prometheus metrics handler
	r.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

	log.Info("starting listening", zap.String("host", config.Server.Host), zap.String("port", config.Server.Port))
	fasthttp.ListenAndServe(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port), func(ctx *fasthttp.RequestCtx) {
		zap.L().Info("new request", zap.String("host", ctx.RemoteIP().String()))
		if bytes.Compare(ctx.Request.Header.Peek("Upgrade"), []byte("websocket")) == 0 {
			ws.Upgrade(ctx)
		} else {
			r.Handler(ctx)
			metrics.ResponseCodes.WithLabelValues(strconv.Itoa(ctx.Response.StatusCode())).Inc()
		}
	})
}
