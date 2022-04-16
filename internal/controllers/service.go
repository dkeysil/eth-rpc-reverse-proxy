package controllers

import (
	backendresolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/backend_resolver"
	"github.com/valyala/fasthttp"
)

type Service struct {
	Client          fasthttp.Client
	BackendResolver backendresolver.BackendResolver
}
