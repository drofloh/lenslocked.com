package models

import "github.com/jinzhu/gorm"

// NewServices ...
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User: NewUserService(db),
	}, nil
}

// Services ...
type Services struct {
	Gallery GalleryService
	User    UserService
}
