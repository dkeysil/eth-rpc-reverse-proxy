package controllers

import (
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/backend_resolver"
	client "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/client/http"
	resolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/id_resolver"
)

type Service struct {
	Client          client.ReverseProxyClient
	BackendResolver backendresolver.BackendResolver
	IDResolver      resolver.IDResolver
}
