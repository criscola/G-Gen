package main

import (
	"html/template"
	"net/http"
)

var (
	templates *template.Template
)

func main() {
	/** ROUTES **/
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(""))))

	http.ListenAndServe(":8080", nil)
}
