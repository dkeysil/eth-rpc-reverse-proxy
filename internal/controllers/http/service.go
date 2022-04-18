package controllers

import (
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	client "github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client/http"
)

type Service struct {
	Client          client.ReverseProxyClient
	BackendResolver backendresolver.BackendResolver
}
