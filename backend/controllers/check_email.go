package controllers

import (
	"backend/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CheckEmailHandler(dbInstance *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.QueryParam("email")

		if email == "" {
			res := ResponseData{IsError: true, Messages: []string{"Valid email address is required"}, Data: nil}
			return c.JSON(http.StatusBadRequest, res)
		}

		var user models.User

		if err := dbInstance.Where("email = ?", email).First(&user).Error; err != nil {
			res := ResponseData{IsError: true, Messages: []string{"User with given email did not exists"}, Data: nil}
			return c.JSON(http.StatusNotFound, res)
		}

		res := ResponseData{IsError: false, Messages: []string{"Email exists"}, Data: map[string]interface{}{"userId": user.ID}}
		return c.JSON(http.StatusOK, res)
	}
}
