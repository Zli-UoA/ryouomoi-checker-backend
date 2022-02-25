package usecase

import (
	"errors"
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"log"
	"strconv"
	"time"
)

type SetLovePointUseCase interface {
	Execute(userID, loverUserID int64, lovePoint int) (*model.Lover, error)
}

type setLovePointUseCaseImpl struct {
	thresholdLovePoint int
	botUserID          int64
	ur                 repository.UserRepository
	cr                 repository.ChatRepository
	tc                 service.TwitterService // bot用のTwitterServiceを受け取る(AccessTokenの権限が違うため)
}

var (
	CoupleAlreadyExistError = errors.New("couple already exist")
)

func createTwitterDMLink(userID int64) string {
	return "https://twitter.com/messages/compose?recipient_id=" + strconv.FormatInt(userID, 10)
}

func createMessage(loverName, talkRoomUrl string) string {
	return loverName + "さんと両思いになりました。トークルームでお話しましょう。\n" + talkRoomUrl
}

func (s *setLovePointUseCaseImpl) Execute(userID, loverUserID int64, lovePoint int) (*model.Lover, error) {
	_, err := s.ur.GetCurrentCouple(userID)
	if err == nil {
		return nil, CoupleAlreadyExistError
	}
	brokenCouple, err := s.ur.GetLatestBrokenCouple(userID)
	if err == nil {
		now := time.Now()
		expireAt := brokenCouple.BrokenAt.AddDate(0, 1, 0)
		if now.Before(expireAt) {
			remainDuration := expireAt.Sub(now)
			remainDays := remainDuration.Hours() / 24
			return nil, &BrokenCoupleNotExpiredError{RemainDays: int(remainDays)}
		}
	}
	userLovePoint := &model.UserLovePoint{
		ID:          0,
		UserID:      userID,
		LoverUserID: loverUserID,
		LovePoint:   lovePoint,
	}
	user, err := s.ur.GetUser(userID)
	if err != nil {
		return nil, err
	}
	_, err = s.ur.SetLovePoint(userLovePoint)
	if err != nil {
		return nil, err
	}
	loverUser, err := s.ur.GetUser(loverUserID)
	if err != nil {
		return nil, nil
	}
	loverUserLovePoint, err := s.ur.GetLovePoint(loverUserID, userID)
	if err != nil {
		return nil, nil
	}
	if lovePoint+loverUserLovePoint.LovePoint < s.thresholdLovePoint {
		return nil, nil
	}
	couple := &model.Couple{
		User1: user,
		User2: loverUser,
	}
	_, err = s.ur.CreateCouple(couple)
	if err != nil {
		return nil, err
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
	_, err = s.cr.CreateChatRoom(couple)
	if err != nil {
		return lover, err
	}
	botUser, err := s.ur.GetUser(s.botUserID)
	if err != nil {
		return lover, err
	}
	user1 := couple.User1
	user2 := couple.User2
	err = s.tc.SendDirectMessage(botUser.TwitterAccessToken, user1.ID, createMessage(user2.DisplayName, createTwitterDMLink(user2.ID)))
	if err != nil {
		log.Println(err)
	}
	err = s.tc.SendDirectMessage(botUser.TwitterAccessToken, user2.ID, createMessage(user1.DisplayName, createTwitterDMLink(user1.ID)))
	if err != nil {
		log.Println(err)
	}
	err = s.tc.SendDirectMessage(user1.TwitterAccessToken, user2.ID, "好きです")
	if err != nil {
		log.Println(err)
	}
	err = s.tc.SendDirectMessage(user2.TwitterAccessToken, user1.ID, "好きです")
	if err != nil {
		log.Println(err)
	}
	return lover, nil
}

func NewSetLovePointUseCase(
	thresholdLovePoint int,
	botUserID int64,
	ur repository.UserRepository,
	cr repository.ChatRepository,
	tc service.TwitterService,
) SetLovePointUseCase {
	return &setLovePointUseCaseImpl{
		thresholdLovePoint: thresholdLovePoint,
		botUserID:          botUserID,
		ur:                 ur,
		cr:                 cr,
		tc:                 tc,
	}
}
