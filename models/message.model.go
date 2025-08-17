package models

import (
	"time"
)

type Message struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"9f1fc230-b312-4c3e-a9ff-5c8e12345678"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	SenderID   string    `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID string    `gorm:"type:uuid;not null" json:"recipient_id"`
	SentAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	ReadAt     time.Time `json:"read_at"`
}
