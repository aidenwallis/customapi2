package twitchapi

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type bearerTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (t *twitchAPI) getBearerToken() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	values := url.Values{}
	values.Add("client_id", t.clientID)
	values.Add("client_secret", t.clientSecret)
	values.Add("grant_type", "client_credentials")

	var body bearerTokenResponse
	if err := t.request(ctx, http.MethodPost, "https://id.twitch.tv/oauth2/token", values, &body); err != nil {
		return err
	}

	if body.AccessToken == "" {
		return errors.New("access_token not in response")
	}

	log.Println("twitchapi: generated bearer token")
	t.bearerToken = body.AccessToken
	return nil
}

func (t *twitchAPI) startBearerLoop() {
	for {
		if err := t.getBearerToken(); err != nil {
			log.Println("Failed to get bearer token: " + err.Error())
			time.Sleep(time.Second * 5)
			continue
		}

		select {
		case <-t.stopCh:
			return

		// give bearer tokens a 24 hour lifetime
		case <-time.After(time.Hour * 24):
			continue
		}
	}
}
