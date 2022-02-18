package repository

import (
	"time"

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
	GetCurrentCouple(userID int64) (*model.Couple, error)
	CreateCouple(couple *model.Couple) (*model.Couple, error)
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
func (u *userRepositoryImpl) GetLovePoint(userID, loverUserID int64) (*model.UserLovePoint, error) { // test done
	userLovePoint := UserLovePoint{}
	err := u.db.Get(&userLovePoint, "SELECT u.id id, t.twitter_id \"user.twitter_id\", t.screen_name \"user.screen_name\", t.display_name \"user.display_name\", t.profile_image_url \"user.profile_image_url\", t.biography \"user.biography\", u.lover_user_id lover_user_id, u.love_point love_point FROM user_love_points u JOIN twitter_users t ON t.twitter_id = u.user_id WHERE user_id = ? AND lover_user_id = ?", userID, loverUserID)
	if err != nil {
		return nil, err
	}
	return convertToUserLovePoint(&userLovePoint), nil
}
func (u *userRepositoryImpl) SetLovePoint(point *model.UserLovePoint) (*model.UserLovePoint, error) { //test done
	res, err := u.db.Exec("UPDATE user_love_points SET love_point= ? where user_id= ? AND lover_user_id= ?", point.LovePoint, point.ID, point.LoverUserID)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		res, err = u.db.Exec("INSERT INTO user_love_points (user_id, lover_user_id, love_point) VALUES (?, ?, ?)", point.UserID, point.LoverUserID, point.LovePoint)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		point.ID = id
	}
	return point, nil
}

func (u *userRepositoryImpl) UpdateCouple(couple *model.Couple) (*model.Couple, error) { //test done
	userID1 := couple.User1.ID
	userID2 := couple.User2.ID
	time_now := time.Now()
	_, err := u.db.Exec("UPDATE couples SET broken_at =  ? where user_id_1 = ? AND user_id_2= ?", time_now, userID1, userID2)
	if err != nil {
		return nil, err
	}
	couple.BrokenAt = &time_now
	return couple, nil
}

func (u *userRepositoryImpl) CreateCouple(couple *model.Couple) (*model.Couple, error) { //test done
	userID1 := couple.User1.ID
	userID2 := couple.User2.ID
	time_now := time.Now()
	_, err := u.db.Exec("INSERT INTO couples (user_id_1,user_id_2,created_at) VALUES (?,?,?)", userID1, userID2, time_now)
	if err != nil {
		return nil, err
	}
	couple.CreatedAt = &time_now
	return couple, nil
}
func (u *userRepositoryImpl) GetLatestBrokenCouple(userID int64) (*model.Couple, error) { //未テスト TwitterUserのとこどうしよう
	cp := Couple{} //user1でいい?user2の可能性もありそう
	err := u.db.Get(&cp, "SELECT * FROM couples WHERE user_id_1 = ? OR user_id_2 = ? ORDER BY broken_at DESC LIMIT 1", userID, userID)
	if err != nil {
		return nil, err
	}
	return convertToCouple(&cp), nil
}
func (u *userRepositoryImpl) GetCurrentCouple(userID int64) (*model.Couple, error) { //未テスト 今のlover
	//一件もなかったらnil
	cp := Couple{}
	err := u.db.Get(&cp, "SELECT * FROM couples WHERE broken_at IS NULL AND (user_id_1 = ? OR user_id_2 = ?);", userID, userID) //sqlはオッケー
	if err != nil {
		return nil, err
	}
	return convertToCouple(&cp), err
}
func (u *userRepositoryImpl) GetUser(id int64) (*model.User, error) { //未テスト そのまま(?のとこだけ変えた)
	twitterUser := TwitterUser{}
	err := u.db.Get(&twitterUser, "SELECT * FROM twitter_users WHERE twitter_id=?", id)
	if err != nil {
		return nil, err
	}
	return convertToUser(&twitterUser), nil
}

func (u *userRepositoryImpl) CreateUser(user *model.User) (*model.User, error) { //未テスト　そのまま
	twitterUser := convertToTwitterUser(user)
	_, err := u.db.NamedExec("INSERT INTO twitter_users (twitter_id, screen_name, display_name, profile_image_url, biography, access_token, access_token_secret) VALUES (:twitter_id, :screen_name, :display_name, :profile_image_url, :biography, :access_token, :access_token_secret)", twitterUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepositoryImpl) UpdateUser(user *model.User) (*model.User, error) { //未テスト そのまま
	twitterUser := convertToTwitterUser(user)
	_, err := u.db.NamedExec("UPDATE twitter_users SET screen_name=:screen_name, display_name=:display_name, profile_image_url=:profile_image_url, biography=:biography, access_token=:access_token, access_token_secret=:access_token_secret WHERE twitter_id=:twitter_id", twitterUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}