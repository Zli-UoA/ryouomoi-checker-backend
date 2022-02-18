package controller

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FriendsController struct {
	ujs service.UserJWTService
	fsu usecase.FriendsSearchUseCase
	slpu usecase.SetLovePointUseCase
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

func (f *FriendsController) SetLovePoint(c *gin.Context) {
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
	loverUserIDStr := c.Param("id")
	loverUserID, err := strconv.ParseInt(loverUserIDStr, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	req := &LovePoint{}
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	}
	matchSuccess, err := f.slpu.Execute(userID, loverUserID, req.LovePoint)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	res := MatchResult{MatchSuccess: matchSuccess}
	c.JSON(200, res)
}

func NewFriendsController(ujs service.UserJWTService, fsu usecase.FriendsSearchUseCase, slpu usecase.SetLovePointUseCase) *FriendsController {
	return &FriendsController{
		ujs: ujs,
		fsu: fsu,
		slpu: slpu,
	}
}
