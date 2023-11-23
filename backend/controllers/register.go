package controllers

import (
	"backend/db"
	"backend/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c echo.Context) error {

	println(c.Request().Header.Get("Content-Type"))
	if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Error when parsing form data"}, Data: nil}
		log.Println(err)

		return c.JSON(http.StatusInternalServerError, res)
	}

	email := strings.ToLower(c.Request().FormValue("input-email"))
	password := c.Request().FormValue("input-password")
	name := c.Request().FormValue("input-name")

	log.Println(email, password, name)

	// validasi user input
	if email == "" || password == "" || name == "" {
		res := ResponseData{IsError: true, Messages: []string{"Email, name, and password is required"}, Data: nil}

		return c.JSON(http.StatusBadRequest, res)
	}

	// cek panjang karakter passowrd
	if len(password) < 8 {
		res := ResponseData{IsError: true, Messages: []string{"Password must have at least 8 characters"}, Data: nil}

		return c.JSON(http.StatusBadRequest, res)
	}

	// generate hash
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	if err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Something went wrong"}, Data: nil}

		return c.JSON(http.StatusInternalServerError, res)
	}

	newUser := models.User{ID: uuid.NewString(), CreatedAt: time.Now(), Name: name, Email: email, Password: string(encryptedPass)}

	dbInstance, err := db.Database()

	if err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Something went wrong"}, Data: nil}

		return c.JSON(http.StatusInternalServerError, res)
	}

	// create user data
	if err := dbInstance.Create(&newUser).Error; err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Email already used"}, Data: nil}

		return c.JSON(http.StatusInternalServerError, res)
	}

	res := ResponseData{IsError: false, Messages: []string{"User created"}, Data: map[string]interface{}{"user": newUser}}

	return c.JSON(http.StatusCreated, res)

}
