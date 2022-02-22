package usecase

import (
	"errors"
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
)

type GetCurrentLoverUsecase interface {
	CheckBreakFirst(userID int64) (*model.BrokeReport, error)
	Execute(userID int64)(*model.TwitterUser, error)
}

type getCurrentLoverUsecaseImpl struct {
	ur repository.UserRepository
}

var (
	BrokenCoupleError = errors.New("NoRows_latestBrokenCouple")
	BrokenReportError = errors.New("NoRows_BrokenReport")
)

func (g *getCurrentLoverUsecaseImpl)CheckBreakFirst(userID int64) (*model.BrokeReport, error) {
	couple, err := g.ur.GetLatestBrokenCouple(userID)
	if err != nil {
		return nil, BrokenCoupleError
	}
	brokeReport, err := g.ur.GetBrokeReport(userID, couple.ID)
	if err != nil {
		return nil, BrokenReportError
	}
	return brokeReport, nil

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

