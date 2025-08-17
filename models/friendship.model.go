package models

import "time"

type Friendship struct {
	ID        uint      `gorm:"primaryKey"`
	User1ID   string    `gorm:"not null;uniqueIndex:idx_users"`
	User2ID   string    `gorm:"not null;uniqueIndex:idx_users"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
