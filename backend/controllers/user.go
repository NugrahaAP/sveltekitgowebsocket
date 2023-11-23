package controllers

import (
	"backend/misc"
	"backend/models"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetUser(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {

	userIdParam := c.QueryParam("userId")
	listParam := c.QueryParam("listUser")

	if userIdParam != "" {
		if !misc.ValidateUUID(userIdParam) {
			res := ResponseData{IsError: true, Messages: []string{"Valid userId is required"}}
			return http.StatusBadRequest, res
		}
		var user models.User
		if err := dbInstance.Where("id = ?", userIdParam).First(&user).Error; err != nil {

			res := ResponseData{IsError: true, Messages: []string{"User did not exists"}, Data: nil}
			return http.StatusNotFound, res
		}

		res = ResponseData{IsError: false, Messages: []string{"Success"}, Data: user}
		return http.StatusOK, res

	} else if listParam != "" && listParam == "true" {

		var users []models.User
		userId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)

		if err := dbInstance.Where("id != ?", userId).Find(&users).Error; err != nil {
			res := ResponseData{IsError: true, Messages: []string{"Something went wrong"}, Data: nil}
			return http.StatusInternalServerError, res
		}

		res = ResponseData{IsError: false, Messages: []string{"Success"}, Data: map[string]interface{}{"users": users}}

		return http.StatusOK, res

	} else {
		requesterUserId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)

		var user models.User
		if err := dbInstance.Where("id = ?", requesterUserId).First(&user).Error; err != nil {

			res := ResponseData{IsError: true, Messages: []string{"User did not exists"}, Data: nil}
			return http.StatusNotFound, res
		}

		res = ResponseData{IsError: false, Messages: []string{"Success"}, Data: user}
		return http.StatusOK, res
	}

}

func UserHandler(dbInstance *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		//Get user where user id is requester userid
		if c.Request().Method == "GET" {
			status, res := GetUser(dbInstance, c)
			return c.JSON(status, res)
		}

		// method not allowed
		res := ResponseData{IsError: true, Messages: []string{"Method is not allowed"}, Data: nil}
		return c.JSON(http.StatusMethodNotAllowed, res)

	}
}
