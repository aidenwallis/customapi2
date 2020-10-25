package twitchapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aidenwallis/customapi2/pkg/config"
	"github.com/pkg/errors"
)

type TwitchAPI interface {
	GetUsers(ctx context.Context, req *GetUsersRequest) ([]*User, error)
	Start()
	Stop()
}

var ErrInvalidRequest = errors.New("twitchapi: invalid request")

type twitchAPI struct {
	bearerToken  string
	bearerTicker *time.Ticker
	clientID     string
	clientSecret string
	client       *http.Client
	stopCh       chan bool
}

func New(cfg *config.TwitchConfig) TwitchAPI {
	return &twitchAPI{
		bearerTicker: time.NewTicker(time.Hour * 24),
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		client:       &http.Client{Timeout: time.Second * 15},
		stopCh:       make(chan bool, 1),
	}
}

func (t *twitchAPI) Start() {
	t.startBearerLoop()
}

func (t *twitchAPI) Stop() {
	t.stopCh <- true
}

func (t *twitchAPI) prepareReq(ctx context.Context, method, uri string, values url.Values) (*http.Request, error) {
	encodedValues := values.Encode()
	if len(encodedValues) > 0 {
		uri += "?" + encodedValues
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+t.bearerToken)
	req.Header.Add("Client-ID", t.clientID)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (t *twitchAPI) request(ctx context.Context, method, uri string, values url.Values, out interface{}) error {
	req, err := t.prepareReq(ctx, method, uri, values)
	if err != nil {
		return err
	}

	res, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= 300 {
		return fmt.Errorf("Non 2xx http status code returned: %d", res.StatusCode)
	}

	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return errors.Wrap(err, "failed to decode response")
		}
	}

	return nil
}
