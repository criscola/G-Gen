package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var (
	templates *template.Template
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	/** ROUTES **/
	router.GET("/", IndexHandler)
	router.GET("/generator", GeneratorHandler)
	router.POST("/generator/imageUpload", ImageUploadHandler)

	http.ListenAndServe(":8080", router)
}

func IndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templates = template.Must(template.ParseFiles("templates/home/index.tmpl", "templates/base.tmpl"))
	err := templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GeneratorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templates = template.Must(template.ParseFiles("templates/generator/index.tmpl", "templates/base.tmpl"))
	err := templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ImageUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("method:", r.Method)
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
