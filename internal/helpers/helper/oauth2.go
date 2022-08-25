package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

func GenerateOauthConfig() *oauth2.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found", err)
	}
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	return &oauth2.Config{
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveReadonlyScope},
		RedirectURL:  "https://ahorrojs.io/sync/callback",
	}
}
