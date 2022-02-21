package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
)

type GetCurrentLoverUsecase interface {
	Execute(userID int64)(*model.TwitterUser, error)
}

type getCurrentLoverUsecaseImpl struct {
	ur repository.UserRepository
}

func (g *getCurrentLoverUsecaseImpl) Execute(userID int64) (*model.TwitterUser, error) {
	couple, err := g.ur.GetCurrentCouple(userID)
	if err != nil {
		return nil, err
	}
	var lover *model.User
	if couple.User1.ID == userID {
		lover = couple.User2
	} else {
		lover = couple.User1
	}
	return &model.TwitterUser{
		ID:              lover.ID,
		ScreenName:      lover.ScreenName,
		DisplayName:     lover.DisplayName,
		ProfileImageUrl: lover.ProfileImageUrl,
		Biography:       lover.Biography,
	}, nil
}

func NewGetCurrentLover(ur repository.UserRepository) GetCurrentLoverUsecase {
	return &getCurrentLoverUsecaseImpl{
		ur: ur,
	}
}

