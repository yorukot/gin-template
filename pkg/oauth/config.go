package config

import (
	"fmt"
	"os"

	"github.com/yorukot/go-template/pkg/logger"
	"github.com/yorukot/go-template/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
)

// Init oauth for goth
func init() {
	err := godotenv.Load("template.env")
	if err != nil {
		logger.Log.Fatal("Error loading .env file")
	}

	// Change your url
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			fmt.Sprintf("%s/auth/oauth/google/callback", utils.BackendURL),
			"email",
			"profile",
		),

		github.New(
			os.Getenv("GITHUB_CLIENT_ID"),
			os.Getenv("GITHUB_CLIENT_SECRET"),
			fmt.Sprintf("%s/auth/oauth/github/callback", utils.BackendURL),
		),

		gitlab.New(
			os.Getenv("GITLAB_CLIENT_ID"),
			os.Getenv("GITLAB_CLIENT_SECRET"),
			fmt.Sprintf("%s/auth/oauth/gitlab/callback", utils.BackendURL),
		),
	)
}
