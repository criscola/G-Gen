package main

import (
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("templates/*"))
)

func main() {
	http.HandleFunc("/", IndexHandler)

	http.ListenAndServe(":8080", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// you access the cached templates with the defined name, not the filename
	err := templates.ExecuteTemplate(w, "homepage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
