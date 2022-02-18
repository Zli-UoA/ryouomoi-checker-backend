package controller

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
)

type FriendsController struct {
	ujs service.UserJWTService
	fsu usecase.FriendsSearchUseCase
}

func convertToJson(twitterUser *model.TwitterUser) *TwitterUser {
	return &TwitterUser{
		ID:          twitterUser.ID,
		ScreenName:  twitterUser.ScreenName,
		DisplayName: twitterUser.DisplayName,
		ImageUrl:    twitterUser.ProfileImageUrl,
		Biography:   twitterUser.Biography,
	}
}

func (f *FriendsController) FriendsSearch(c *gin.Context) {
	token, err := GetAuthToken(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	userID, err := f.ujs.GetUserIDFromJWT(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	query := c.Query("query")
	searchResult, err := f.fsu.Execute(userID, query)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	jsonArray := make([]*TwitterUser, len(searchResult))
	for i, user := range searchResult {
		jsonArray[i] = convertToJson(user)
	}
	c.JSON(200, jsonArray)
}

func NewFriendsController(ujs service.UserJWTService, fsu usecase.FriendsSearchUseCase) *FriendsController {
	return &FriendsController{
		ujs: ujs,
		fsu: fsu,
	}
}
