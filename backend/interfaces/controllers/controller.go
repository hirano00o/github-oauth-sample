package controllers

import (
	"backend/domain"
	"backend/interfaces/database"
	"backend/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OauthController is ...
type OauthController struct {
	Interactor usecases.OauthInteractor
}

// NewController is ...
func NewController(db database.DBHandler) *OauthController {
	return &OauthController{
		Interactor: usecases.OauthInteractor{
			OauthRepository: &database.OauthRepository{
				DBHandler: db,
			},
		},
	}
}

// Login returns github login url
func (oc *OauthController) Login(c Context, conf domain.ServerConf) {
	s := domain.Session{}
	c.Bind(&s)
	logininfo, err := oc.Interactor.SetupGithubLogin(conf, s)
	if err != nil {
		zap.S().Errorw(err.Error())
		c.SecureJSON(http.StatusInternalServerError, NewError(err))
	}
	c.Header("Location", logininfo.URL)
	c.SecureJSON(http.StatusTemporaryRedirect, logininfo)
}

// Callback authenticate user and return token
func (oc *OauthController) Callback(c *gin.Context, conf domain.ServerConf) {
	callback := domain.Callback{}
	c.Bind(&callback)
	user, err := oc.Interactor.RegistToken(c, conf, callback)
	if err != nil {
		zap.S().Errorw(err.Error())
		c.SecureJSON(http.StatusInternalServerError, NewError(err))
	}
	c.SecureJSON(http.StatusOK, user)
}

// Auth authenticate user and return github token
func (oc *OauthController) Auth(c Context) {
	a := domain.Auth{}
	c.Bind(&a)
	token, err := oc.Interactor.AuthToken(a)
	if err != nil {
		zap.S().Errorw(err.Error())
		c.SecureJSON(http.StatusInternalServerError, NewError(err))
	}
	c.SecureJSON(http.StatusOK, token)
}
