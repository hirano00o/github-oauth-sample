package infrastructure

import (
	"backend/domain"
	"os"

	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
)

func NewConf() (conf domain.ServerConf) {
	conf.DBConf = getDBConf()
	conf.Github = getGithubConf()
	return
}

func getDBConf() (db domain.DBConf) {
	db.Database = os.Getenv("DATABASE")
	db.DSN = os.Getenv("DSN")
	return
}

func getGithubConf() (github oauth2.Config) {
	scopes := []string{"repo"}
	github = oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("SERVER_HOST"),
		Scopes:       scopes,
		Endpoint:     oauth2github.Endpoint,
	}
	return
}
