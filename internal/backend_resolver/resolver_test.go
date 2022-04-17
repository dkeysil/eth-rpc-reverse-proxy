package backendresolver

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewResolver(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("code did not panic")
		}
	}()

	NewResolver(map[string][]string{"not_base_path": {"localhost:8080"}})
}

func TestGetUpstreamHost(t *testing.T) {
	is := is.New(t)

	upstreams := map[string][]string{
		"*":        {"localhost:8001", "localhost:8002", "localhost:8003"},
		"eth_call": {"localhost:8001", "localhost:8002"},
	}
	resolver := NewResolver(upstreams)

	t.Run("path not in upstreams map", func(t *testing.T) {
		expectedBaseHost := []string{"localhost:8001", "localhost:8002", "localhost:8003", "localhost:8001"}

		for _, host := range expectedBaseHost {
			is.Equal(resolver.GetUpstreamHost(""), host)
		}
	})

	t.Run("specify upstream host list by path", func(t *testing.T) {
		expectedBaseHost := []string{"localhost:8001", "localhost:8002", "localhost:8001", "localhost:8002"}

		for _, host := range expectedBaseHost {
			is.Equal(resolver.GetUpstreamHost("eth_call"), host)
		}
	})

}
