package domain

import (
	"time"

	"golang.org/x/oauth2"
)

// DBConf is config using DB
type DBConf struct {
	Database string
	DSN      string
}

// ServerConf is above all
type ServerConf struct {
	DBConf DBConf
	Github oauth2.Config
}

// Login is auth login info
type Login struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"token_expiry"`
	URL    string    `json:"redirect_url"`
}

// Callback is callback param after github login
type Callback struct {
	User
	Code  string `json:"code"`
	State string `json:"state"`
}

// User is token for uesr
type User struct {
	Token string `json:"token"`
}

// Token is token expiry
type Token struct {
	Expiry time.Time `json:"expiry"`
}

// GithubToken is github token
type GithubToken struct {
	oauth2.Token
}
