package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"log"
	"time"
)

type DeleteCurrentLoverUseCase interface {
	Execute(userID int64, brokeReason int, allowShare bool) error
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

// 最初に破局手続きをするほう
func (d *deleteCurrentLoverUseCaseImpl) BreakFirst(userID int64, couple *model.Couple, brokeReason int, allowShare bool) error {
	var user *model.User
	var lover *model.User
	if couple.User1.ID == userID {
		user = couple.User1
		lover = couple.User2
	} else {
		user = couple.User2
		lover = couple.User1
	}
	now := time.Now()
	couple.BrokenAt = &now
	_, err := d.ur.UpdateCouple(couple)
	if err != nil {
		return err
	}
	brokeReport := &model.BrokeReport{
		Couple:        couple,
		User:          user,
		BrokeReasonID: brokeReason,
		AllowShare:    allowShare,
	}
	_, err = d.ur.CreateBrokeReport(brokeReport)
	if err != nil {
		return err
	}
	err = d.ur.DeleteLovePoint(userID, lover.ID)
	if err != nil {
		return err
	}
	return nil
}

var ErrorBrokeReportAlreadyExists = errors.New("user broke report already exists")

// あとに破局手続きをするほう
func (d *deleteCurrentLoverUseCaseImpl) BreakSecond(userID int64, couple *model.Couple, brokeReason int, allowShare bool) error {
	_, err := d.ur.GetBrokeReport(userID, couple.ID)
	if err == nil {
		return ErrorBrokeReportAlreadyExists
	}
	var user *model.User
	var lover *model.User
	if couple.User1.ID == userID {
		user = couple.User1
		lover = couple.User2
	} else {
		user = couple.User2
		lover = couple.User1
	}
	brokeReport := &model.BrokeReport{
		Couple:        couple,
		User:          user,
		BrokeReasonID: brokeReason,
		AllowShare:    allowShare,
	}
	_, err = d.ur.CreateBrokeReport(brokeReport)
	if err != nil {
		return err
	}
	err = d.ur.DeleteLovePoint(userID, lover.ID)
	if err != nil {
		return err
	}
	existingBrokeReport, err := d.ur.GetBrokeReport(lover.ID, couple.ID)
	if err != nil {
		log.Println("broke report is nil")
		return err
	}
	botUser, err := d.ur.GetUser(d.botUserID)
	if err != nil {
		log.Println("bot user is nil")
		return nil
	}
	if allowShare && existingBrokeReport.AllowShare {
		err = d.ts.SendTweet(botUser.TwitterAccessToken, createTweetContent(couple))
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (d *deleteCurrentLoverUseCaseImpl) Execute(userID int64, brokeReason int, allowShare bool) error {
	couple, err := d.ur.GetCurrentCouple(userID)
	if err == nil {
		err = d.BreakFirst(userID, couple, brokeReason, allowShare)
	} else {
		couple, err = d.ur.GetLatestBrokenCouple(userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("latest broken couple is nil. %w", err)
			} else {
				return err
			}
		}
		err = d.BreakSecond(userID, couple, brokeReason, allowShare)
	}
	if err != nil {
		return err
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
