package feature

import (
	"net/http"

	"github.com/aidenwallis/customapi2/pkg/cache"
	"github.com/aidenwallis/customapi2/pkg/data"
	"github.com/aidenwallis/customapi2/pkg/twitchapi"
)

type Config struct {
	Cache     cache.Cache
	Data      *data.Data
	TwitchAPI twitchapi.TwitchAPI
}

type Feature struct {
	routes []*Route

	Cache     cache.Cache
	Data      *data.Data
	TwitchAPI twitchapi.TwitchAPI
}

func New(cfg *Config) *Feature {
	return &Feature{
		Cache:     cfg.Cache,
		Data:      cfg.Data,
		TwitchAPI: cfg.TwitchAPI,
	}
}

func (f *Feature) Get(pattern string, handler handlerFunc) *Route {
	return f.registerRoute(http.MethodGet, pattern, handler)
}

func (f *Feature) Post(pattern string, handler handlerFunc) *Route {
	return f.registerRoute(http.MethodPost, pattern, handler)
}

func (f *Feature) Routes() []*Route {
	return f.routes
}

func (f *Feature) registerRoute(method, pattern string, handler handlerFunc) *Route {
	r := &Route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
	f.routes = append(f.routes, r)
	return r
}
