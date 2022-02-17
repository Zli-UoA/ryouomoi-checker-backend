package main

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/controller"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

func connectDB() (*sqlx.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	var err error
	for i := 0; i < 10; i++ {
		db, err := sqlx.Connect("mysql", dbUser+":"+dbPassword+"@tcp("+dbAddress+":"+dbPort+")/"+dbName)
		if err == nil {
			return db, nil
		}
	}
	return nil, err
}

func main() {
	apiKey := os.Getenv("API_KEY")
	apiKeySecret := os.Getenv("API_KEY_SECRET")
	callbackUrl := os.Getenv("CALLBACK_URL")
	frontRedirectUrl := os.Getenv("FRONT_REDIRECT_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}

	ur := repository.NewUserRepository(db)
	ts := service.NewTwitterService(apiKey, apiKeySecret, callbackUrl)
	ujs := service.NewUserJWTService(jwtSecret)
	gtluu := usecase.NewGetTwitterLoginUrlUseCase(ts)
	htcu := usecase.NewHandleTwitterCallbackUseCase(ts, ujs, ur)
	tc := controller.NewTwitterController(frontRedirectUrl, gtluu, htcu)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/twitter/login", tc.GetTwitterLoginUrl)
	r.GET("/twitter/callback", tc.HandleTwitterCallback)

	r.Run(":8080")
}
