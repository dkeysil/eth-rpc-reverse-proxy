package controllers

import (
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/pkg/client"
)

type Service struct {
	Client          client.ReverseProxyClient
	BackendResolver backendresolver.BackendResolver
}
