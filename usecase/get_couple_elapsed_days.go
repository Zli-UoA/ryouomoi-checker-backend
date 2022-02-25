package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"time"
)

type GetCoupleElapsedDaysUseCase interface {
	Execute(userID int64) (int, error)
}

type getCoupleElapsedDaysUseCaseImpl struct {
	ur repository.UserRepository
}

func (g *getCoupleElapsedDaysUseCaseImpl) Execute(userID int64) (int, error) {
	couple, err := g.ur.GetCurrentCouple(userID)
	if err != nil {
		return 0, err
	}
	now := time.Now()
	duration := now.Sub(couple.CreatedAt)
	days := duration.Hours() / 24
	return int(days), err
}

func NewGetCoupleElapsedDaysUseCase(ur repository.UserRepository) GetCoupleElapsedDaysUseCase {
	return &getCoupleElapsedDaysUseCaseImpl{ur: ur}
}
