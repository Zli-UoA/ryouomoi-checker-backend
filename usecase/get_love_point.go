package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
)

type GetLovePointUsecase interface {
	Execute(userID, loverUserID int64) (*model.UserLovePoint, error)
}

type getLovePointUseCaseImpl struct {
	ur repository.UserRepository
}

func (g *getLovePointUseCaseImpl) Execute(userID, loverID int64) (*model.UserLovePoint, error) {
	lover, err := g.ur.GetLovePoint(userID, loverID)
	if err != nil {
		return nil, err
	}
	return lover, nil
}

func NewGetLovePointUsecase(ur repository.UserRepository) GetLovePointUsecase {
	return &getLovePointUseCaseImpl{
		ur: ur,
	}
}
