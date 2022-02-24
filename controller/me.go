package controller

import (
	"database/sql"
	"errors"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MeController struct {
	ujs    service.UserJWTService
	glpu   usecase.GetLovePointUsecase
	glptsu usecase.GetLovePointsUseCase
	gclu   usecase.GetCurrentLoverUsecase
	dclu   usecase.DeleteCurrentLoverUseCase
	gliuu  usecase.GetLoggedInUserUseCase
}

func (m *MeController) GetLoggedInUser(c *gin.Context) {
	token, err := GetAuthToken(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	userID, err := m.ujs.GetUserIDFromJWT(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, err := m.gliuu.Execute(userID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	userJson := &TwitterUser{
		ID:          user.ID,
		ScreenName:  user.ScreenName,
		DisplayName: user.DisplayName,
		ImageUrl:    user.ProfileImageUrl,
		Biography:   user.Biography,
	}
	c.JSON(200, userJson)
}

func (m *MeController) GetLovePoints(c *gin.Context) {
	token, err := GetAuthToken(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	userID, err := m.ujs.GetUserIDFromJWT(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	lovePoints, err := m.glptsu.Execute(userID)
	res := make([]*UserLovePoint, len(lovePoints))
	for i, point := range lovePoints {
		res[i] = &UserLovePoint{
			LoverUser: convertToJson(point.LoverUser),
			LovePoint: point.LovePoint,
		}
	}
	c.JSON(200, res)
}

func (m *MeController) GetLovePoint(c *gin.Context) {
	token, err := GetAuthToken(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	userID, err := m.ujs.GetUserIDFromJWT(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	loverIDStr := c.Param("id")
	loverID, err := strconv.ParseInt(loverIDStr, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	lovePoint, err := m.glpu.Execute(userID, loverID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	jsonLP := LovePoint{LovePoint: lovePoint.LovePoint}
	c.JSON(200, jsonLP)
}

func (m *MeController) GetCurrentLover(c *gin.Context) {
	token, err := GetAuthToken(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	userID, err := m.ujs.GetUserIDFromJWT(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	currentLover, err := m.gclu.Execute(userID)
	if err != nil {
		if errors.Is(err, usecase.BrokenCoupleError) {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, usecase.BrokenReportError) {
			c.JSON(410, gin.H{
				"message": "相手に振られてしまいました。",
			})
			return
		}
		if errors.Is(err, usecase.BrokenCoupleNotExpiredError) {
			c.JSON(425, gin.H{
				"message": "破局してから1ヶ月以上経過していません。",
			})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{
				"message": "現在、彼氏・彼女はいません。",
			})
			return
		}
		return
	}
	jsonLover := &TwitterUser{
		ID:          currentLover.ID,
		ScreenName:  currentLover.ScreenName,
		DisplayName: currentLover.DisplayName,
		ImageUrl:    currentLover.ProfileImageUrl,
		Biography:   currentLover.Biography,
	}
	c.JSON(200, jsonLover)
}

func (m *MeController) DeleteCurrentLover(c *gin.Context) {
	token, err := GetAuthToken(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	userID, err := m.ujs.GetUserIDFromJWT(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	req := &BrokeReport{}
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = m.dclu.Execute(userID, req.ReasonID, req.AllowShare)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, usecase.ErrorBrokeReportAlreadyExists) {
			c.JSON(404, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	c.Status(200)
}

func NewMeController(
	ujs service.UserJWTService,
	glpu usecase.GetLovePointUsecase,
	glptsu usecase.GetLovePointsUseCase,
	gclu usecase.GetCurrentLoverUsecase,
	dclu usecase.DeleteCurrentLoverUseCase,
	gliuu usecase.GetLoggedInUserUseCase,
) *MeController {
	return &MeController{
		ujs:    ujs,
		glpu:   glpu,
		glptsu: glptsu,
		gclu:   gclu,
		dclu:   dclu,
		gliuu:  gliuu,
	}
}
