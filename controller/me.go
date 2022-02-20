package controller

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
)

type MeController struct {
	ujs    service.UserJWTService
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

func NewMeController(ujs service.UserJWTService, glptsu usecase.GetLovePointsUseCase, dclu usecase.DeleteCurrentLoverUseCase) *MeController {
	return &MeController{
		ujs:    ujs,
		glptsu: glptsu,
		dclu:   dclu,
	}
}
