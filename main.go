package main

import (
	"net/http"

	"github.com/drofloh/lenslocked.com/views"
	"github.com/gorilla/mux"
)

// func handlerFunc(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	if r.URL.Path == "/" {
// 		fmt.Fprint(w, "<h1>welcome to my site</h1>")
// 	} else if r.URL.Path == "/contact" {
// 		fmt.Fprint(w, "To get in touch, please email <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>")
// 	} else {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprintf(w, "<h1>Page Not Found</h1><p>Please email support if the issue continues</p>")
// 	}

// }
var (
	homeView    *views.View
	contactView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
	// fmt.Fprint(w, "<h1>welcome to my site</h1>")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// fmt.Fprint(w, "To get in touch, please email <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>")
	err := contactView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("views/home.gohtml")
	contactView = views.NewView("views/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}
