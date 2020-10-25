package twitchapi

import (
	"context"
	"net/http"
	"net/url"
)

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
}

type usersResponse struct {
	Data []*User `json:"data"`
}

type GetUsersRequest struct {
	IDs    []string
	Logins []string
}

func (t *twitchAPI) GetUsers(ctx context.Context, req *GetUsersRequest) ([]*User, error) {
	if len(req.IDs)+len(req.Logins) > 100 {
		return nil, ErrInvalidRequest
	}

	values := url.Values{}

	for _, id := range req.IDs {
		values.Add("id", id)
	}

	for _, login := range req.Logins {
		values.Add("login", login)
	}

	var res usersResponse
	if err := t.request(ctx, http.MethodGet, "https://api.twitch.tv/helix/users", values, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
