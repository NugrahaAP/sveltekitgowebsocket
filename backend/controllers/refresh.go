package controllers

import (
	"backend/auth"
	"backend/misc"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func RefreshHandler(c echo.Context) error {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	refreshToken := c.Request().FormValue("refresh-token")

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// cek isi token
	if err != nil || !token.Valid {
		res := ResponseData{IsError: true, Messages: []string{"Invalid refresh token"}, Data: nil}
		return c.JSON(http.StatusUnauthorized, res)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["TokenType"].(string) != "refresh" {
		res := ResponseData{IsError: true, Messages: []string{"Invalid refresh token"}, Data: nil}
		return c.JSON(http.StatusUnauthorized, res)
	}

	log.Println(claims)

	// ambil user id
	userId := claims["userId"].(string)

	// generate access token baru
	newToken, err := auth.GenerateToken(userId)

	if err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Something went wrong"}, Data: nil}
		return c.JSON(http.StatusInternalServerError, res)
	}

	// generate refresh token baru
	newRefreshToken, err := auth.GenerateRefreshToken(userId)

	if err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Something went wrong"}, Data: nil}
		return c.JSON(http.StatusInternalServerError, res)
	}

	// set header dan cookie
	c.Response().Header().Set(echo.HeaderAuthorization, newToken)
	c.SetCookie(misc.CreateCookie("Session", "Authorization: "+newToken))

	res := ResponseData{IsError: false, Messages: []string{"Token refreshed"}, Data: map[string]interface{}{"accessToken": newToken, "refreshToken": newRefreshToken, "userId": userId}}
	return c.JSON(http.StatusCreated, res)
}
