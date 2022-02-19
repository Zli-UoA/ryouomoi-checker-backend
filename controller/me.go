package controller

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
)

type MeController struct {
	ujs  service.UserJWTService
	dclu usecase.DeleteCurrentLoverUseCase
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

func NewMeController(ujs service.UserJWTService, dclu usecase.DeleteCurrentLoverUseCase) *MeController {
	return &MeController{
		ujs:  ujs,
		dclu: dclu,
	}
}
