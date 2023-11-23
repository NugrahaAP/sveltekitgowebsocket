package main

import (
	"backend/controllers"
	"backend/db"
	"backend/websocket"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error when loading .env file: ", err)
	}

	var secretKey = []byte(os.Getenv("SECRET_KEY"))

	e := echo.New()
	e.Debug = true

	wsServer := websocket.NewWSServer()

	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	//
	// 		log.Println("Authorization Header:", c.Request().Header.Get("Authorization"))
	// 		return next(c)
	// 	}
	// })

	dbInstance, err := db.Database()
	if err != nil {
		panic(err)
	}

	e.GET("/health_check", controllers.HealthCheckHandler)

	authGroup := e.Group("/auth")
	authGroup.POST("/login", controllers.LoginHandler)
	authGroup.POST("/register", controllers.RegisterHandler)
	authGroup.POST("/refresh", controllers.RefreshHandler)

	e.GET("/ws/:chat_room_id", wsServer.HandleWS)

	apiGroup := e.Group("/backend/api/v1")
	apiGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    secretKey,
		TokenLookup:   "header:Authorization",
		SigningMethod: "HS256",
	}))

	// check token expiry
	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenExpiry := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["StandarClaim"].(map[string]interface{})["exp"].(float64)
			log.Println(tokenExpiry)

			if time.Now().Unix() > int64(tokenExpiry) {
				log.Println("is expired: ", time.Now().Unix() > int64(tokenExpiry))
				res := controllers.ResponseData{IsError: true, Messages: []string{"Access token is expired"}, Data: nil}
				return c.JSON(http.StatusUnauthorized, res)
			}

			return next(c)
		}
	})

	apiGroup.GET("/rahasia", func(c echo.Context) error {
		// userid := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)

		return c.JSON(http.StatusOK, controllers.ResponseData{IsError: false, Messages: []string{"Access token is valid"}, Data: nil})
	})

	apiGroup.Any("/message", controllers.MessageHandler(dbInstance))
	apiGroup.Any("/group_chat_room", controllers.GroupChatRoomHandler(dbInstance))
	apiGroup.Any("/chat_room", controllers.ChatRoomHandler(dbInstance))
	apiGroup.Any("/user", controllers.UserHandler(dbInstance))
	apiGroup.GET("/checkEmail", controllers.CheckEmailHandler(dbInstance))

	e.Logger.Fatal(e.Start(":1437"))

}
