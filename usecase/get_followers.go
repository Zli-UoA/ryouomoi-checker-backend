package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
)

type GetFollowersUseCase interface {
	Execute(userID int64) ([]*model.TwitterUser, error)
}

type getFollowersUseCaseImpl struct {
	ts service.TwitterService
	ur repository.UserRepository
}

func (g *getFollowersUseCaseImpl) Execute(userID int64)([]*model.TwitterUser, error) {
	user, err := g.ur.GetUser(userID)
	if err != nil {
		return nil, err
	}
	followers, err := g.ts.GetFollowers(user.TwitterAccessToken)
	if err != nil {
		return nil, err
	}
	return followers, err
}

func NewGetFollowersUseCase(ts service.TwitterService, ur repository.UserRepository) GetFollowersUseCase {
	return &getFollowersUseCaseImpl{
		ts: ts,
		ur: ur,
	}
}