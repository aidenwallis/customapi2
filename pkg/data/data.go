package data

import (
	"github.com/aidenwallis/customapi2/pkg/cache"
	"github.com/aidenwallis/customapi2/pkg/data/twitch"
	"github.com/aidenwallis/customapi2/pkg/twitchapi"
)

type Config struct {
	Cache     cache.Cache
	TwitchAPI twitchapi.TwitchAPI
}

type Data struct {
	Twitch twitch.Twitch
}

func New(cfg *Config) *Data {
	return &Data{
		Twitch: twitch.New(cfg.Cache, cfg.TwitchAPI),
	}
}
