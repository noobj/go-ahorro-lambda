package helper

import (
	"github.com/noobj/go-serverless-services/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

func GenerateOauthConfig() *oauth2.Config {
	env := config.GetInstance()
	googleClientId := env.GoogleClientId
	googleClientSecret := env.GoogleClientSecret

	return &oauth2.Config{
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveReadonlyScope},
		RedirectURL:  "https://ahorrojs.io/sync/callback",
	}
}
