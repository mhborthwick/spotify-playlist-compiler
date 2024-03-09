package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var authConfig *oauth2.Config

func GetRandomString() string {
	return uuid.NewString()
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Placeholder")
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	state := GetRandomString()
	http.Redirect(w, r, authConfig.AuthCodeURL(state), http.StatusSeeOther)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		fmt.Println("Unauthorized")
		return
	}
	token, err := authConfig.Exchange(r.Context(), code)
	if err != nil {
		fmt.Println("Unauthorized")
		return
	}
	fmt.Fprint(w, token.AccessToken)
}

func main() {
	id := os.Getenv("CLIENT_ID")
	secret := os.Getenv("CLIENT_SECRET")

	if id == "" || secret == "" {
		fmt.Println("Missing Env Vars")
		return
	}

	authConfig = &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  "http://localhost:1337/callback",
		Endpoint:     spotify.Endpoint,
		Scopes: []string{
			"user-read-email",
			"user-read-private",
			"playlist-modify-public",
			"playlist-modify-private",
		},
	}

	http.HandleFunc("/", Index)
	http.HandleFunc("/login", Authorize)
	http.HandleFunc("/callback", Callback)

	log.Fatal(http.ListenAndServe(":1337", nil))
}
