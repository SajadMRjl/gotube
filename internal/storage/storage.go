package storage

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStorage struct {
	db *gorm.DB
}

func NewGormStorage(connStr string) (*GormStorage, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&Track{}, &User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &GormStorage{db: db}, nil
}

func (s *GormStorage) CreateTrack(ctx context.Context, track *Track) error {
	return s.db.WithContext(ctx).Create(track).Error
}

func (s *GormStorage) GetTrackBySpotifyID(ctx context.Context, spotifyID string) (*Track, error) {
	var track Track
	err := s.db.WithContext(ctx).Where("spotify_id = ?", spotifyID).First(&track).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &track, err
}

func (s *GormStorage) GetTrackByTelegramFileID(ctx context.Context, fileID string) (*Track, error) {
	var track Track
	err := s.db.WithContext(ctx).Where("telegram_file_id = ?", fileID).First(&track).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &track, err
}

func (s *GormStorage) UpdateTrack(ctx context.Context, track *Track) error {
	return s.db.WithContext(ctx).Save(track).Error
}

func (s *GormStorage) DeleteTrack(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&Track{}, id).Error
}

func (s *GormStorage) CreateUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *GormStorage) GetUser(ctx context.Context, telegramID int64) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (s *GormStorage) UpdateUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

func (s *GormStorage) DeleteUser(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&User{}, id).Error
}

func (s *GormStorage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}