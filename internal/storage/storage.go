package storage

import (
	"context"
	"errors"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStorage struct {
	DB *gorm.DB
}

func NewGormStorage(connStr string) (*GormStorage, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &GormStorage{DB: db}, nil
}

func (s *GormStorage) GetOrCreateUser(ctx context.Context, tgUser *tgbotapi.User) (*User, error) {
	var user User
	result := s.DB.WithContext(ctx).Where("telegram_id = ?", tgUser.ID).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new user
		newUser := User{
			TelegramID:    tgUser.ID,
			Username:      tgUser.UserName,
			FirstName:     tgUser.FirstName,
			LastName:      tgUser.LastName,
			LanguageCode:  tgUser.LanguageCode,
			LastMessageAt: time.Now(),
		}
		if err := s.DB.WithContext(ctx).Create(&newUser).Error; err != nil {
			return nil, err
		}
		return &newUser, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	// Update existing user
	user.LastMessageAt = time.Now()
	if err := s.DB.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
