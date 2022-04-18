package main

import (
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/config"
	httpControllers "github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers/http"
	wsControllers "github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers/ws"
	resolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/id_resolver"
	client "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client/http"
	wsClient "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client/ws"
	"github.com/valyala/fasthttp"
)

type Clients struct {
	httpClient      client.ReverseProxyClient
	websocketClient wsClient.WSReverseProxyClient
}

type Services struct {
	httpService      httpControllers.Service
	websocketService wsControllers.Service
}

type Resolvers struct {
	backendResolver backendresolver.BackendResolver
	idResolver      resolver.IDResolver
}

func NewClients(resolvers Resolvers) Clients {
	httpClient := client.NewReverseProxyClient(&fasthttp.Client{})
	wsClient := wsClient.NewWSReverseProxyClient(append(resolvers.backendResolver.GetAllUpstreams("ws:*"), resolvers.backendResolver.GetAllUpstreams("ws:eth_call")...), resolvers.idResolver)

	return Clients{
		httpClient:      httpClient,
		websocketClient: wsClient,
	}
}

func NewServices(clients Clients, resolvers Resolvers) Services {
	httpService := httpControllers.Service{
		Client:          clients.httpClient,
		BackendResolver: resolvers.backendResolver,
	}
	wsService := wsControllers.Service{
		Client:          clients.websocketClient,
		BackendResolver: resolvers.backendResolver,
		IDResolver:      resolvers.idResolver,
	}

	return Services{
		httpService:      httpService,
		websocketService: wsService,
	}
}

func NewResolvers(config config.Config) Resolvers {
	return Resolvers{
		backendResolver: backendresolver.NewResolver(config.Upstreams),
		idResolver:      resolver.NewIDResolver(),
	}

}
