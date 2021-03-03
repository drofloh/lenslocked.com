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

// User represents the user model stored in our database.
// This is used for user accounts, storing both an email address and a password
// so users can log in and gain access to thier content.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UserDB is used to interact with the users database.
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// User to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods used to manipulate and work with the
// user model.
type UserService interface {
	// Authenticate will verify the provided email address and password
	// are correct. If they are correct, the user corresponding to that email
	// will be returned. Otherwise you will receive either:
	// ErrNotFound, ErrInvalidPassword or another error if something goes
	// wrong.
	Authenticate(email, password string) (*User, error)
	UserDB
}

// NewUserService ...
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB: &userValidator{
			UserDB: uv,
		},
	}, nil

}

// UserService ...
type userService struct {
	UserDB
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// ByRemember will hash the remember token and then call ByRemember on the
// subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create will create the provided user and backfill data
// like the ID, created_at, deleted_at fields etc.
func (uv *userValidator) Create(user *User) error {
	if user.Remember != "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	err := runUserValFuncs(user, uv.bcryptPassword, uv.hmacRemember)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update will hash a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user, uv.bcryptPassword, uv.hmacRemember)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with the provided ID
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash a users password with a pre defined pepper and
// bcrypt if the password is not the empty string.
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	//defer db.Close()
	return &userGorm{
		db: db,
	}, nil
}

// Assign userGorm pbject to blank UserDB interface to catch if userGorm
// changes to not satisfy the interface.
var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

// ByID will look up by the ID provided.
// 1 - user, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user with the given email address and returns the user.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token and returns that user.
// This method expects the remember token to already be hashed.
// Errors are the same as ByEmail.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
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
func (us *userService) Authenticate(email, password string) (*User, error) {
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
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update with update the user with all the data in the
// provided user object.
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close the UserService database connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drops and rebuilds the user table.
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
