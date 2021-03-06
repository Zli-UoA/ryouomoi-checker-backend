package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/mrjones/oauth"
	"net/url"
	"strconv"
	"strings"
)

type TwitterService interface {
	GetLoginUrl() (string, error)
	AuthorizeToken(token, verificationCode string) (*oauth.AccessToken, error)
	GetUser(token *oauth.AccessToken) (*model.TwitterUser, error)
	GetFollowees(token *oauth.AccessToken) ([]*model.TwitterUser, error)
	GetFollowers(token *oauth.AccessToken) ([]*model.TwitterUser, error)
	Search(token *oauth.AccessToken, query string) ([]*model.TwitterUser, error)
	Lookup(token *oauth.AccessToken, ids []int64) ([]*model.TwitterUser, error)
	Show(token *oauth.AccessToken, id int64) (*model.TwitterUser, error)
	SendTweet(token *oauth.AccessToken, content string) error
	SendDirectMessage(token *oauth.AccessToken, toUserID int64, content string) error
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
	Users      []userObject `json:"users"`
	NextCursor int          `json:"next_cursor"`
}

type twitterServiceImpl struct {
	consumer        *oauth.Consumer
	callbackUrl     string
	requestTokenMap map[string]*oauth.RequestToken
}

func convertToModel(userObj *userObject) *model.TwitterUser {
	normalUrl := userObj.ProfileImageURL
	originUrl := strings.Join(strings.Split(normalUrl, "_normal"), "")
	return &model.TwitterUser{
		ID:              userObj.ID,
		DisplayName:     userObj.Name,
		ScreenName:      userObj.ScreenName,
		ProfileImageUrl: originUrl,
		Biography:       userObj.Description,
	}
}

const (
	requestTokenURL    = "https://api.twitter.com/oauth/request_token"
	authorizationURL   = "https://api.twitter.com/oauth/authenticate"
	accessTokenURL     = "https://api.twitter.com/oauth/access_token"
	twitterAPIEndpoint = "https://api.twitter.com/1.1"
)

func CreateAccessToken(accessToken, accessTokenSecret string) *oauth.AccessToken {
	return &oauth.AccessToken{
		Token:  accessToken,
		Secret: accessTokenSecret,
	}
}

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

	user := convertToModel(userObj)
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
			users = append(users, convertToModel(&userObj))
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
			users = append(users, convertToModel(&userObj))
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

	escapedQuery := url.QueryEscape(query)
	res, err := httpClient.Get(twitterAPIEndpoint + "/users/search.json?count=20&q=" + escapedQuery)
	if err != nil {
		return nil, err
	}
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
		users = append(users, convertToModel(&userObj))
	}
	return users, nil
}

func (c *twitterServiceImpl) Lookup(token *oauth.AccessToken, ids []int64) ([]*model.TwitterUser, error) {
	httpClient, err := c.consumer.MakeHttpClient(token)
	if err != nil {
		return nil, err
	}
	var urlBuilder strings.Builder
	urlBuilder.WriteString(twitterAPIEndpoint + "/users/lookup.json?user_id=")
	for i, id := range ids {
		urlBuilder.WriteString(strconv.FormatInt(id, 10))
		if i < len(ids)-1 {
			urlBuilder.WriteString(",")
		}
	}
	res, err := httpClient.Get(urlBuilder.String())
	if err != nil {
		return nil, err
	}
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
		users = append(users, convertToModel(&userObj))
	}
	return users, nil
}

func (c *twitterServiceImpl) Show(token *oauth.AccessToken, id int64) (*model.TwitterUser, error) {
	httpClient, err := c.consumer.MakeHttpClient(token)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Get(twitterAPIEndpoint + "/users/show.json?user_id=" + strconv.FormatInt(id, 10))
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

	user := convertToModel(userObj)
	return user, nil
}

func (c *twitterServiceImpl) SendTweet(token *oauth.AccessToken, content string) error {
	httpClient, err := c.consumer.MakeHttpClient(token)
	if err != nil {
		return err
	}
	status := url.QueryEscape(content)
	_, err = httpClient.Post(
		twitterAPIEndpoint+"/statuses/update.json?status="+status,
		"application/json",
		nil,
	)
	return nil
}

func (c *twitterServiceImpl) SendDirectMessage(token *oauth.AccessToken, toUserID int64, content string) error {
	httpClient, err := c.consumer.MakeHttpClient(token)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"event": map[string]interface{}{
			"type": "message_create",
			"message_create": map[string]interface{}{
				"target": map[string]interface{}{
					"recipient_id": toUserID,
				},
				"message_data": map[string]interface{}{
					"text": content,
				},
			},
		},
	}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = httpClient.Post(
		twitterAPIEndpoint+"/direct_messages/events/new.json",
		"application/json",
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return err
	}
	return nil
}

func NewTwitterService(consumerKey, consumerSecret, callbackUrl string) TwitterService {
	consumer := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   requestTokenURL,
			AuthorizeTokenUrl: authorizationURL,
			AccessTokenUrl:    accessTokenURL,
		})
	return &twitterServiceImpl{
		consumer:        consumer,
		callbackUrl:     callbackUrl,
		requestTokenMap: map[string]*oauth.RequestToken{},
	}
}
