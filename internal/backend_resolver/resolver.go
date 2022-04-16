package backendresolver

import (
	_ "github.com/golang/mock/mockgen"
)

//go:generate mockgen -source=resolver.go -destination=mocks/resolver_mock.go BackendResolver
type BackendResolver interface {
	GetUpstreamHost() string
	GetEthCallUpstreamHost() string
}

type backends struct {
	upstream        []string
	ethCallUpstream []string
}

func NewBackends(upstream, ethCallUpstream []string) BackendResolver {
	if len(upstream) == 0 {
		panic("upstream list is empty")
	}

	if len(ethCallUpstream) == 0 {
		panic("eth call upstream list is empty")
	}

	return &backends{
		upstream:        upstream,
		ethCallUpstream: ethCallUpstream,
	}
}

func (br *backends) GetUpstreamHost() string {
	return br.upstream[len(br.upstream)-1]
}

func (br *backends) GetEthCallUpstreamHost() string {
	return br.ethCallUpstream[len(br.ethCallUpstream)-1]
}
