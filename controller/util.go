package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetAuthToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer") {
		return "", errors.New("AuthorizationヘッダーはBearer形式である必要があります")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}
