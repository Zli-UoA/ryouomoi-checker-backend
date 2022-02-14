package service

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/mrjones/oauth"
)

type TwitterService interface {
	GetLoginUrl() (string, error)
	AuthorizeToken(token, verificationCode string) (*oauth.AccessToken, error)
	GetUser(token *oauth.AccessToken) (*model.TwitterUser, error)
	GetFollowees(token *oauth.AccessToken) ([]*model.TwitterUser, error)
	GetFollowers(token *oauth.AccessToken) ([]*model.TwitterUser, error)
	Search(query string) ([]*model.TwitterUser, error)
}
