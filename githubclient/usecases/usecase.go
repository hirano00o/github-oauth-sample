package usecases

import (
	"backend/domain"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// OauthRepository is interface repositories
type OauthRepository interface {
	StoreState(string, string, time.Time) error
	FindStateByUserToken(string) (int, string, error)
	FindGithubTokenIDByUserToken(string) (time.Time, int, error)
	FindByUserID(int) (string, string, string, time.Time, error)
	UpdateUserTokenExpiry(string, time.Time) (int, error)
	StoreGithubToken(string, string, string, time.Time, int) error
}

// OauthInteractor is usecase interactor
type OauthInteractor struct {
	OauthRepository OauthRepository
}

// CreateGithubLoginURL creates github login url
func (oi *OauthInteractor) CreateGithubLoginURL(conf domain.ServerConf) (login domain.Login, err error) {
	login.Token = createRand(32)
	state := createRand(32)
	login.URL = conf.Github.AuthCodeURL(state)
	login.Expiry = time.Now().Add(10 * time.Minute)
	err = oi.OauthRepository.StoreState(login.Token, state, login.Expiry)
	return
}

const (
	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	indexBit  = 6
	indexMask = 1<<indexBit - 1
	indexMax  = 63 / indexBit
)

func createRand(n int) (randVal string) {
	randSource := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	cache, remain := randSource.Int63(), indexMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSource.Int63(), indexMax
		}
		index := int(cache & indexMask)
		if index < len(letters) {
			b[i] = letters[index]
			i--
		}
		cache >>= indexBit
		remain--
	}
	randVal = string(b)
	return
}

// GetUserFromGithub gets username and avater url
func (oi *OauthInteractor) GetUserFromGithub(g domain.GithubToken) (err error) {
}

// RegistGithubToken creates and registration github token
func (oi *OauthInteractor) RegistGithubToken(ctx *gin.Context, conf domain.ServerConf, c domain.Callback) (user domain.Token, err error) {
	// recieved state is expected or not
	id, state, err := oi.OauthRepository.FindStateByUserToken(c.Token)
	if err != nil {
		return
	}
	if state != c.State {
		err = errors.New("not match state")
		return
	}
	// make github token
	githubToken, err := conf.Github.Exchange(ctx, c.Code)
	if err != nil {
		return
	}
	err = oi.OauthRepository.StoreGithubToken(githubToken.AccessToken, githubToken.TokenType, githubToken.RefreshToken, githubToken.Expiry, id)
	if err != nil {
		return
	}

	// update user token expiry
	expiry := time.Now().Add(7 * 24 * time.Hour)
	count, err := oi.OauthRepository.UpdateUserTokenExpiry(c.Token, expiry)
	if err != nil {
		return
	}
	if count == 0 {
		err = errors.New("no user token")
		return
	}
	user.Expiry = expiry
	return
}

func checkBearer(token string) (err error) {
	if strings.HasPrefix(token, "Bearer") {
		err = errors.New("Token is not verified.")
		return
	}
	return
}

// AuthToken authenticate user token and returns github token
func (oi *OauthInteractor) AuthToken(a domain.User) (token domain.GithubToken, err error) {
	err = checkBearer(a.Token)
	if err != nil {
		return
	}
	expiry, id, err := oi.OauthRepository.FindGithubTokenIDByUserToken(a.Token)
	if expiry.After(time.Now()) {
		err = errors.New("user token expiry")
		return
	}
	accessToken, tktype, refToken, expiry, err := oi.OauthRepository.FindByUserID(id)
	if err != nil {
		return
	}
	token.AccessToken = accessToken
	token.TokenType = tktype
	token.RefreshToken = refToken
	token.Expiry = expiry
	return
}
