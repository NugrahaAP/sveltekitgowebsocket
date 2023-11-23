package controllers

import (
	"backend/misc"
	"backend/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateMessage(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {

	var newMessage models.Message

	chatRoomId := c.Request().FormValue("input-chatRoomId")
	senderId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)
	messageType := c.Request().FormValue("input-messageType")
	messageBody := c.Request().FormValue("input-messageBody")
	messageLink := c.Request().FormValue("input-messageLink")

	// validasi user input
	if chatRoomId == "" || senderId == "" || messageType == "" || messageBody == "" {

		res := ResponseData{IsError: true, Messages: []string{"ChatRoomId, SenderId, MessageType is required"}, Data: nil}
		return http.StatusBadRequest, res
	}

	// create message
	newMessage = models.Message{ID: uuid.NewString(), CreatedAt: time.Now(), MessageBody: messageBody, MessageType: messageType, MessageLink: messageLink, UserId: senderId, ChatRoomId: chatRoomId}
	if err := dbInstance.Create(&newMessage).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Error when creating message"}, Data: nil}
		return http.StatusBadRequest, res
	}

	// get chatroom
	var chatRoom models.ChatRoom
	if err := dbInstance.Where("id = ?", chatRoomId).First(&chatRoom).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Chat room did not exists"}, Data: nil}
		return http.StatusNotFound, res
	}

	chatRoom.Messages = append(chatRoom.Messages, newMessage)

	if err := dbInstance.Save(&chatRoom).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Error when updating chat room"}, Data: nil}
		return http.StatusInternalServerError, res
	}
	newMessage.Sender.ID = senderId

	// return newmessage
	res = ResponseData{IsError: false, Messages: []string{"Message created"}, Data: map[string]interface{}{"message": newMessage}}
	return http.StatusCreated, res
}

// not tested yet!
func DeleteMessage(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {

	messageId := c.QueryParam("msgid")
	userId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)

	// validasi message id
	if messageId == "" || !misc.ValidateUUID(messageId) {

		res := ResponseData{IsError: true, Messages: []string{"Valid messageId is required"}, Data: nil}
		return http.StatusBadRequest, res
	}

	// delete message
	if err := dbInstance.Where("id = ? and sender_id", messageId, userId).Delete(&models.Message{}).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Message did not exists"}, Data: nil}
		return http.StatusNotFound, res
	}

	res = ResponseData{IsError: false, Messages: []string{"Message deleted"}, Data: nil}
	return http.StatusOK, res
}

func UpdateMessage(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {
	var message models.Message

	messageId := c.Request().FormValue("input-messageId")
	messsageBody := c.Request().FormValue("input-messageBody")

	if messageId == "" || messsageBody == "" || !misc.ValidateUUID(messageId) {

		res := ResponseData{IsError: true, Messages: []string{"valid message id and message body is required"}, Data: nil}
		return http.StatusBadRequest, res
	}

	// update message
	if err := dbInstance.Where("id = ?", messageId).First(&message).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Message did not exists"}, Data: nil}
		return http.StatusNotFound, res
	}

	message.MessageBody = messsageBody
	dbInstance.Save(&message)

	res = ResponseData{IsError: false, Messages: []string{"Message updated"}, Data: map[string]interface{}{"message": message}}
	return http.StatusOK, res
}

func MessageHandler(dbInstance *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		// create message
		if c.Request().Method == "POST" {

			status, res := CreateMessage(dbInstance, c)
			return c.JSON(status, res)
		}

		// delete /messsage?msgid=uuid
		if c.Request().Method == "DELETE" {

			status, res := DeleteMessage(dbInstance, c)
			return c.JSON(status, res)
		}

		//update message
		if c.Request().Method == "PATCH" {

			status, res := UpdateMessage(dbInstance, c)
			return c.JSON(status, res)
		}

		// method is not allowed
		res := ResponseData{IsError: true, Messages: []string{"Method is not allowed"}, Data: nil}
		return c.JSON(http.StatusMethodNotAllowed, res)
	}
}
