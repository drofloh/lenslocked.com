package models

import (
	"errors"

	"github.com/drofloh/lenslocked.com/hash"
	"github.com/drofloh/lenslocked.com/rand"
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

	// ErrInvalidPassword is returned when an invalid password is used
	// when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: invalid password provided")
)

const userPwPepper = "Nhfsf632dfsg"
const hmacSecretKey = "secret-key"

// NewUserService ...
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	//defer db.Close()
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil

}

// UserService ...
type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
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

// ByRemember looks up a user with the given remember token and returns that user.
// This method will handle hashing the token for us.
// Errors are the same as ByEmail.
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	err := first(us.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Authenticate can be used to authenticate a user with the provided email
// and password.
// If the email address provided is invalid, this will return
// 		nil, ErrNotFound
// If the password provided is invalid, this will return
// 		nil, ErrInvalidPassword
// If the email and password are both valid, this will return
// 		user, nil
// Otherwise if another error is encountered this will return
// 		nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
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

	if user.Remember != "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)

	return us.db.Create(user).Error
}

// Update with update the user with all the data in the
// provided user object.
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
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
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}
