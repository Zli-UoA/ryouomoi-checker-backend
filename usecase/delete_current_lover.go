package usecase

import (
	"fmt"
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"log"
	"time"
)

type DeleteCurrentLoverUseCase interface {
	Execute(userID int64, brokeReport *model.BrokeReport) error
}

type deleteCurrentLoverUseCaseImpl struct {
	botUserID int64
	ur        repository.UserRepository
	ts        service.TwitterService // bot
}

func createTweetContent(couple *model.Couple) string {
	duration := couple.BrokenAt.Sub(couple.CreatedAt)
	months := int(duration.Hours()) / 24 / 30
	var durationStr string
	if months > 0 {
		durationStr = fmt.Sprintf("%dヶ月間", months)
	} else {
		days := int(duration.Hours()) / 24
		durationStr = fmt.Sprintf("%d日間", days)
	}
	content := fmt.Sprintf(
		"速報！\n%s付き合った%sさん(@%s)と%sさん(@%s)が破局しました。",
		durationStr,
		couple.User1.DisplayName,
		couple.User1.ScreenName,
		couple.User2.DisplayName,
		couple.User2.ScreenName,
	)
	return content
}

func (d *deleteCurrentLoverUseCaseImpl) Execute(userID int64, brokeReport *model.BrokeReport) error {
	couple, err := d.ur.GetCurrentCouple(userID)
	if err != nil {
		return err
	}
	now := time.Now()
	couple.BrokenAt = &now
	_, err = d.ur.UpdateCouple(couple)
	_, err = d.ur.CreateBrokeReport(brokeReport)
	var lover *model.User
	if couple.User1.ID == userID {
		lover = couple.User2
	} else {
		lover = couple.User1
	}
	err = d.ur.DeleteLovePoint(userID, lover.ID)
	if err != nil {
		return err
	}
	existingBrokeReport, err := d.ur.GetBrokeReport(lover.ID, couple.ID)
	if err != nil {
		return nil
	}
	botUser, err := d.ur.GetUser(d.botUserID)
	if err != nil {
		log.Println(err)
		return nil
	}
	if brokeReport.AllowShare && existingBrokeReport.AllowShare {
		err = d.ts.SendTweet(botUser.TwitterAccessToken, createTweetContent(couple))
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func NewDeleteCurrentLover(botUseID int64, ur repository.UserRepository, ts service.TwitterService) DeleteCurrentLoverUseCase {
	return &deleteCurrentLoverUseCaseImpl{
		botUserID: botUseID,
		ur:        ur,
		ts:        ts,
	}
}
