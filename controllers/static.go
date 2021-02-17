package controllers

import "github.com/drofloh/lenslocked.com/views"

// NewStatic ...
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

// Static ...
type Static struct {
	Home    *views.View
	Contact *views.View
}
