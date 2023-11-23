package controllers

import (
	"backend/misc"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateGroupChatRoom(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {
	var newGroupChatRoom models.GroupChatRoom
	var newChatRoom models.ChatRoom
	var participant []models.User
	var listOfAdmin []models.User
	var addAdmin models.User

	// get list of participant user_id
	requesterUserId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)

	//formdata json string {participantId:["uuid","uuid"]}
	inputParticipantUserId := c.Request().FormValue("input-participant")
	log.Println(inputParticipantUserId)
	roomName := c.Request().FormValue("input-roomName")

	if inputParticipantUserId == "" || roomName == "" {

		res := ResponseData{IsError: true, Messages: []string{"List of userId and room name is required"}, Data: nil}
		return http.StatusBadRequest, res
	}

	var userIds map[string][]string

	if err := json.Unmarshal([]byte(inputParticipantUserId), &userIds); err != nil {

		res := ResponseData{IsError: true, Messages: []string{"List of userId is not valid JSON string"}, Data: nil}
		return http.StatusBadRequest, res
	}

	// validasi userid is valid uuid string
	for _, userId := range userIds["participantId"] {
		if !misc.ValidateUUID(userId) {

			res := ResponseData{IsError: true, Messages: []string{"Invalid userId"}, Data: nil}
			return http.StatusBadRequest, res
		}
		participant = append(participant, models.User{ID: userId})
	}

	addAdmin = models.User{ID: requesterUserId}
	listOfAdmin = append(listOfAdmin, addAdmin)
	newChatRoom = models.ChatRoom{ID: uuid.NewString(), CreatedAt: time.Now(), Participant: participant, ChatRoomType: "group"}
	newGroupChatRoom = models.GroupChatRoom{ID: uuid.NewString(), CreatedAt: time.Now(), ChatRoom: newChatRoom, RoomName: roomName, RoleAdmin: listOfAdmin, ChatRoomId: newChatRoom.ID}

	if err := dbInstance.Create(&newGroupChatRoom).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Error when creating group chat room "}, Data: nil}
		return http.StatusInternalServerError, res
	}

	var returnGroupChatRoom models.GroupChatRoom
	if err := dbInstance.Preload("ChatRoom").Preload("RoleAdmin").Where("id = ?", newGroupChatRoom.ID).First(&returnGroupChatRoom).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Error when retrieveing group chat room"}, Data: nil}
		return http.StatusInternalServerError, res
	}

	log.Println("group created")

	res = ResponseData{IsError: false, Messages: []string{"Group chat room created"}, Data: map[string]interface{}{"groupChatRoom": returnGroupChatRoom}}
	return http.StatusCreated, res
}

func GetGroupChatRoom(groupChatRoomId string, dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {

	requesterUserId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)
	var groupChatRoom models.GroupChatRoom

	if !misc.ValidateUUID(groupChatRoomId) {

		res := ResponseData{IsError: true, Messages: []string{"Invalid group chat room id"}, Data: nil}
		return http.StatusBadRequest, res
	}

	// get group chat room where id = gcrid and requester user in participant
	if err := dbInstance.Joins("JOIN chat_rooms ON group_chat_rooms.chat_room_id = chat_rooms.id").
		Joins("JOIN user_chatroom ON chat_rooms.id = user_chatroom.chat_room_id").
		Preload("ChatRoom.Participant").
		Preload("ChatRoom.Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("messages.created_at DESC").Limit(1).Preload("Sender")

		}).
		First(&groupChatRoom, "group_chat_rooms.id = ? AND user_chatroom.user_id = ?", groupChatRoomId, requesterUserId).Error; err != nil {

		res := ResponseData{IsError: false, Messages: []string{"Group chat room does not exist or user is not a participant"}, Data: map[string]interface{}{"err": err}}

		return http.StatusNotFound, res
	}

	res = ResponseData{IsError: false, Messages: []string{"Success"}, Data: map[string]interface{}{"groupChatRoom": groupChatRoom}}
	return http.StatusOK, res
}

