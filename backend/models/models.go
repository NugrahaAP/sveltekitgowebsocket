package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
}

type ChatRoom struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Participant  []User         `gorm:"many2many:user_chatroom; not null" json:"participant"`
	Messages     []Message      `gorm:"null" json:"messages"`
	ChatRoomType string         `gorm:"default:personal" json:"chatRoomType"`
}

type Message struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ChatRoomId  string         `json:"chatRoomId"`
	MessageBody string         `gorm:"not null" json:"messageBody"`
	MessageType string         `gorm:"default:message" json:"messageType"`
	MessageLink string         `gorm:"null" json:"messageLink"`
	Sender      User           `gorm:"null;foreignKey:UserId" json:"sender"`
	UserId      string         `gorm:"null" json:"userId"`
}

type GroupChatRoom struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ChatRoom   ChatRoom       `gorm:"foreignKey:ChatRoomId" json:"chatRoom"`
	RoomName   string         `gorm:"null" json:"roomName"`
	RoleAdmin  []User         `gorm:"many2many:group_chat_room_admins;" json:"roleAdmin"`
	ChatRoomId string         `gorm:"null" json:"chatRoomId"`
}
