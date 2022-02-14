package service

type UserJWTService interface {
	CreateUserIDJWT(userID int64) (string, error)
	GetUserIDFromJWT(token string) (int64, error)
}
