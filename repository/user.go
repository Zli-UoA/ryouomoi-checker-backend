package repository

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetUser(id int64) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	GetLovePoint(userID, loverUserID int64) (*model.UserLovePoint, error)
	SetLovePoint(point *model.UserLovePoint) (*model.UserLovePoint, error)
	GetCouple(userID int64) (*model.Couple, error)
	CreateCouple(userID1, userID2 int64) (*model.Couple, error)
	UpdateCouple(couple *model.Couple) (*model.Couple, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

func convertToUser(twitterUser *TwitterUser) *model.User {
	user := model.User{
		ID:                 twitterUser.TwitterID,
		ScreenName:         twitterUser.ScreenName,
		DisplayName:        twitterUser.DisplayName,
		ProfileImageUrl:    twitterUser.ProfileImageUrl,
		Biography:          twitterUser.Biography,
		TwitterAccessToken: service.CreateAccessToken(twitterUser.AccessToken, twitterUser.AccessTokenSecret),
	}
	return &user
}

func convertToTwitterUser(user *model.User) *TwitterUser {
	twitterUser := TwitterUser{
		TwitterID:         user.ID,
		ScreenName:        user.ScreenName,
		DisplayName:       user.DisplayName,
		ProfileImageUrl:   user.ProfileImageUrl,
		Biography:         user.Biography,
		AccessToken:       user.TwitterAccessToken.Token,
		AccessTokenSecret: user.TwitterAccessToken.Secret,
	}
	return &twitterUser
}

func (u *userRepositoryImpl) GetUser(id int64) (*model.User, error) {
	twitterUser := TwitterUser{}
	err := u.db.Get(&twitterUser, "SELECT * FROM twitter_users WHERE twitter_id=$1", id)
	if err != nil {
		return nil, err
	}
	return convertToUser(&twitterUser), nil
}

func (u *userRepositoryImpl) CreateUser(user *model.User) (*model.User, error) {
	twitterUser := convertToTwitterUser(user)
	_, err := u.db.NamedExec("INSERT INTO twitter_users (twitter_id, screen_name, display_name, profile_image_url, biography, access_token, access_token_secret) VALUES (:twitter_id, :screen_name, :display_name, :profile_image_url, :biography, :access_token, :access_token_secret)", twitterUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepositoryImpl) UpdateUser(user *model.User) (*model.User, error) {
	twitterUser := convertToTwitterUser(user)
	_, err := u.db.NamedExec("UPDATE twitter_users SET screen_name=:screen_name, display_name=:display_name, profile_image_url=:profile_image_url, biography=:biography, access_token=:access_token, access_token_secret=:access_token_secret WHERE twitter_id=:twitter_id", twitterUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepositoryImpl) GetLovePoint(userID, loverUserID int64) (*model.UserLovePoint, error) {
	//userLoverPoint := UserLoverPoint{}
	//err := u.db.Get(&userLoverPoint, "SELECT ulp.id id, lover_user_id, love_point, tu.twitter_id , screen_name, display_name, profile_image_url, biography, access_token, access_token_secret")
	panic("implement me")
}

func (u *userRepositoryImpl) SetLovePoint(point *model.UserLovePoint) (*model.UserLovePoint, error) {
	panic("implement me")
}

func (u *userRepositoryImpl) GetCouple(userID int64) (*model.Couple, error) {
	panic("implement me")
}

func (u *userRepositoryImpl) CreateCouple(userID1, userID2 int64) (*model.Couple, error) {
	panic("implement me")
}

func (u *userRepositoryImpl) UpdateCouple(couple *model.Couple) (*model.Couple, error) {
	panic("implement me")
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}
