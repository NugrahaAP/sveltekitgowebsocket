package controllers

import (
	"backend/misc"
	"backend/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateChatRoom(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {
	var chatRoom models.ChatRoom
	var participant []models.User

	senderId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)
	receiverId := c.Request().FormValue("input-receiverId")
	roomId := uuid.NewString()

	log.Println("sender ", senderId)
	log.Println("receiver ", receiverId)
	// validasi user input
	if senderId == "" || receiverId == "" || !misc.ValidateUUID(senderId) || !misc.ValidateUUID(receiverId) {

		res := ResponseData{IsError: true, Messages: []string{"Invalid receiverId"}}
		return http.StatusBadRequest, res
	}

	// check apakah sudah ada room dengan participant HANYA 2 user tersebut, jika ada return roomnya
	existingChatRoom := models.ChatRoom{}
	if err := dbInstance.Preload("Participant").Joins("JOIN user_chatroom on user_chatroom.chat_room_id = chat_rooms.id").Where("user_chatroom.user_id IN ? AND chat_rooms.chat_room_type = 'personal'", []string{senderId, receiverId}).Group("chat_rooms.id").Having("COUNT(DISTINCT user_chatroom.user_id) = 2").First(&existingChatRoom).Error; err == nil {

		res := ResponseData{IsError: false, Messages: []string{"Chat room already exists"}, Data: map[string]interface{}{"chatRoom": existingChatRoom}}
		return http.StatusOK, res
	}

	// get user
	if err := dbInstance.Where("id IN ?", []string{senderId, receiverId}).Find(&participant).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Users did not exists"}, Data: nil}
		return http.StatusNotFound, res
	}

	chatRoom = models.ChatRoom{
		ID:           roomId,
		CreatedAt:    time.Now(),
		Participant:  participant,
		Messages:     nil,
		ChatRoomType: "personal",
	}

	// create chatroom
	if err := dbInstance.Create(&chatRoom).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Error when creating chatroom"}, Data: nil}
		return http.StatusInternalServerError, res
	}

	// return chatroom data
	res = ResponseData{IsError: false, Messages: []string{"Chatroom created"}, Data: map[string]interface{}{"chatRoom": chatRoom}}
	return http.StatusCreated, res
}

func GetChatRoom(chatRoomId string, dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {
	var chatRoom models.ChatRoom
	requesterUserId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)
	// get messages in chatroom and requested user in participant
	if err := dbInstance.Preload("Participant").
		Preload("Messages.Sender").
		Joins("Join user_chatroom ON chat_rooms.id = user_chatroom.chat_room_id").
		Where("chat_rooms.id = ? AND user_chatroom.user_id = ?", chatRoomId, requesterUserId).
		First(&chatRoom).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Chat room did not exists"}, Data: nil}
		return http.StatusNotFound, res
	}

	res = ResponseData{IsError: true, Messages: []string{"Success"}, Data: map[string]interface{}{"chatRoom": chatRoom}}
	return http.StatusOK, res
}

func ListChatRoom(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {
	var chatRoom []models.ChatRoom
	userid := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)

	// get list of chatroom
	if err := dbInstance.Preload("Participant").Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("messages.created_at DESC").Preload("Sender")
	}).
		Joins("JOIN user_chatroom ON chat_rooms.id = user_chatroom.chat_room_id").
		Where("chat_rooms.chat_room_type = ? AND user_chatroom.user_id = ?", "personal", userid).
		Find(&chatRoom).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Chat room did not exists"}, Data: map[string]interface{}{"err": err}}
		return http.StatusNotFound, res
	}

	res = ResponseData{IsError: false, Messages: []string{"Success"}, Data: map[string]interface{}{"chatrooms": chatRoom}}
	return http.StatusOK, res
}

// not tested yet!
func LeaveChatRoom(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {
	var UpdatedChatRoom models.ChatRoom
	var UserToRemove models.User

	chatRoomId := c.QueryParam("crid")
	userId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)

	if chatRoomId == "" || !misc.ValidateUUID(userId) || !misc.ValidateUUID(chatRoomId) {
		res := ResponseData{IsError: true, Messages: []string{"Valid"}}
		return http.StatusBadRequest, res
	}

	// get chatroom
	if err := dbInstance.Preload("Participant").Where("id = ?", chatRoomId).First(&UpdatedChatRoom).Error; err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Chat room did not exists"}, Data: nil}
		return http.StatusNotFound, res
	}

	// get user
	if err := dbInstance.Where("id = ?", userId).First(&UserToRemove).Error; err != nil {
		res := ResponseData{IsError: true, Messages: []string{"User did not exists"}, Data: nil}
		return http.StatusNotFound, res
	}

	// pop user dari Participant
	if err := dbInstance.Model(&UpdatedChatRoom).Association("Participant").Delete(&UserToRemove); err != nil {
		res := ResponseData{IsError: true, Messages: []string{"Failed when removing user from chat room"}, Data: nil}

		return http.StatusInternalServerError, res
	}

	// create new message type info
	newMessage := models.Message{ID: uuid.NewString(), CreatedAt: time.Now(), MessageBody: UserToRemove.Name + " is leaving the conversation", MessageType: "info"}

	// update chatroom message
	UpdatedChatRoom.Messages = append(UpdatedChatRoom.Messages, newMessage)
	if err := dbInstance.Save(&UpdatedChatRoom).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Error when adding message to chat room"}, Data: nil}
		return http.StatusInternalServerError, res
	}

	// return message type info "namauser is leaving this conversation"
	res = ResponseData{IsError: false, Messages: []string{"Success"}, Data: map[string]interface{}{"message": newMessage}}
	return http.StatusOK, res
}

func ChatRoomHandler(dbInstance *gorm.DB) echo.HandlerFunc {

	return func(c echo.Context) error {

		// create
		if c.Request().Method == "POST" {

			status, res := CreateChatRoom(dbInstance, c)
			return c.JSON(status, res)
		}

		if c.Request().Method == "GET" {
			// get /chat_room?crid=uuid
			chatRoomId := c.QueryParam("crid")
			if chatRoomId != "" {
				if misc.ValidateUUID(chatRoomId) {
					status, res := GetChatRoom(chatRoomId, dbInstance, c)
					return c.JSON(status, res)
				} else {
					res := ResponseData{IsError: true, Messages: []string{"Valid room id is required"}, Data: nil}
					return c.JSON(http.StatusBadRequest, res)
				}

			} else {
				// list

				status, res := ListChatRoom(dbInstance, c)
				return c.JSON(status, res)
			}
		}

		// delete/chat_room?crid=uuid
		if c.Request().Method == "DELETE" {

			status, res := LeaveChatRoom(dbInstance, c)
			return c.JSON(status, res)
		}

		// method is not allowed
		res := ResponseData{IsError: true, Messages: []string{"Method is not allowed"}, Data: nil}
		return c.JSON(http.StatusMethodNotAllowed, res)
	}
}
