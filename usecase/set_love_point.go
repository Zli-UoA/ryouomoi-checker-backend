package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"log"
	"strconv"
)

type SetLovePointUseCase interface {
	Execute(userID, loverUserID int64, lovePoint int) (bool, error)
}

type setLovePointUseCaseImpl struct {
	thresholdLovePoint int
	botUserID          int64
	ur                 repository.UserRepository
	tc                 service.TwitterService // bot用のTwitterServiceを受け取る(AccessTokenの権限が違うため)
}

func createTwitterDMLink(userID int64) string {
	return "https://twitter.com/messages/compose?recipient_id=" + strconv.FormatInt(userID, 10)
}

func createMessage(loverName, talkRoomUrl string) string {
	return loverName + "さんと両思いになりました。トークルームでお話しましょう。\n" + talkRoomUrl
}

func (s *setLovePointUseCaseImpl) Execute(userID, loverUserID int64, lovePoint int) (bool, error) {
	userLovePoint := &model.UserLovePoint{
		ID:          0,
		UserID:      userID,
		LoverUserID: loverUserID,
		LovePoint:   lovePoint,
	}
	_, err := s.ur.SetLovePoint(userLovePoint)
	if err != nil {
		return false, err
	}
	loverUserLovePoint, err := s.ur.GetLovePoint(loverUserID, userID)
	if err != nil {
		return false, nil
	}
	if lovePoint+loverUserLovePoint.LovePoint < s.thresholdLovePoint {
		return false, nil
	}
	couple, err := s.ur.CreateCouple(userID, loverUserID)
	if err != nil {
		return false, err
	}
	botUser, err := s.ur.GetUser(s.botUserID)
	if err != nil {
		return false, err
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
	return true, nil
}

func NewSetLovePointUseCase(thresholdLovePoint int, botUserID int64, ur repository.UserRepository, tc service.TwitterService) SetLovePointUseCase {
	return &setLovePointUseCaseImpl{
		thresholdLovePoint: thresholdLovePoint,
		botUserID:          botUserID,
		ur:                 ur,
		tc:                 tc,
	}
}
