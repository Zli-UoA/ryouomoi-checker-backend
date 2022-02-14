package repository

import "github.com/Zli-UoA/ryouomoi-checker-backend/model"

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
