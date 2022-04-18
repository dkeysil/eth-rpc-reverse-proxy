package controllers

import (
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	resolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/id_resolver"
	c "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client/ws"
)

type Service struct {
	Client          c.WSReverseProxyClient
	BackendResolver backendresolver.BackendResolver
	IDResolver      resolver.IDResolver
}
