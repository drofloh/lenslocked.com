package models

import (
	"github.com/jinzhu/gorm"
	// load postgres dialect from gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewServices ...
func NewServices(dialect, connectionInfo string) (*Services, error) {
	db, err := gorm.Open(dialect, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		Image:   NewImageService(),
		db:      db,
	}, nil
}

// Services ...
type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

// Close the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// DestructiveReset drops all table and rebuilds.
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate all tables.
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}
