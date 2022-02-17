package model

import "github.com/mrjones/oauth"

type User struct {
	ID                 int64
	ScreenName         string
	DisplayName        string
	ProfileImageUrl    string
	Biography          string
	TwitterAccessToken *oauth.AccessToken
}

type TwitterUser struct {
	ID              int64
	ScreenName      string
	DisplayName     string
	ProfileImageUrl string
	Biography       string
}
