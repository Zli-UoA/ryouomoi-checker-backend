package controller

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
)

type TwitterController struct {
	usecase usecase.GetTwitterLoginUrlUseCase
}

func (tc *TwitterController) GetTwitterLoginUrl(c *gin.Context) {
	url, err := tc.usecase.Execute()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	json := TwitterLoginUrlJson{
		LoginUrl: url,
	}
	c.JSON(200, json)
}

func NewTwitterController(usecase usecase.GetTwitterLoginUrlUseCase) *TwitterController {
	return &TwitterController{
		usecase: usecase,
	}
}
