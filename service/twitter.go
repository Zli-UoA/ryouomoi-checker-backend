package service

import (
	"encoding/json"
	"errors"
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/mrjones/oauth"
	"log"
	"strconv"
)

type TwitterService interface {
	GetLoginUrl() (string, error)
	AuthorizeToken(token, verificationCode string) (*oauth.AccessToken, error)
	GetUser(token *oauth.AccessToken) (*model.TwitterUser, error)
	GetFollowees(token *oauth.AccessToken) ([]*model.TwitterUser, error)
	GetFollowers(token *oauth.AccessToken) ([]*model.TwitterUser, error)
	Search(token *oauth.AccessToken, query string) ([]*model.TwitterUser, error)
}

type userObject struct {
	ID              int64  `json:"id"`
	IDStr           string `json:"id_str"`
	Name            string `json:"name"`
	ScreenName      string `json:"screen_name"`
	ProfileImageURL string `json:"profile_image_url_https"`
	Description     string `json:"description"`
}

type usersObject struct {
	Users []userObject `json:"users"`
	NextCursor int `json:"next_cursor"`
}

type twitterServiceImpl struct {
	consumer *oauth.Consumer
	callbackUrl string
	requestTokenMap map[string]*oauth.RequestToken
}

const (
	requestTokenURL = "https://api.twitter.com/oauth/request_token"
	authorizationURL = "https://api.twitter.com/oauth/authenticate"
	accessTokenURL = "https://api.twitter.com/oauth/access_token"
	twitterAPIEndpoint = "https://api.twitter.com/1.1"
)

func (c *twitterServiceImpl) GetLoginUrl() (string, error) {
	rToken, loginUrl, err := c.consumer.GetRequestTokenAndUrl(c.callbackUrl)
	if err != nil {
		return "", err
	}
	c.requestTokenMap[rToken.Token] = rToken
	return loginUrl, err
}

func (c *twitterServiceImpl) AuthorizeToken(token string, verificationCode string) (*oauth.AccessToken, error) {
	rToken, ok := c.requestTokenMap[token]
	if !ok {
		return nil, errors.New("request token not found")
	}
	delete(c.requestTokenMap, token)
	aToken, err := c.consumer.AuthorizeToken(rToken, verificationCode)
	if err != nil {
		return nil, err
	}
	return aToken, err
}

func (c *twitterServiceImpl) GetUser(token *oauth.AccessToken) (*model.TwitterUser, error) {
	httpClient, err := c.consumer.MakeHttpClient(token)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Get(twitterAPIEndpoint + "/account/verify_credentials.json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("twitter api response error")
	}

	userObj := &userObject{}
	err = json.NewDecoder(res.Body).Decode(userObj)
	if err != nil {
		return nil, err
	}

	user := &model.TwitterUser{
		ID: userObj.ID,
		DisplayName: userObj.Name,
		ScreenName: userObj.ScreenName,
		ProfileImageUrl: userObj.ProfileImageURL,
		Biography: userObj.Description,
	}
	return user, nil
}

func (c *twitterServiceImpl) GetFollowers(token *oauth.AccessToken) ([]*model.TwitterUser, error) {
	httpClient, err := c.consumer.MakeHttpClient(token)
	if err != nil {
		return nil, err
	}

	var users []*model.TwitterUser
	nextCursor := -1

	for {
		res, err := httpClient.Get(twitterAPIEndpoint + "/followers/list.json?count=200&cursor=" + strconv.Itoa(nextCursor))
		if err != nil {
			return nil, err
		}
		log.Println(res.StatusCode)
		if res.StatusCode != 200 {
			return nil, errors.New("twitter api response error")
		}

		usersObj := &usersObject{}
		err = json.NewDecoder(res.Body).Decode(usersObj)
		res.Body.Close()

		if err != nil {
			return nil, err
		}

		for _, userObj := range usersObj.Users {
			users = append(users, &model.TwitterUser{
				ID: userObj.ID,
				DisplayName: userObj.Name,
				ScreenName: userObj.ScreenName,
				ProfileImageUrl: userObj.ProfileImageURL,
				Biography: userObj.Description,
			})
		}

		if usersObj.NextCursor == 0 {
			break
		} else {
			nextCursor = usersObj.NextCursor
		}
	}
	return users, nil
}

func (c *twitterServiceImpl) GetFollowees(token *oauth.AccessToken) ([]*model.TwitterUser, error) {
	httpClient, err := c.consumer.MakeHttpClient(token)
	if err != nil {
		return nil, err
	}

	var users []*model.TwitterUser
	nextCursor := -1

	for {
		res, err := httpClient.Get(twitterAPIEndpoint + "/friends/list.json?count=200&cursor=" + strconv.Itoa(nextCursor))
		if err != nil {
			return nil, err
		}
		log.Println(res.StatusCode)
		if res.StatusCode != 200 {
			return nil, errors.New("twitter api response error")
		}

		usersObj := &usersObject{}
		err = json.NewDecoder(res.Body).Decode(usersObj)
		res.Body.Close()

		if err != nil {
			return nil, err
		}

		for _, userObj := range usersObj.Users {
			users = append(users, &model.TwitterUser{
				ID: userObj.ID,
				DisplayName: userObj.Name,
				ScreenName: userObj.ScreenName,
				ProfileImageUrl: userObj.ProfileImageURL,
				Biography: userObj.Description,
			})
		}

		if usersObj.NextCursor == 0 {
			break
		} else {
			nextCursor = usersObj.NextCursor
		}
	}
	return users, nil
}

func (c *twitterServiceImpl) Search(token *oauth.AccessToken, query string) ([]*model.TwitterUser, error) {
	httpClient, err := c.consumer.MakeHttpClient(token)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Get(twitterAPIEndpoint + "/users/search.json?count=20&q=" + query)
	if err != nil {
		return nil, err
	}
	log.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return nil, errors.New("twitter api response error")
	}

	usersObj := &[]userObject{}
	err = json.NewDecoder(res.Body).Decode(usersObj)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	var users []*model.TwitterUser

	for _, userObj := range *usersObj {
		users = append(users, &model.TwitterUser{
			ID: userObj.ID,
			DisplayName: userObj.Name,
			ScreenName: userObj.ScreenName,
			ProfileImageUrl: userObj.ProfileImageURL,
			Biography: userObj.Description,
		})
	}
	return users, nil
}

func NewTwitterService(consumerKey, consumerSecret, callbackUrl string) TwitterService {
	consumer := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl: requestTokenURL,
			AuthorizeTokenUrl: authorizationURL,
			AccessTokenUrl: accessTokenURL,
		})
	return &twitterServiceImpl{
		consumer:        consumer,
		callbackUrl:     callbackUrl,
		requestTokenMap: map[string]*oauth.RequestToken{},
	}
}
