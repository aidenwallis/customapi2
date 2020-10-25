package version

import (
	"github.com/aidenwallis/customapi2/pkg/cache"
	"github.com/aidenwallis/customapi2/pkg/data"
	"github.com/aidenwallis/customapi2/pkg/feature"
	"github.com/aidenwallis/customapi2/pkg/twitchapi"
	"github.com/go-chi/chi"
)

type Version struct {
	*chi.Mux

	featureCfg *feature.Config
	Routes     []*feature.Route
}

type Config struct {
	Cache     cache.Cache
	Data      *data.Data
	TwitchAPI twitchapi.TwitchAPI
}

func New(cfg *Config) *Version {
	r := chi.NewRouter()

	return &Version{
		Mux: r,

		featureCfg: &feature.Config{
			Cache:     cfg.Cache,
			Data:      cfg.Data,
			TwitchAPI: cfg.TwitchAPI,
		},
	}
}

func (v *Version) Register(factory func(*feature.Feature)) {
	f := feature.New(v.featureCfg)
	factory(f)

	for _, r := range f.Routes() {
		r.Inject(v.Mux)
		v.Routes = append(v.Routes, r)
	}
}
