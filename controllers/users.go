package controllers

import (
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

// New ...
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}