func ListGroupChatRoom(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {

	var listOfGroupChatRoom []models.GroupChatRoom
	requesterUserId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)

	// get group chat where user in participant
	if err := dbInstance.Preload("ChatRoom.Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("messages.created_at DESC").Preload("Sender")
	}).Preload("RoleAdmin").
		Joins("JOIN chat_rooms ON group_chat_rooms.chat_room_id = chat_rooms.id").
		Joins("JOIN user_chatroom ON chat_rooms.id = user_chatroom.chat_room_id").
		Preload("ChatRoom.Participant").
		Where("chat_rooms.chat_room_type = ? AND user_chatroom.user_id = ?", "group", requesterUserId).
		Find(&listOfGroupChatRoom).Error; err != nil {

		res := ResponseData{IsError: false, Messages: []string{"Group chat room does not exist or user is not a participant"}, Data: map[string]interface{}{"err": err}}
		return http.StatusNotFound, res
	}

	res = ResponseData{IsError: true, Messages: []string{"Success"}, Data: map[string]interface{}{"groupChatRooms": listOfGroupChatRoom}}
	return http.StatusOK, res
}

// not tested yet!
func DeleteGroupChatRoom(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {
	// delete where id = group chat room id and userid is in role admin

	requesterUserId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)
	groupChatRoomId := c.QueryParam("gcrid")
	if groupChatRoomId == "" || !misc.ValidateUUID(groupChatRoomId) {

		res := ResponseData{IsError: true, Messages: []string{"Valid group chat room id is required"}, Data: nil}
		return http.StatusBadRequest, res
	}

	if err := dbInstance.Where("id = ? AND ? IN (SELECT user_id FROM group_chat_room_admins WHERE group_chat_room_id = group_chat_rooms.id)", groupChatRoomId, requesterUserId).Delete(&models.GroupChatRoom{}).Error; err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Error when trying to delete group chat room"}, Data: nil}
		return http.StatusInternalServerError, res
	}

	res = ResponseData{IsError: false, Messages: []string{"Success"}, Data: nil}
	return http.StatusOK, res
}

func InsertAdmin(dbInstance *gorm.DB, c echo.Context) (httpStatus int, res ResponseData) {

	inputNewAdminId := c.Request().FormValue("input-newAdminId")
	requesterUserId := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["UserId"].(string)
	groupChatRoomId := c.Request().FormValue("input-groupChatRoomId")

	if inputNewAdminId == "" || groupChatRoomId == "" || !misc.ValidateUUID(inputNewAdminId) || !misc.ValidateUUID(groupChatRoomId) {

		res := ResponseData{IsError: true, Messages: []string{"Valid user id and group chat room id is required"}}
		return http.StatusBadRequest, res
	}

	// updatedGroupChatRoom.RoleAdmin
	if err := dbInstance.Model(&models.GroupChatRoom{}).Where("id = ? AND ? IN (SELECT user_id FROM group_chat_room_admins WHERE group_chat_room_id = group_chat_rooms.id)", groupChatRoomId, requesterUserId).Association("RoleAdmin").Append(&models.User{ID: inputNewAdminId}); err != nil {

		res := ResponseData{IsError: true, Messages: []string{"Error when updating group chat room admin"}, Data: nil}
		return http.StatusInternalServerError, res
	}

	res = ResponseData{IsError: false, Messages: []string{"Successfully add user to group chat admin"}, Data: nil}
	return http.StatusOK, res
}

func GroupChatRoomHandler(dbInstance *gorm.DB) echo.HandlerFunc {

	return func(c echo.Context) error {
		// create
		if c.Request().Method == "POST" {

			status, res := CreateGroupChatRoom(dbInstance, c)
			return c.JSON(status, res)

		}

		if c.Request().Method == "GET" {

			// get group_chat_room?gcrid=uuid
			groupChatRoomId := c.QueryParam("gcrid")
			if groupChatRoomId != "" {
				log.Println("get group")

				status, res := GetGroupChatRoom(groupChatRoomId, dbInstance, c)
				return c.JSON(status, res)
			} else {

				log.Println("list group")
				// list group_chat_room where current user in participant
				status, res := ListGroupChatRoom(dbInstance, c)
				return c.JSON(status, res)
			}
		}

		if c.Request().Method == "DELETE" {

			status, res := DeleteGroupChatRoom(dbInstance, c)
			return c.JSON(status, res)
		}

		if c.Request().Method == "PATCH" {

			status, res := InsertAdmin(dbInstance, c)
			return c.JSON(status, res)
		}

		// method is not allowed
		res := ResponseData{IsError: true, Messages: []string{"Method is not allowed"}, Data: nil}
		return c.JSON(http.StatusMethodNotAllowed, res)
	}
}
