package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
)

type FriendsSearchUseCase interface {
	Execute(userID int64, query string) ([]*model.TwitterUser, error)
}

type friendsSearchUseCaseImpl struct {
	ts service.TwitterService
	ur repository.UserRepository
}

func (f *friendsSearchUseCaseImpl) Execute(userID int64, query string) ([]*model.TwitterUser, error) {
	user, err := f.ur.GetUser(userID)
	if err != nil {
		return nil, err
	}
	searchResult, err := f.ts.Search(user.TwitterAccessToken, query)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func NewFriendsSearchUseCase(ts service.TwitterService, ur repository.UserRepository) FriendsSearchUseCase {
	return &friendsSearchUseCaseImpl{
		ts: ts,
		ur: ur,
	}
}
