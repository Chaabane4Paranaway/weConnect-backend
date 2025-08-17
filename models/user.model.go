package models

// Account model info
// @Description User account information
// @Description with user id and username
import "time"

type User struct {
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"9f1fc230-b312-4c3e-a9ff-5c8e12345678"`
	// pseudo        string    `gorm:"unique;not null" json:"pseudo" example:"user@example.com"`
	Pseudo       string    `gorm:"uniqueIndex;not null" json:"pseudo" example:"user"`
	Password     string    `gorm:"not null" json:"password" example:"$2a$12$ABC..."`
	Otp          string    `json:"otp,omitempty" example:"123456"`
	OtpExpiresAt time.Time `json:"otp_expires_at" example:"2025-07-08T14:00:00Z"`
	Verified     bool      `gorm:"default:false" json:"verified" example:"false"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at" example:"2025-07-08T13:55:00Z"`
}

// RegisterRequest model info
// @Description User registration request information
