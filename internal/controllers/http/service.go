package controllers

import (
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/backend_resolver"
	client "github.com/dkeysil/eth-rpc-reverse-proxy/internal/infrastructure/client/http"
)

type Service struct {
	Client          client.ReverseProxyClient
	BackendResolver backendresolver.BackendResolver
}
