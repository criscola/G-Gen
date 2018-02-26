package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"ggen/utils/config"
	"ggen/utils/consts"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

var (
	templates *template.Template
	store     = sessions.NewCookieStore([]byte(config.CookieStoreKey))
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	/** ROUTES **/
	router.GET("/", IndexHandler)
	router.GET("/generator", GeneratorHandler)
	router.GET("/uploads/:filename", UploadGetterHandler)
	router.POST("/generator/imageUpload", ImageUploadHandler)

	http.ListenAndServe(":80", context.ClearHandler(router))

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
	// Loads and parses the template file for this handler
	templates = template.Must(template.ParseFiles("templates/generator/index.tmpl", "templates/base.tmpl"))
	err := templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get or create new session having name consts.SessionName
	session, err := store.Get(r, consts.SessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
}

func UploadGetterHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == "GET" {
		session, err := store.Get(r, consts.SessionName)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(consts.InternalServerError)
			return
		}

		if strings.EqualFold(ps.ByName("filename"), session.Values[consts.SessionImageFilename].(string)) {
			extension := filepath.Ext(ps.ByName("filename"))
			if strings.EqualFold(extension, "jpg") {
				extension = "jpeg"
			}
			w.Header().Set(consts.HttpContentType, "image/"+extension)

			buff := bytes.NewBuffer(ioutil.ReadFile(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string))))

		} else {
			w.WriteHeader(consts.AccessForbidden)
		}
	}
}

func ImageUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// If user didn't upload any image before now, generate new random filename and put it in the session store
		session, err := store.Get(r, consts.SessionName)
		if session.Values[consts.SessionImageFilename] == nil {
			tempImageFilename := TempFileName() + filepath.Ext(handler.Filename)
			session.Values[consts.SessionImageFilename] = tempImageFilename
			sessions.Save(r, w)

		} else if !strings.EqualFold(filepath.Ext(session.Values[consts.SessionImageFilename].(string)), filepath.Ext(handler.Filename)) {
			// Checks if extension of the uploaded file is equals to the extension of the filename in the session store.
			// If filename is NOT equals, rewrite filename in session store and deletes old file from memory
			err = os.Remove(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string)))
			if err != nil {
				fmt.Println(err)
				return
			}
			tempImageFilename := TempFileName() + filepath.Ext(handler.Filename)
			session.Values[consts.SessionImageFilename] = tempImageFilename
			sessions.Save(r, w)
		}
		// Open/Create temp file for image using name stored in the session store
		f, err := os.OpenFile(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string)), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)

		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
		w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(session.Values[consts.SessionImageFilename].(string))))
		w.Write([]byte(session.Values[consts.SessionImageFilename].(string)))
	}

}

// TempFileName generates a temporary filename for use in testing or whatever
func TempFileName() string {
	// Check if file exists in folder... otherwise generate another tempfilename
	files, err := ioutil.ReadDir("./uploads/")
	if err != nil {
		log.Fatal(err)
	}
again:
	randBytes := make([]byte, 16)
	rand.Read(randBytes)

	temp := hex.EncodeToString(randBytes)
	for _, f := range files {
		if strings.EqualFold(f.Name(), temp) {
			goto again
		}
	}
	return temp
}
