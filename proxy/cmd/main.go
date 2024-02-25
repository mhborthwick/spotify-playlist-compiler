package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mhborthwick/spotify-playlist-compiler/.config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var authConfig *oauth2.Config

func GetRandomString() string {
	return uuid.NewString()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hey")
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	state := GetRandomString()
	http.Redirect(w, r, authConfig.AuthCodeURL(state), http.StatusSeeOther)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, _ := authConfig.Exchange(r.Context(), code)
	fmt.Fprint(w, token.AccessToken)
}

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}

	authConfig = &oauth2.Config{
		ClientID:     cfg.Id,
		ClientSecret: cfg.Secret,
		RedirectURL:  "http://localhost:1337/callback",
		Endpoint:     spotify.Endpoint,
		Scopes: []string{
			"user-read-email",
			"user-read-private",
		},
	}

	http.HandleFunc("/", Handler)
	http.HandleFunc("/login", Authorize)
	http.HandleFunc("/callback", Callback)

	log.Fatal(http.ListenAndServe(":1337", nil))
}
