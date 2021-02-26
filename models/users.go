package models

import (
	"errors"

	"github.com/jinzhu/gorm"

	// load postgres dialect from gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to a method
	// like Delete
	ErrInvalidID = errors.New("modles: ID provided was invalid")
)

const userPwPepper = "Nhfsf632dfsg"

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
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user with the given email address and returns the user.
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// first will query using the provided gorm.DB and it will get the
// first item returned and place it into dst. If nothing is found
// in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create will create the provided user and backfill data
// like the ID, created_at, deleted_at fields etc.
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

// Update with update the user with all the data in the
// provided user object.
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// Close the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drops and rebuilds the user table.
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// User ...
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}
