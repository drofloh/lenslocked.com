package main

import (
	"fmt"
	"net/http"
	"text/template"

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
var homeTemplate *template.Template

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
	// fmt.Fprint(w, "<h1>welcome to my site</h1>")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "To get in touch, please email <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>")
}

func main() {
	var err error
	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}
