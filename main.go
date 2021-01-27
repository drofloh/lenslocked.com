package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>welcome to my site</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "To get in touch, please email <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>Page Not Found</h1><p>Please email support if the issue continues</p>")
	}

}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
