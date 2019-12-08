package controllers

import (
	"backend/domain"
	"backend/interfaces/database"
	"backend/usecases"
	"net/http"

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

// GetUser returns github login url
func (oc *OauthController) GetUser(c Context) {
	g := domain.GithubToken{}
	c.Bind(&g)
	_, err := oc.Interactor.GetUserFromGithub(g)
	if err != nil {
		zap.S().Errorw(err.Error())
		c.SecureJSON(http.StatusInternalServerError, NewError(err))
	}
	c.Header("Location", logininfo.URL)
	c.SecureJSON(http.StatusTemporaryRedirect, logininfo)
}
