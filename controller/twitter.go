package controller

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
)

type TwitterController struct {
	frontRedirectUrl string
	gtluu            usecase.GetTwitterLoginUrlUseCase
	htcu             usecase.HandleTwitterCallbackUseCase
}

func (tc *TwitterController) GetTwitterLoginUrl(c *gin.Context) {
	url, err := tc.gtluu.Execute()
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

func (tc *TwitterController) HandleTwitterCallback(c *gin.Context) {
	oauthToken := c.Query("oauth_token")
	oauthVerifier := c.Query("oauth_verifier")
	jwt, err := tc.htcu.Execute(oauthToken, oauthVerifier)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.Redirect(302, tc.frontRedirectUrl+"?auth_token="+jwt)
}

func NewTwitterController(frontCallbackUrl string, gtluu usecase.GetTwitterLoginUrlUseCase, htcu usecase.HandleTwitterCallbackUseCase) *TwitterController {
	return &TwitterController{
		frontRedirectUrl: frontCallbackUrl,
		gtluu:            gtluu,
		htcu:             htcu,
	}
}
