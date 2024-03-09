package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
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
	dir, _ := os.Getwd()

	err := godotenv.Load(path.Join(dir, "auth", ".env"))

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	id := os.Getenv("CLIENT_ID")
	secret := os.Getenv("CLIENT_SECRET")

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
