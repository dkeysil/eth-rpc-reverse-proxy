package main

import (
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/config"
	httpControllers "github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers/http"
	wsControllers "github.com/dkeysil/eth-rpc-reverse-proxy/internal/controllers/ws"
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/backend_resolver"
	client "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/client/http"
	wsClient "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/client/ws"
	resolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/id_resolver"
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
	httpBackendResolver backendresolver.BackendResolver
	wsBackendResolver   backendresolver.BackendResolver
	idResolver          resolver.IDResolver
}

func NewClients(resolvers Resolvers) Clients {
	fasthttpClient := &fasthttp.Client{
		RetryIf: func(request *fasthttp.Request) bool {
			// if I correctly understand - all methods in eth-rpc is idempotent
			// and for example, you can't double spend eth because signed transaction can be executed only once
			return true
		},
	}
	httpClient := client.NewReverseProxyClient(fasthttpClient, resolvers.httpBackendResolver)
	wsClient := wsClient.NewWSReverseProxyClient(
		append(resolvers.wsBackendResolver.GetAllUpstreams("*"), resolvers.wsBackendResolver.GetAllUpstreams("eth_call")...),
		resolvers.idResolver,
		resolvers.wsBackendResolver,
	)

	return Clients{
		httpClient:      httpClient,
		websocketClient: wsClient,
	}
}

func NewServices(clients Clients, resolvers Resolvers) Services {
	httpService := httpControllers.Service{
		Client:          clients.httpClient,
		BackendResolver: resolvers.httpBackendResolver,
	}
	wsService := wsControllers.Service{
		Client:          clients.websocketClient,
		BackendResolver: resolvers.wsBackendResolver,
		IDResolver:      resolvers.idResolver,
	}

	return Services{
		httpService:      httpService,
		websocketService: wsService,
	}
}

func NewResolvers(config config.Config) Resolvers {
	return Resolvers{
		httpBackendResolver: backendresolver.NewResolver(config.HTTPUpstreams),
		wsBackendResolver:   backendresolver.NewResolver(config.WSUpstreams),
		idResolver:          resolver.NewIDResolver(),
	}

}
