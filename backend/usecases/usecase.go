package usecases

import (
	"backend/domain"
	"errors"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// OauthRepository is ...
type OauthRepository interface {
	StoreState(string, string, time.Time) error
	FindBySessionID(string) (string, error)
	FindBySessionIDAndUserToken(string, string) (time.Time, int, error)
	FindByUserTokenID(int) (string, string, string, time.Time, error)
	StoreUserToken(string, string, time.Time, int) (int, error)
	StoreGithubToken(string, string, string, time.Time) (int, error)
}

// OauthInteractor is ...
type OauthInteractor struct {
	OauthRepository OauthRepository
}

// SetupGithubLogin is ...
func (oi *OauthInteractor) SetupGithubLogin(conf domain.ServerConf, s domain.Session) (login domain.Login, err error) {
	login.State = createRand()
	login.URL = conf.Github.AuthCodeURL(login.State)
	expiry := time.Now().Add(10 * time.Minute)
	err = oi.OauthRepository.StoreState(login.State, s.ID, expiry)
	return
}

const (
	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	indexBit  = 6
	indexMask = 1<<indexBit - 1
	indexMax  = 63 / indexBit
)

func createRand() (randVal string) {
	randSource := rand.NewSource(time.Now().UnixNano())
	n := 32
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

// RegistToken is ...
func (oi *OauthInteractor) RegistToken(ctx *gin.Context, conf domain.ServerConf, c domain.Callback) (user domain.User, err error) {
	// recieved state is expected or not
	state, err := oi.OauthRepository.FindBySessionID(c.ID)
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
	id, err := oi.OauthRepository.StoreGithubToken(githubToken.AccessToken, githubToken.TokenType, githubToken.RefreshToken, githubToken.Expiry)
	if err != nil {
		return
	}

	// make user token
	token := createRand()
	expiry := time.Now().Add(7 * 24 * time.Hour)
	count, err := oi.OauthRepository.StoreUserToken(c.ID, token, expiry, id)
	if err != nil {
		return
	}
	if count == 0 {
		err = errors.New("no user session")
		return
	}
	user.Token = token
	return
}

// AuthToken is ...
func (oi *OauthInteractor) AuthToken(a domain.Auth) (token domain.GithubToken, err error) {
	expiry, id, err := oi.OauthRepository.FindBySessionIDAndUserToken(a.ID, a.Token)
	if expiry.After(time.Now()) {
		err = errors.New("user token expiry")
		return
	}
	accessToken, tktype, refToken, expiry, err := oi.OauthRepository.FindByUserTokenID(id)
	if err != nil {
		return
	}
	token.AccessToken = accessToken
	token.TokenType = tktype
	token.RefreshToken = refToken
	token.Expiry = expiry
	return
}
