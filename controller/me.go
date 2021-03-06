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
	dlpu   usecase.DeleteLovePointUseCase
	glptsu usecase.GetLovePointsUseCase
	gclu   usecase.GetCurrentLoverUsecase
	gcedu  usecase.GetCoupleElapsedDaysUseCase
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

func (m *MeController) DeleteLovePoint(c *gin.Context) {
	token, err := GetAuthToken(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	userID, err := m.ujs.GetUserIDFromJWT(token)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	loverIDStr := c.Param("id")
	loverID, err := strconv.ParseInt(loverIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = m.dlpu.Execute(userID, loverID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{
				"message": "??????????????????????????????????????????",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.Status(200)
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
			c.JSON(404, gin.H{
				"message": "??????????????????????????????????????????",
			})
			return
		}
		var brokeReportNotFoundError *usecase.BrokeReportNotFoundError
		if errors.As(err, &brokeReportNotFoundError) {
			c.JSON(410, gin.H{
				"message": "??????????????????????????????????????????",
				"lover": &TwitterUser{
					ID:          brokeReportNotFoundError.Lover.ID,
					ScreenName:  brokeReportNotFoundError.Lover.ScreenName,
					DisplayName: brokeReportNotFoundError.Lover.DisplayName,
					ImageUrl:    brokeReportNotFoundError.Lover.ProfileImageUrl,
					Biography:   brokeReportNotFoundError.Lover.Biography,
				},
			})
			return
		}
		var brokenCoupleNotExpiredError *usecase.BrokenCoupleNotExpiredError
		if errors.As(err, &brokenCoupleNotExpiredError) {
			c.JSON(425, gin.H{
				"message":    "??????????????????1???????????????????????????????????????",
				"remainDays": brokenCoupleNotExpiredError.RemainDays,
			})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{
				"message": "??????????????????????????????????????????",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	jsonLover := &Lover{
		User:        convertToJson(currentLover.User),
		TalkRoomUrl: currentLover.TalkRoomUrl,
	}
	c.JSON(200, jsonLover)
}

func (m *MeController) GetCoupleElapsedDays(c *gin.Context) {
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
	days, err := m.gcedu.Execute(userID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"days": days,
	})
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
		if errors.Is(err, usecase.ErrorBrokeReportAlreadyExists) {
			c.JSON(410, gin.H{
				"message": "???????????????????????????????????????",
			})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{
				"message": err.Error(),
			})
			return
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
	dlpu usecase.DeleteLovePointUseCase,
	glptsu usecase.GetLovePointsUseCase,
	gclu usecase.GetCurrentLoverUsecase,
	gcedu usecase.GetCoupleElapsedDaysUseCase,
	dclu usecase.DeleteCurrentLoverUseCase,
	gliuu usecase.GetLoggedInUserUseCase,
) *MeController {
	return &MeController{
		ujs:    ujs,
		glpu:   glpu,
		dlpu:   dlpu,
		glptsu: glptsu,
		gclu:   gclu,
		gcedu:  gcedu,
		dclu:   dclu,
		gliuu:  gliuu,
	}
}
