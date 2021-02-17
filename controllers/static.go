package controllers

import "github.com/drofloh/lenslocked.com/views"

// NewStatic ...
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "views/static/contact.gohtml"),
	}
}

// Static ...
type Static struct {
	Home    *views.View
	Contact *views.View
}
