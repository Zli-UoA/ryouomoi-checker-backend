package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
)

type GetLovePointsUseCase interface {
	Execute(userID int64) ([]*model.UserLovePointWithTwitterUser, error)
}

type getLovePointsUseCaseImpl struct {
	ur repository.UserRepository
	ts service.TwitterService
}

func (g *getLovePointsUseCaseImpl) Execute(userID int64) ([]*model.UserLovePointWithTwitterUser, error) {
	user, err := g.ur.GetUser(userID)
	if err != nil {
		return nil, err
	}
	lovePoints, err := g.ur.GetLovePoints(userID)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(lovePoints))
	for i, lovePoint := range lovePoints {
		ids[i] = lovePoint.LoverUserID
	}
	users, err := g.ts.Lookup(user.TwitterAccessToken, ids)
	if err != nil {
		return nil, err
	}
	result := make([]*model.UserLovePointWithTwitterUser, len(lovePoints))
	for i, lovePoint := range lovePoints {
		result[i] = &model.UserLovePointWithTwitterUser{
			ID:        lovePoint.ID,
			UserID:    lovePoint.UserID,
			LoverUser: users[i],
			LovePoint: lovePoint.LovePoint,
		}
	}
	return result, nil
}

func NewGetLovePointsUseCase(ur repository.UserRepository, ts service.TwitterService) GetLovePointsUseCase {
	return &getLovePointsUseCaseImpl{
		ur: ur,
		ts: ts,
	}
}
