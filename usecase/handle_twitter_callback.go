package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/mrjones/oauth"
)

type HandleTwitterCallbackUseCase interface {
	Execute(oauthToken, oauthVerifier string) (string, error) // returns user auth jwt
}

type handleTwitterCallbackUseCaseImpl struct {
	ts  service.TwitterService
	ujs service.UserJWTService
	ur  repository.UserRepository
}

func convertToUser(twitterUser *model.TwitterUser, accessToken *oauth.AccessToken) *model.User {
	user := model.User{
		ID:                 twitterUser.ID,
		ScreenName:         twitterUser.ScreenName,
		DisplayName:        twitterUser.DisplayName,
		ProfileImageUrl:    twitterUser.ProfileImageUrl,
		Biography:          twitterUser.Biography,
		TwitterAccessToken: accessToken,
	}
	return &user
}

func (h *handleTwitterCallbackUseCaseImpl) Execute(oauthToken, oauthVerifier string) (string, error) {
	accessToken, err := h.ts.AuthorizeToken(oauthToken, oauthVerifier)
	if err != nil {
		return "", err
	}
	twitterUser, err := h.ts.GetUser(accessToken)
	if err != nil {
		return "", err
	}
	user := convertToUser(twitterUser, accessToken)
	_, err = h.ur.GetUser(twitterUser.ID)
	if err == nil {
		_, err = h.ur.UpdateUser(user)
		if err != nil {
			return "", err
		}
	} else {
		_, err = h.ur.CreateUser(user)
		if err != nil {
			return "", err
		}
	}
	jwt, err := h.ujs.CreateUserIDJWT(user.ID)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func NewHandleTwitterCallbackUseCase(ts service.TwitterService, ujs service.UserJWTService, ur repository.UserRepository) HandleTwitterCallbackUseCase {
	return &handleTwitterCallbackUseCaseImpl{
		ts:  ts,
		ujs: ujs,
		ur:  ur,
	}
}
