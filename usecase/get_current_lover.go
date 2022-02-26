package usecase

import (
	"errors"
	"fmt"
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"time"
)

type GetCurrentLoverUsecase interface {
	CheckBrokeCouple(userID int64) (*model.BrokeReport, error)
	Execute(userID int64) (*model.Lover, error)
}

type getCurrentLoverUsecaseImpl struct {
	ur repository.UserRepository
}

type BrokenCoupleNotExpiredError struct {
	RemainDays int
}

func (b *BrokenCoupleNotExpiredError) Error() string {
	return fmt.Sprintf("破局期間残り%v日", b.RemainDays)
}

var (
	BrokenCoupleError = errors.New("NoRows_latestBrokenCouple")
)

type BrokeReportNotFoundError struct {
	Lover *model.User
}

func (c *BrokeReportNotFoundError) Error() string {
	return fmt.Sprintf("%vにフラれました。破局手続きを行ってください。", c.Lover.DisplayName)
}

func (g *getCurrentLoverUsecaseImpl) CheckBrokeCouple(userID int64) (*model.BrokeReport, error) {
	couple, err := g.ur.GetLatestBrokenCouple(userID)
	if err != nil {
		return nil, BrokenCoupleError
	}
	brokeReport, err := g.ur.GetBrokeReport(userID, couple.ID)
	if err != nil {
		var lover *model.User
		if couple.User1.ID == userID {
			lover = couple.User2
		} else {
			lover = couple.User1
		}
		return nil, &BrokeReportNotFoundError{Lover: lover}
	}
	now := time.Now()
	expiredAt := couple.BrokenAt.AddDate(0, 1, 0)
	if now.Before(expiredAt) {
		remainDuration := expiredAt.Sub(now)
		remainDays := remainDuration.Hours() / 24
		return nil, &BrokenCoupleNotExpiredError{
			RemainDays: int(remainDays),
		}
	}
	return brokeReport, nil
}

func (g *getCurrentLoverUsecaseImpl) getTalkRoomUrl(loverUser *model.User) string {
	return createTwitterDMLink(loverUser.ID)
}

func (g *getCurrentLoverUsecaseImpl) Execute(userID int64) (*model.Lover, error) {
	couple, err := g.ur.GetCurrentCouple(userID)
	if err != nil {
		_, err := g.CheckBrokeCouple(userID)
		return nil, err
	}
	var loverUser *model.User
	if couple.User1.ID == userID {
		loverUser = couple.User2
	} else {
		loverUser = couple.User1
	}
	lover := &model.Lover{
		User: &model.TwitterUser{
			ID:              loverUser.ID,
			ScreenName:      loverUser.ScreenName,
			DisplayName:     loverUser.DisplayName,
			ProfileImageUrl: loverUser.ProfileImageUrl,
			Biography:       loverUser.Biography,
		},
		TalkRoomUrl: createTwitterDMLink(loverUser.ID),
	}
	return lover, nil
}

func NewGetCurrentLover(ur repository.UserRepository) GetCurrentLoverUsecase {
	return &getCurrentLoverUsecaseImpl{
		ur: ur,
	}
}
