package storage

import (
	"time"

	"gorm.io/gorm"
)

type Track struct {
	gorm.Model
	SpotifyID         string `gorm:"uniqueIndex;not null"`
	Title             string `gorm:"not null"`
	Artist            string `gorm:"not null"`
	Duration          int    // in seconds
	TelegramFileID    string `gorm:"index"`
	TelegramMessageID int
	ChannelID         int64
}

type User struct {
	gorm.Model
	TelegramID    int64 `gorm:"uniqueIndex;not null"`
	Username      string
	FirstName     string
	LastName      string
	LanguageCode  string
	LastMessageAt time.Time
}
