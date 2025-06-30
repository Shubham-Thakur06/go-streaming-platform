package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Media     []Media   `json:"media,omitempty" gorm:"foreignKey:UserID"`
}

type Media struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description"`
	Filename     string    `json:"filename" gorm:"not null"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type" gorm:"not null"` // "mp3", "mp4", etc.
	Genre        string    `json:"genre"`
	Tags         string    `json:"tags"`     // comma-separated tags
	Duration     int       `json:"duration"` // in seconds
	StorageURL   string    `json:"storage_url" gorm:"not null"`
	ThumbnailURL string    `json:"thumbnail_url"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User         User      `json:"user,omitempty"`
	IsPublic     bool      `json:"is_public" gorm:"default:true"`
	ViewCount    int       `json:"view_count" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserActivity struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	MediaID   uuid.UUID `json:"media_id" gorm:"type:uuid;not null"`
	Action    string    `json:"action" gorm:"not null"` // "play", "pause", "stop", "share"
	Duration  int       `json:"duration"`               // seconds watched
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user,omitempty"`
	Media     Media     `json:"media,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return nil
}

func (media *Media) BeforeCreate(tx *gorm.DB) error {
	if media.ID == uuid.Nil {
		media.ID = uuid.New()
	}
	return nil
}

func (activity *UserActivity) BeforeCreate(tx *gorm.DB) error {
	if activity.ID == uuid.Nil {
		activity.ID = uuid.New()
	}
	return nil
}
