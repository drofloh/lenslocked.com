package views

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/drofloh/lenslocked.com/context"
)

var (
	// LayoutDir ...
	LayoutDir string = "views/layouts/"
	// TemplateDir ...
	TemplateDir string = "views/"
	// TemplateExt ...
	TemplateExt string = ".gohtml"
)

// NewView ...
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View ...
type View struct {
	Template *template.Template
	Layout   string
}

// ServeHTTP ...
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render used to render a view with the predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. Please email support", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// layoutFiles - returns a slice of strings representing the layoutfiles
// used in out application.
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings representing file paths
// for templates and it prepends the TemplateDir to each string in the
// slice
//
// Eg this input {"home"} would result in the output {"views/home"}
// if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings representing file paths for
// templates and it appends the TemplateExt extention to each string in
// the slice.
//
// Eg ths input {"home"} would result in the output {"home.gohtml"}
// if TemplateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
