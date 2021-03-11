package controllers

import (
	"github.com/drofloh/lenslocked.com/models"
	"github.com/drofloh/lenslocked.com/views"
)

// NewGalleries ...
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),

		gs: gs,
	}
}

// Galleries ...
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
