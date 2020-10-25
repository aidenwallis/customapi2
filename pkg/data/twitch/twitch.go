package twitch

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aidenwallis/customapi2/pkg/cache"
	"github.com/aidenwallis/customapi2/pkg/twitchapi"
)

type Twitch interface {
	GetUserByLogin(ctx context.Context, login string) (*twitchapi.User, error)
}

type twitch struct {
	cache     cache.Cache
	twitchAPI twitchapi.TwitchAPI
}

var ErrUserNotFound = errors.New("twitch: user not found")

const (
	cachePrefix        = "twitchapi::"
	usersByIDPrefix    = cachePrefix + "users-by-id::"
	usersByLoginPrefix = cachePrefix + "users-by-login::"
)

func New(cacheImpl cache.Cache, twitchAPI twitchapi.TwitchAPI) Twitch {
	return &twitch{
		cache:     cacheImpl,
		twitchAPI: twitchAPI,
	}
}

func (t *twitch) GetUserByLogin(ctx context.Context, login string) (*twitchapi.User, error) {
	login = strings.ToLower(strings.TrimSpace(login))
	cacheKey := usersByLoginPrefix + login

	var cacheResp twitchapi.User
	err := t.cache.Get(ctx, cacheKey, &cacheResp)
	if err == nil {
		return &cacheResp, nil
	}
	if err != nil {
		if err == cache.ErrNilResult {
			return nil, ErrUserNotFound
		}
		if err != cache.ErrNil {
			log.Println("failed to get twitch user from redis: " + err.Error())
		}
	}

	users, err := t.twitchAPI.GetUsers(ctx, &twitchapi.GetUsersRequest{Logins: []string{login}})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		err = t.cache.Set(ctx, cacheKey, nil, cache.OneHour)
		if err != nil {
			log.Println("failed to set nil twitch user in redis: " + err.Error())
		}
		return nil, ErrUserNotFound
	}

	user := users[0]
	_ = t.cache.Set(ctx, usersByLoginPrefix+user.Login, user, cache.OneDay)
	_ = t.cache.Set(ctx, usersByIDPrefix+user.ID, user, cache.TwoWeeks)
	return users[0], nil
}
