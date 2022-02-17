package model

type User struct {
	ID                int64
	ScreenName        string
	DisplayName       string
	ProfileImageUrl   string
	Biography         string
	AccessToken       string
	AccessTokenSecret string
}

type TwitterUser struct {
	ID              int64
	ScreenName      string
	DisplayName     string
	ProfileImageUrl string
	Biography       string
}
