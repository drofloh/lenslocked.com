package controllers

import (
	"fmt"
	"net/http"

	"github.com/drofloh/lenslocked.com/views"
)

// NewUsers used to create a new users controller.
// this tempplate with panic if the templates are not parsed correctly,
// and should only be used during initial setup.
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// Users ...
type Users struct {
	NewView *views.View
}

// New used to render the form where a user can creaste a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create is used to process the signup form when a user submits it.
// Used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a temp response.")

}