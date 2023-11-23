package controllers

import (
	"backend/auth"
	"backend/db"
	"backend/misc"
	"backend/models"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c echo.Context) error {

	var user models.User

	email := strings.ToLower(c.Request().FormValue("input-email"))
	password := c.Request().FormValue("input-password")

	// kalo email atau password empty string
	if email == "" || password == "" {
		res := ResponseData{IsError: true, Messages: []string{"Email and Password is required"}, Data: nil}
		return c.JSON(http.StatusBadRequest, res)
	}

	// kalo panjang password kecil dari 8
	if len(password) < 8 {
		res := ResponseData{IsError: true, Messages: []string{"Password must have at least 8 characters"}, Data: nil}

		return c.JSON(http.StatusBadRequest, res)
	}

	// db connection
	dbInstance, err := db.Database()

	if err != nil {
		log.Fatal("Error when creating connection to db: ", err)

		res := ResponseData{IsError: true, Messages: []string{"Something went wrong"}, Data: nil}
		return c.JSON(http.StatusInternalServerError, res)
	}

	// get user
	if err := dbInstance.Where("email = ?", email).First(&user).Error; err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Invalid username/password"}, Data: nil}
		return c.JSON(http.StatusNotFound, res)
	}

	// bandingkan password
	comparePass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if comparePass != nil {
		res := ResponseData{IsError: true, Messages: []string{"Invalid username/password"}, Data: nil}
		return c.JSON(http.StatusBadRequest, res)
	}

	// generate JWT Token and Refrest Token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Something went wrong"}, Data: nil}
		return c.JSON(http.StatusInternalServerError, res)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Something went wrong"}, Data: nil}
		return c.JSON(http.StatusInternalServerError, res)
	}

	c.Response().Header().Set(echo.HeaderAuthorization, token)
	c.SetCookie(misc.CreateCookie("Session", "Authorization: "+token))

	res := ResponseData{IsError: false, Messages: []string{"Success"}, Data: map[string]interface{}{"user": user, "accessToken": token, "refreshToken": refreshToken}}
	return c.JSON(http.StatusOK, res)

}
