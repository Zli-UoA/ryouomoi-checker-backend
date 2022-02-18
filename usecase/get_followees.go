package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
)

type GetFolloweesUseCase interface {
	Execute(userID int64) ([]*model.TwitterUser, error)
}

type getFolloweesUseCaseImpl struct {
	ts service.TwitterService
	ur repository.UserRepository
}

func (g *getFolloweesUseCaseImpl) Execute(userID int64)([]*model.TwitterUser, error) {
	user, err := g.ur.GetUser(userID)
	if err != nil {
		return nil, err
	}
	followees, err := g.ts.GetFollowees(user.TwitterAccessToken)
	if err != nil {
		return nil, err
	}
	return followees, err
}

func NewGetFolloweesUseCase(ts service.TwitterService, ur repository.UserRepository) GetFolloweesUseCase {
	return &getFolloweesUseCaseImpl{
		ts: ts,
		ur: ur,
	}
}