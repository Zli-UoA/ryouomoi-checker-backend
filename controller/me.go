package controller

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MeController struct {
	ujs    service.UserJWTService
	glpu   usecase.GetLovePointUsecase
	glptsu usecase.GetLovePointsUseCase
	dclu   usecase.DeleteCurrentLoverUseCase
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
	}
	err = m.dclu.Execute(userID, req.ReasonID, req.AllowShare)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.Status(200)
}

func NewMeController(ujs service.UserJWTService, glpu usecase.GetLovePointUsecase, glptsu usecase.GetLovePointsUseCase, dclu usecase.DeleteCurrentLoverUseCase) *MeController {
	return &MeController{
		ujs:    ujs,
		glpu:   glpu,
		glptsu: glptsu,
		dclu:   dclu,
	}
}
