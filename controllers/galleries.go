package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drofloh/lenslocked.com/context"
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

// GalleryForm ...
type GalleryForm struct {
	Title string `schema:"title"`
}

// Create ...
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/logn", http.StatusFound)
		return
	}
	fmt.Println("Create got the user:", user)
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
