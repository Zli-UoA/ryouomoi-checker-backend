package controller

import (
	"errors"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FriendsController struct {
	ujs  service.UserJWTService
	fsu  usecase.FriendsSearchUseCase
	slpu usecase.SetLovePointUseCase
	gfu  usecase.GetFolloweesUseCase
	fru  usecase.GetFollowersUseCase
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

func (f *FriendsController) GetFollowees(c *gin.Context) {
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
	searchResult, err := f.gfu.Execute(userID)
	jsonArray := make([]*TwitterUser, len(searchResult))
	for i, user := range searchResult {
		jsonArray[i] = convertToJson(user)
	}
	c.JSON(200, jsonArray)
}

func (f *FriendsController) GetFollowers(c *gin.Context) {
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
	searchResult, err := f.fru.Execute(userID)
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
		return
	}
	lover, err := f.slpu.Execute(userID, loverUserID, req.LovePoint)
	if err != nil {
		var target *usecase.BrokenCoupleNotExpiredError
		if errors.As(err, &target) {
			c.JSON(425, gin.H{
				"message":    "破局してから1ヶ月以上経過していません。",
				"remainDays": target.RemainDays,
			})
			return
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	res := MatchResult{
		MatchSuccess: lover != nil,
	}
	if lover != nil {
		res.Lover = &Lover{
			User:        convertToJson(lover.User),
			TalkRoomUrl: lover.TalkRoomUrl,
		}
	}
	c.JSON(200, res)
}

func NewFriendsController(
	ujs service.UserJWTService,
	fsu usecase.FriendsSearchUseCase,
	slpu usecase.SetLovePointUseCase,
	gfu usecase.GetFolloweesUseCase,
	fru usecase.GetFollowersUseCase,
) *FriendsController {
	return &FriendsController{
		ujs:  ujs,
		fsu:  fsu,
		slpu: slpu,
		gfu:  gfu,
		fru:  fru,
	}
}
