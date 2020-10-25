package twitch

import (
	"regexp"

	"github.com/aidenwallis/customapi2/pkg/feature"
	"github.com/aidenwallis/customapi2/pkg/twitchapi"
)

var loginRegex = regexp.MustCompile(`^[\w\d_]{2,}$`)

type twitchUser struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
	Bio         string `json:"bio"`
	ViewCount   int    `json:"viewCount"`
}

type twitchFeature struct {
	*feature.Feature
}

func New(parent *feature.Feature) {
	f := &twitchFeature{Feature: parent}

	f.Get("/twitch/id/{login}", f.getUserID)
}

func makeTwitchUser(user *twitchapi.User) *twitchUser {
	return &twitchUser{
		ID:          user.ID,
		Login:       user.Login,
		DisplayName: user.DisplayName,
		Avatar:      user.ProfileImageURL,
		Bio:         user.Description,
		ViewCount:   user.ViewCount,
	}
}
