package main

import (
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
	router.GET("/uploads/:filename", ImageGetHandler)
	router.GET("generator/queue/:queueId", QueueHandler)
	router.POST("/generator/imageUpload", ImagePostHandler)
	router.POST("generator/generate", StartGeneratorJobHandler)
	router.DELETE("/generator/imageRemove", ImageRemoveHandler)

	http.ListenAndServe(":80", context.ClearHandler(router))

}

func IndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templates = template.Must(template.ParseFiles("templates/home/index.tmpl", "templates/base.tmpl"))
	err := templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}

func GeneratorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Loads and parses the template file for this handler
	templates = template.Must(template.ParseFiles("templates/generator/index.tmpl", "templates/base.tmpl"))
	err := templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	// Get or create new session having name consts.SessionName
	session, err := store.Get(r, consts.SessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	// Save it before we write to the response/ from the handler.
	session.Save(r, w)
}

func ImageGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodGet {
		session, err := store.Get(r, consts.SessionName)

		if err != nil {
			panic(err)
			w.WriteHeader(http.StatusInternalServerError)

		}

		if strings.EqualFold(ps.ByName("filename"), session.Values[consts.SessionImageFilename].(string)) {
			extension := strings.TrimPrefix(filepath.Ext(ps.ByName("filename")), ".")
			fmt.Println("ext: " + extension)
			if strings.EqualFold(extension, "jpg") {
				extension = "jpeg"
			}
			w.Header().Set(consts.HttpContentType, "image/"+extension)

			data, err := ioutil.ReadFile(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string)))
			if err != nil {
				panic(err)
			}

			w.Write(data)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}
}

func QueueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodGet {
		// check if there is the parameter
		// if not, return bad request, otherwise:
		// check if parameter matches the corresponding session variable
		// if not, return access forbidden, otherwise:
		// return number (0-100)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ImagePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodPost {
		fmt.Println("ok1")
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("image")
		if err != nil {
			panic(err)

		}
		defer file.Close()

		// Generate random filename
		tempImageFilename := TempFileName() + filepath.Ext(handler.Filename)

		session, err := store.Get(r, consts.SessionName)
		if err != nil {
			panic(err)
		}

		// This can happen if user has reloaded the web page but didn't trigger the proper ImageRemoveHandler, we still have its
		// old session variable, so we can safetly know that we can delete his old image (also, if he is editing multiple files, they
		// will be already loaded in the corresponding web page, so we don't need it anymore)
		if session.Values[consts.SessionImageFilename] != nil {
			if _, err := os.Stat(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string))); !os.IsNotExist(err) {
				err = os.Remove(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string)))
				if err != nil {
					panic(err)

				}
			}
		}

		// Save random filename in the session store
		session.Values[consts.SessionImageFilename] = tempImageFilename
		sessions.Save(r, w)

		// Open/Create temp file for image using name stored in the session store
		f, err := os.OpenFile(filepath.Join("./uploads/", tempImageFilename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)

		}
		defer f.Close()
		_, err = io.Copy(f, file)

		if err != nil {
			panic(err)

		}
		w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
		w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(session.Values[consts.SessionImageFilename].(string))))
		w.Write([]byte(session.Values[consts.SessionImageFilename].(string)))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func StartGeneratorJobHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodPost {
		// Get session
		// Generate id for job and add to session
		// Write in response body the id
		// Start goroutine for the generator job
		// Send back 202 (accepted) status code
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ImageRemoveHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodDelete {
		session, err := store.Get(r, consts.SessionName)
		if err != nil {
			panic(err)

		}
		fmt.Printf(session.Values[consts.SessionImageFilename].(string))
		err = os.Remove(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string)))

		if err != nil {
			// If error... something is wrong and this should be logged
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
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
