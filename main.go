package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Zli-UoA/ryouomoi-checker-backend/controller"
	"github.com/Zli-UoA/ryouomoi-checker-backend/repository"
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
	"github.com/Zli-UoA/ryouomoi-checker-backend/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

func connectDB() (*sqlx.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=UTC", dbUser, dbPassword, dbAddress, dbPort, dbName)

	var err error
	for i := 0; i < 10; i++ {
		db, err := sqlx.Connect("mysql", dsn)
		if err == nil {
			return db, nil
		}
		time.Sleep(100 * time.Millisecond)
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
	cr := repository.NewChatRepository(db)

	ts := service.NewTwitterService(apiKey, apiKeySecret, callbackUrl)
	bts := service.NewTwitterService(botApiKey, botApiKeySecret, botCallbackUrl)
	ujs := service.NewUserJWTService(jwtSecret)

	gtluu := usecase.NewGetTwitterLoginUrlUseCase(ts)
	htcu := usecase.NewHandleTwitterCallbackUseCase(ts, ujs, ur)
	fsu := usecase.NewFriendsSearchUseCase(ts, ur)
	slpu := usecase.NewSetLovePointUseCase(6, botUserID, ur, cr, bts)
	gfu := usecase.NewGetFolloweesUseCase(ts, ur)
	fru := usecase.NewGetFollowersUseCase(ts, ur)
	glpu := usecase.NewGetLovePointUsecase(ur)
	glptsu := usecase.NewGetLovePointsUseCase(ur, ts)
	gclu := usecase.NewGetCurrentLover(ur)
	gcedu := usecase.NewGetCoupleElapsedDaysUseCase(ur)
	dclu := usecase.NewDeleteCurrentLover(botUserID, ur, bts)
	gliuu := usecase.NewGetLoggedInUserUseCase(ur, ts)

	tc := controller.NewTwitterController(frontRedirectUrl, gtluu, htcu)
	btc := controller.NewTwitterController(
		"https://example.com/",
		usecase.NewGetTwitterLoginUrlUseCase(bts),
		usecase.NewHandleTwitterCallbackUseCase(bts, ujs, ur),
	)
	fc := controller.NewFriendsController(ujs, fsu, slpu, gfu, fru)
	mc := controller.NewMeController(ujs, glpu, glptsu, gclu, gcedu, dclu, gliuu)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization", "Accept")
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.ExposeHeaders = []string{"Link"}
	r.Use(cors.New(corsConfig))

	r.GET("/twitter/login", tc.GetTwitterLoginUrl)
	r.GET("/twitter/callback", tc.HandleTwitterCallback)

	r.GET("/bot/login", btc.GetTwitterLoginUrl)
	r.GET("/bot/callback", btc.HandleTwitterCallback)

	r.GET("/friends/search", fc.FriendsSearch)
	r.GET("/friends/follower", fc.GetFollowers)
	r.GET("/friends/followee", fc.GetFollowees)
	r.POST("/friends/:id", fc.SetLovePoint)

	r.GET("/me", mc.GetLoggedInUser)
	r.GET("/me/lovers", mc.GetLovePoints)
	r.GET("/me/lovers/:id", mc.GetLovePoint)
	r.GET("/me/lover", mc.GetCurrentLover)
	r.GET("/me/lover/days", mc.GetCoupleElapsedDays)
	r.DELETE("/me/lover", mc.DeleteCurrentLover)

	r.Run(":8080")
}
