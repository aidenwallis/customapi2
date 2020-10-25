package v2

import (
	"github.com/aidenwallis/customapi2/api/v2/twitch"
	"github.com/aidenwallis/customapi2/pkg/version"
)

func New(cfg *version.Config) *version.Version {
	v := version.New(cfg)

	v.Register(twitch.New)

	return v
}
