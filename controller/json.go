package controller

import "github.com/Zli-UoA/ryouomoi-checker-backend/model"

type TwitterLoginUrlJson struct {
	LoginUrl string `json:"loginUrl"`
}

type TwitterUser struct {
	ID          int64  `json:"id"`
	ScreenName  string `json:"screenName"`
	DisplayName string `json:"displayName"`
	ImageUrl    string `json:"imageUrl"`
	Biography   string `json:"biography"`
}

func convertToJson(twitterUser *model.TwitterUser) *TwitterUser {
	return &TwitterUser{
		ID:          twitterUser.ID,
		ScreenName:  twitterUser.ScreenName,
		DisplayName: twitterUser.DisplayName,
		ImageUrl:    twitterUser.ProfileImageUrl,
		Biography:   twitterUser.Biography,
	}
}

type Lover struct {
	User        *TwitterUser `json:"user"`
	TalkRoomUrl string       `json:"talkRoomUrl"`
}

type LovePoint struct {
	LovePoint int `json:"lovePoint"`
}

type UserLovePoint struct {
	LoverUser *TwitterUser `json:"user"`
	LovePoint int          `json:"lovePoint"`
}

type MatchResult struct {
	MatchSuccess bool   `json:"matchSuccess"`
	Lover        *Lover `json:"lover"`
}

type BrokeReport struct {
	ReasonID   int  `json:"reasonId"`
	AllowShare bool `json:"allowShare"`
}
