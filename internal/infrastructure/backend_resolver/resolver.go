package backendresolver

import (
	"sync"
	"sync/atomic"
)

//go:generate mockgen -source=resolver.go -destination=mocks/resolver_mock.go BackendResolver
type BackendResolver interface {
	GetUpstreamHost(path string) string
	GetAllUpstreams(path string) []string
	RemoveHost(host string)
}

type resolver struct {
	upstreams map[string][]string
	counters  map[string]*uint64
	sync.RWMutex
}

// NewResolver backend resolver constructor
// upstreams map must have "*" key for the base upstream host getter
func NewResolver(upstreams map[string][]string) BackendResolver {
	if _, ok := upstreams["*"]; !ok {
		panic("upstreams map must have \"*\" key for the base upstream host getter")
	}

	resolver := &resolver{
		upstreams: upstreams,
		counters:  make(map[string]*uint64),
	}

	for key := range upstreams {
		counter := uint64(0)
		resolver.counters[key] = &counter
	}

	return resolver
}

// GetUpstreamHost returns upstream host selected with round robin
// if upstreams list not found by path - will return host from base upstream list
func (r *resolver) GetUpstreamHost(path string) string {
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()

	if _, ok := r.upstreams[path]; !ok {
		path = "*"
	}

	defer atomic.AddUint64(r.counters[path], 1)
	// first instance will have slightly more load
	return r.upstreams[path][atomic.LoadUint64(r.counters[path])%uint64(len(r.upstreams[path]))]
}

func (r *resolver) GetAllUpstreams(path string) []string {
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()

	return r.upstreams[path]
}

func (r *resolver) RemoveHost(host string) {
	r.RWMutex.Lock()
	for path, upstreams := range r.upstreams {
		for i, h := range upstreams {
			if h == host {
				if v := r.upstreams[path]; len(v) == 1 {
					delete(r.upstreams, h)
				} else {
					r.upstreams[path] = append(r.upstreams[path][:i], r.upstreams[path][i+1:]...)
				}
			}
		}

	}

	if _, ok := r.upstreams["*"]; !ok {
		panic("all backends removed")
	}
	r.RWMutex.Unlock()
}
