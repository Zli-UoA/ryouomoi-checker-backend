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
	"strconv"
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
	botApiKey := os.Getenv("BOT_API_KEY")
	botApiKeySecret := os.Getenv("BOT_API_KEY_SECRET")
	botCallbackUrl := os.Getenv("BOT_CALLBACK_URL")
	botUserID, _ := strconv.ParseInt(os.Getenv("BOT_USER_ID"), 10, 64)
	frontRedirectUrl := os.Getenv("FRONT_REDIRECT_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}

	ur := repository.NewUserRepository(db)

	ts := service.NewTwitterService(apiKey, apiKeySecret, callbackUrl)
	bts := service.NewTwitterService(botApiKey, botApiKeySecret, botCallbackUrl)
	ujs := service.NewUserJWTService(jwtSecret)

	gtluu := usecase.NewGetTwitterLoginUrlUseCase(ts)
	htcu := usecase.NewHandleTwitterCallbackUseCase(ts, ujs, ur)
	fsu := usecase.NewFriendsSearchUseCase(ts, ur)
	slpu := usecase.NewSetLovePointUseCase(5, botUserID, ur, bts)

	tc := controller.NewTwitterController(frontRedirectUrl, gtluu, htcu)
	btc := controller.NewTwitterController(
		"https://example.com/",
		usecase.NewGetTwitterLoginUrlUseCase(bts),
		usecase.NewHandleTwitterCallbackUseCase(bts, ujs, ur),
	)
	fc := controller.NewFriendsController(ujs, fsu, slpu)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/twitter/login", tc.GetTwitterLoginUrl)
	r.GET("/twitter/callback", tc.HandleTwitterCallback)

	r.GET("/bot/login", btc.GetTwitterLoginUrl)
	r.GET("/bot/callback", btc.HandleTwitterCallback)

	r.GET("/friends/search", fc.FriendsSearch)
	r.POST("/friends/:id", fc.SetLovePoint)

	r.Run(":8080")
}
