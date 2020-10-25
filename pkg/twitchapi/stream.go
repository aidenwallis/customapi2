package twitchapi

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type Stream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

type streamsResponse struct {
	Data []*Stream `json:"data"`
}

type GetStreamsRequest struct {
	UserIDs    []string
	UserLogins []string
}

func (t *twitchAPI) GetStreams(ctx context.Context, req *GetStreamsRequest) ([]*Stream, error) {
	if len(req.UserIDs)+len(req.UserLogins) > 100 {
		return nil, ErrInvalidRequest
	}

	values := url.Values{}

	for _, id := range req.UserIDs {
		values.Add("user_id", id)
	}

	for _, login := range req.UserLogins {
		values.Add("user_login", login)
	}

	var body streamsResponse
	if err := t.request(ctx, http.MethodGet, "https://api.twitch.tv/helix/streams", values, &body); err != nil {
		return nil, err
	}

	return body.Data, nil
}
