package main

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/controller"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	apiKey := os.Getenv("API_KEY")
	apiKeySecret := os.Getenv("API_KEY_SECRET")
	callbackUrl := os.Getenv("CALLBACK_URL")
	ts := service.NewTwitterService(apiKey, apiKeySecret, callbackUrl)
	gtluu := usecase.NewGetTwitterLoginUrlUseCase(ts)
	tc := controller.NewTwitterController(gtluu)

	r := gin.Default()
	r.GET("/twitter/login", tc.GetTwitterLoginUrl)

	r.Run(":8080")
}
