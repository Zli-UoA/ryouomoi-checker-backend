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

//dbとアプリ内のuserの変換
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
func convertToUserLovePoint(userLovePoint *UserLovePoint) *model.UserLovePoint {
	point := model.UserLovePoint{
		ID:          userLovePoint.ID,
		UserID:      userLovePoint.User.TwitterID,
		LoverUserID: userLovePoint.LoverUserID,
		LovePoint:   userLovePoint.LovePoint,
	}
	return &point

}
func convertToCouple(couple *Couple) *model.Couple {
	cp := model.Couple{
		ID:        couple.ID,
		User1:     convertToUser(&couple.UserID1),
		User2:     convertToUser(&couple.UserID2),
		CreatedAt: &couple.CreatedAt,
		BrokenAt:  &couple.BrokenAt,
	}
	return &cp
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
	err := u.db.Get(&twitterUser, "SELECT * FROM twitter_users WHERE twitter_id=?", id)
	if err != nil {
		return nil, err
	}
	return convertToUser(&twitterUser), nil
}

func (u *userRepositoryImpl) CreateUser(user *model.User) (*model.User, error) { //未テスト
	twitterUser := convertToTwitterUser(user)
	_, err := u.db.NamedExec("INSERT INTO twitter_users (twitter_id, screen_name, display_name, profile_image_url, biography, access_token, access_token_secret) VALUES (:twitter_id, :screen_name, :display_name, :profile_image_url, :biography, :access_token, :access_token_secret)", twitterUser)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (u *userRepositoryImpl) UpdateUser(user *model.User) (*model.User, error) { //未テスト
	twitterUser := convertToTwitterUser(user)
	_, err := u.db.NamedExec("UPDATE twitter_users SET screen_name=:screen_name, display_name=:display_name, profile_image_url=:profile_image_url, biography=:biography, access_token=:access_token, access_token_secret=:access_token_secret WHERE twitter_id=:twitter_id", twitterUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepositoryImpl) GetLovePoint(userID, loverUserID int64) (*model.UserLovePoint, error) { // test done
	userLovePoint := UserLovePoint{}
	err := u.db.Get(&userLovePoint, "SELECT * FROM user_lover_points WHERE user_id = ? AND lover_user_id = ?", userID, loverUserID)
	if err != nil {
		return nil, err
	}
	return convertToUserLovePoint(&userLovePoint), nil
}
func (u *userRepositoryImpl) SetLovePoint(point *UserLovePoint) (*model.UserLovePoint, error) {//done
	_, err := u.db.Exec("UPDATE user_love_points SET love_point= ? where user_id= ? AND lover_user_id= ?", point.LovePoint,point.ID,point.LoverUserID)
	if err != nil {
		return nil, err
	}
	points := convertToUserLovePoint(point)
	return points, nil
}

func (u *userRepositoryImpl) GetLatestBrokenCouple(userID int64) (*model.Couple, error) { //一番最近の破
	//一件もなかったらnilを返す
	panic("implement me")
}
func (u *userRepositoryImpl) GetCurrentCouple(userID int64) (*model.Couple, error) { //今のlover
	//一件もなかったらnil
	panic("implement me")
}
func (u *userRepositoryImpl) CreateCouple(userID1, userID2 int64) (*model.Couple, error) { //test done
	couple := Couple{}
	_, err := u.db.Exec("INSERT INTO couples (user_id_1,user_id_2,created_at) VALUES ($1,$2,CURRENT_TIMESTAMP)", userID1, userID2)
	if err != nil {
		return nil, err
	}
	return convertToCouple(&couple), nil
}

func (u *userRepositoryImpl) UpdateCouple(couple *model.Couple) (*model.Couple, error) { //破局
	panic("implement me")
}
