package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
)

type GetLoggedInUserUseCase interface {
	Execute(userID int64) (*model.TwitterUser, error)
}

type getLoggedInUserUseCaseImpl struct {
	ur repository.UserRepository
	ts service.TwitterService
}

func (g *getLoggedInUserUseCaseImpl) Execute(userID int64) (*model.TwitterUser, error) {
	user, err := g.ur.GetUser(userID)
	if err != nil {
		return nil, err
	}
	twitterUser, err := g.ts.GetUser(user.TwitterAccessToken)
	if err != nil {
		return nil, err
	}
	return twitterUser, err
}

func NewGetLoggedInUserUseCase(ur repository.UserRepository, ts service.TwitterService) GetLoggedInUserUseCase {
	return &getLoggedInUserUseCaseImpl{
		ur: ur,
		ts: ts,
	}
}
