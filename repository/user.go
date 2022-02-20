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
	DeleteLovePoint(userID, loverUserID int64) error
	GetLovePoints(userID int64) ([]*model.UserLovePoint, error)
	SetLovePoint(point *model.UserLovePoint) (*model.UserLovePoint, error)
	GetCurrentCouple(userID int64) (*model.Couple, error)
	CreateCouple(couple *model.Couple) (*model.Couple, error)
	UpdateCouple(couple *model.Couple) (*model.Couple, error)
	CreateBrokeReport(report *model.BrokeReport) (*model.BrokeReport, error)
	GetBrokeReport(userID, coupleID int64) (*model.BrokeReport, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}
//↓User構築用
type Users_id struct {
	user_id_1 int64 `db:"user_id_1"`
	user_id_2 int64 `db:"user_id_2"`
}
//↓dbから値を抜く時だけ使う
type brokeReport struct {
	id              int64
	couple_id       int64
	user_id         int64
	broke_reason_id int
	allow_share     bool
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
		User1:     convertToUser(couple.User1),
		User2:     convertToUser(couple.User2),
		CreatedAt: couple.CreatedAt,
		BrokenAt:  couple.BrokenAt,
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
func (u *userRepositoryImpl) GenerateTwitterUserFromUserID(userID int64) (*TwitterUser, error) {
	User := TwitterUser{}
	err := u.db.Get(&User, "SELECT * FROM twitter_users WHERE twitter_id=?", userID)
	if err != nil {
		return nil, err
	}
	return &User, nil
}
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}
func (u *userRepositoryImpl) DeleteLovePoint(userID, loverUserID int64) error {
	_, err := u.db.Exec("DELETE FROM user_love_points WHERE user_id = ? AND lover_user_id = ?", userID, loverUserID)
	return err
}

func (u *userRepositoryImpl) GetLovePoints(userID int64) ([]*model.UserLovePoint, error) {
	panic("implement me")
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
	timeNow := time.Now()
	_, err := u.db.Exec("INSERT INTO couples (user_id_1,user_id_2,created_at) VALUES (?,?,?)", userID1, userID2, timeNow)
	if err != nil {
		return nil, err
	}
	couple.CreatedAt = timeNow
	return couple, nil
}
func (u *userRepositoryImpl) GetLovePoint(userID, loverUserID int64) (*model.UserLovePoint, error) { // test done
	userLovePoint := UserLovePoint{}
	err := u.db.Get(&userLovePoint, "SELECT u.id id, t.twitter_id \"user.twitter_id\", t.screen_name \"user.screen_name\", t.display_name \"user.display_name\", t.profile_image_url \"user.profile_image_url\", t.biography \"user.biography\", u.lover_user_id lover_user_id, u.love_point love_point FROM user_love_points u JOIN twitter_users t ON t.twitter_id = u.user_id WHERE user_id = ? AND lover_user_id = ?", userID, loverUserID)
	if err != nil {
		return nil, err
	}
	return convertToUserLovePoint(&userLovePoint), nil
}

func (u *userRepositoryImpl) GetLatestBrokenCouple(userID int64) (*model.Couple, error) { //test done 後で綺麗にする
	cp := Couple{}
	users_id := Users_id{}
	User_1 := TwitterUser{}
	User_2 := TwitterUser{}
	err := u.db.QueryRow("select id,user_id_1 ,user_id_2,created_at,broken_at from twitter_users join couples on twitter_users.twitter_id = couples.user_id_1 WHERE couples.user_id_1 = ? OR couples.user_id_2 = ? ORDER BY broken_at DESC LIMIT 1", 1, 1).Scan(&cp.ID, &users_id.user_id_1, &users_id.user_id_2, &cp.CreatedAt, &cp.BrokenAt)
	if err != nil {
		panic(err)
	}
	err = u.db.Get(&User_1, "SELECT * FROM twitter_users WHERE twitter_id=?", users_id.user_id_1)
	if err != nil {
		return nil, err
	}
	err = u.db.Get(&User_2, "SELECT * FROM twitter_users WHERE twitter_id=?", users_id.user_id_2)
	if err != nil {
		return nil, err
	}
	cp.User1 = &User_1
	cp.User2 = &User_2
	return convertToCouple(&cp), nil
}

func (u *userRepositoryImpl) GetCurrentCouple(userID int64) (*model.Couple, error) { //test done 汚いので後で修正
	//一件もなかったらnil
	cp := Couple{}
	user_id := Users_id{}
	User_1 := TwitterUser{}
	User_2 := TwitterUser{}
	err := u.db.QueryRow("SELECT user_id_1,user_id_2 FROM couples WHERE broken_at IS NULL AND (user_id_1 = ? OR user_id_2 = ?);", userID, userID).Scan(&user_id.user_id_1, &user_id.user_id_2)
	if err != nil {
		return nil, err
	}
	err = u.db.Get(&User_1, "SELECT * FROM twitter_users WHERE twitter_id=?", user_id.user_id_1)
	if err != nil {
		return nil, err
	}
	err = u.db.Get(&User_2, "SELECT * FROM twitter_users WHERE twitter_id=?", user_id.user_id_2)
	if err != nil {
		return nil, err
	}
	cp.User1 = &User_1
	cp.User2 = &User_2
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

func (u *userRepositoryImpl) CreateBrokeReport(report *model.BrokeReport) (*model.BrokeReport, error) { //test done
	_, err := u.db.Exec("INSERT INTO couple_broke_reports (couple_id,user_id, broke_reason_id, allow_share) VALUES (?,?, ?,?)", report.Couple.ID, report.User.ID, report.BrokeReasonID, report.AllowShare)
	if err != nil {
		return nil, err
	}
	return report, nil

}
func (u *userRepositoryImpl) GenerateCoupleFromCoupleID(coupleID int64) (*Couple, error) {
	cp := Couple{}
	users_id := Users_id{}
	User_1 := TwitterUser{}
	User_2 := TwitterUser{}
	err := u.db.QueryRow("SELECT * FROM couples WHERE id=?", coupleID).Scan(&cp.ID, &users_id.user_id_1, &users_id.user_id_2, &cp.CreatedAt, &cp.BrokenAt)
	if err != nil {
		return nil, err
	}
	err = u.db.Get(&User_1, "SELECT * FROM twitter_users WHERE twitter_id=?", users_id.user_id_1)
	if err != nil {
		return nil, err
	}
	err = u.db.Get(&User_2, "SELECT * FROM twitter_users WHERE twitter_id=?", users_id.user_id_2)
	if err != nil {
		return nil, err
	}
	cp.User1 = &User_1
	cp.User2 = &User_2
	return &cp, nil
}
func (u *userRepositoryImpl) GetBrokeReport(userID, coupleID int64) (*model.BrokeReport, error) { //test done?
	repo := brokeReport{}
	err := u.db.QueryRow("select * from couple_broke_reports where user_id = ? and couple_id = ?", userID, coupleID).Scan(&repo.id, &repo.couple_id, &repo.user_id, &repo.broke_reason_id, &repo.allow_share)
	if err != nil {
		return nil, err
	}
	res_repo := model.BrokeReport{}
	res_repo.ID = repo.id
	res_repo.BrokeReasonID = repo.broke_reason_id
	res_repo.AllowShare = repo.allow_share
	cp, err := u.GenerateCoupleFromCoupleID(repo.couple_id)
	if err != nil {
		return nil, err
	}
	res_repo.Couple = convertToCouple(cp)
	user, err := u.GenerateTwitterUserFromUserID(userID)
	res_repo.User = convertToUser(user)
	return &res_repo, nil
}
