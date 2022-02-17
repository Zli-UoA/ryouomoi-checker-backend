package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserJWTService interface {
	CreateUserIDJWT(userID int64) (string, error)
	GetUserIDFromJWT(token string) (int64, error)
}

type userJWTServiceImpl struct {
	secret string
}

func (u *userJWTServiceImpl) CreateUserIDJWT(userID int64) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodHS256.Alg()))
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	tokenString, err := token.SignedString([]byte(u.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *userJWTServiceImpl) GetUserIDFromJWT(token string) (int64, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(u.secret), nil
	})
	if err != nil {
		return 0, err
	}
	return t.Claims.(jwt.MapClaims)["userID"].(int64), nil
}

func NewUserJWTService(secret string) UserJWTService {
	return &userJWTServiceImpl{secret: secret}
}
