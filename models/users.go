package models

import (
	"errors"

	"github.com/jinzhu/gorm"

	// load postgres dialect from gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")
)

// NewUserService ...
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	//defer db.Close()
	return &UserService{
		db: db,
	}, nil

}

// UserService ...
type UserService struct {
	db *gorm.DB
}

// ByID will look up by the ID provided.
// 1 - user, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create will create the provided user and backfill data
// like the ID, created_at, deleted_at fields etc.
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Close the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drops and rebuilds the user table.
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

// User ...
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
