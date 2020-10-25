package twitch

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/aidenwallis/customapi2/pkg/data/twitch"
	"github.com/aidenwallis/customapi2/pkg/responder"
)

func (f *twitchFeature) getUserID(ctx context.Context, w responder.Responder, req *http.Request) {
	login := chi.URLParam(req, "login")
	if login == "" {
		w.BadRequest("Invalid login.", nil)
		return
	}

	if !loginRegex.MatchString(login) {
		w.BadRequest("Invalid login.", nil)
		return
	}

	user, err := f.Data.Twitch.GetUserByLogin(ctx, login)
	if err != nil {
		if err == twitch.ErrUserNotFound {
			w.BadRequest("User not found.", nil)
			return
		}
		log.Println("Failed to get user from Twitch: " + err.Error())
		w.BadRequest("Twitch API returned an error.", nil)
		return
	}

	w.OK(user.ID, makeTwitchUser(user))
}
