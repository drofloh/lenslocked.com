package models

import "strings"

const (
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound modelError = "models: resource not found"

	// ErrPasswordIncorrect is returned when an invalid password is used
	// when attempting to authenticate a user.
	ErrPasswordIncorrect modelError = "models: invalid password provided"

	// ErrEmailRequired is returned when an email address is not provided when
	// creating a user.
	ErrEmailRequired modelError = "models: email address is required"

	// ErrEmailInvalid is returned when an email address provided does not match
	// any of our requirements.
	ErrEmailInvalid modelError = "models: email address is not valid"

	// emailRegex is used to match email addresses. It is not perfect but works
	// well enough for now.
	// emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`

	// ErrEmailTaken is returned when an update or create is attempted with an
	// email address that is already in use.
	ErrEmailTaken modelError = "models: email address is already taken"

	// ErrPasswordTooShort is returned when an update or create is attempted
	// with a password that is less than 8 characters.
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters"

	// ErrPasswordRequired is returned when a create is attempted without a
	// password provided.
	ErrPasswordRequired modelError = "models: password is required"

	// ErrIDInvalid is returned when an invalid ID is provided to a method
	// like Delete.
	ErrIDInvalid privateError = "modles: ID provided was invalid"

	// ErrRememberTooShort is returned when a remember token is not at least 32
	// bytes.
	ErrRememberTooShort privateError = "models: remember token must be at least 32 bytes"

	// ErrRememberRequired is returned when a create or update is attempted
	// without a user remember token hash.
	ErrRememberRequired privateError = "models: remember token is required"

	// ErrUserIDRequired ...
	ErrUserIDRequired privateError = "models: user ID is required"

	// ErrTitleRequired ...
	ErrTitleRequired modelError = "models: title is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
