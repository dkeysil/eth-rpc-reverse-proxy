package controllers

import (
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/backend_resolver"
	c "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/client/ws"
	resolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/id_resolver"
)

type Service struct {
	Client          c.WSReverseProxyClient
	BackendResolver backendresolver.BackendResolver
	IDResolver      resolver.IDResolver
}
