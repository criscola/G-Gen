package main

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"ggen/utils/config"
	"ggen/utils/consts"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"fmt"
)

var (
	templates 		*template.Template
	jobCompletion 	= make(chan *GeneratorJob)
	store     		= sessions.NewCookieStore([]byte(config.CookieStoreKey))

)

func main() {
	router := httprouter.New()
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	/** ROUTES **/
	router.GET("/", IndexHandler)
	router.GET("/generator", GeneratorHandler)
	router.GET("/uploads/:"+consts.RequestFilename, ImageGetHandler)
	router.GET("/generator/queue/:"+consts.RequestQueueId, QueueHandler)
	router.GET("/generator/output/:"+consts.RequestOutputId, OutputHandler)
	router.POST("/generator/imageUpload", ImagePostHandler)
	router.POST("/generator/generate", StartGeneratorJobHandler)
	router.DELETE("/generator/imageRemove", ImageRemoveHandler)

	http.ListenAndServe(":"+config.ServerPort, context.ClearHandler(router))

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if strings.EqualFold(ps.ByName(consts.RequestFilename), session.Values[consts.SessionImageFilename].(string)) {
			extension := strings.TrimPrefix(filepath.Ext(ps.ByName(consts.RequestFilename)), ".")
			if strings.EqualFold(extension, "jpg") {
				extension = "jpeg"
			}
			w.Header().Set(consts.HttpContentType, "image/"+extension)

			data, err := ioutil.ReadFile(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string)))
			checkError(err)

			w.Write(data)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}
}

func QueueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodGet {
		session, err := store.Get(r, consts.SessionName)
		checkError(err)

		id := ps.ByName(consts.RequestQueueId)

		jobs := session.Values[consts.SessionGeneratorJob].(map[string]*GeneratorJob)

		// check if there is the parameter and session
		if id != "" && jobs[id] != nil && jobs[id].Completion != 100 {
			// if there is a job with key as queueId
			selectedJob := <- jobCompletion
			fmt.Println("from QueueHand ", selectedJob.Completion)
			jobs[id] = selectedJob
			sessions.Save(r, w)
			// If the user sending the request actually owns this resource, return the result
			if selectedJob.Id == id {
				// Return % of completion
				temp := strconv.Itoa(selectedJob.Completion)
				w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
				w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(temp)))
				w.Write([]byte(temp))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func OutputHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ImagePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile(consts.FormImage)
		checkError(err)

		defer file.Close()

		// Generate random filename
		tempImageFilename := TempFileName() + filepath.Ext(handler.Filename)

		session, err := store.Get(r, consts.SessionName)
		checkError(err)

		// This can happen if user has reloaded the web page but didn't trigger the proper ImageRemoveHandler, we still have its
		// old session variable, so we can safely know that we can delete his old image (also, if he is editing multiple files, they
		// will be already loaded in the corresponding web page, so we don't need it anymore)
		if session.Values[consts.SessionImageFilename] != nil {
			if _, err := os.Stat(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string))); !os.IsNotExist(err) {
				err = os.Remove(filepath.Join("./uploads/", session.Values[consts.SessionImageFilename].(string)))
				checkError(err)
			}
		}

		// Save random filename in the session store
		session.Values[consts.SessionImageFilename] = tempImageFilename
		sessions.Save(r, w)

		// Open/Create temp file for image using name stored in the session store
		f, err := os.OpenFile(filepath.Join("./uploads/", tempImageFilename), os.O_WRONLY|os.O_CREATE, 0666)
		checkError(err)

		defer f.Close()
		_, err = io.Copy(f, file)
		checkError(err)

		w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
		w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(tempImageFilename)))
		w.Write([]byte(tempImageFilename))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func StartGeneratorJobHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodPost {

		// Start goroutine for the generator job
		// Add goroutine to map holding all jobs with keys as id
		// Send back 202 (accepted) status code
		// Get session
		session, err := store.Get(r, consts.SessionName)

		// Generate id for job
		id := GetRandomString()

		// Generate id for job and add to session struct
		var jobs map[string]*GeneratorJob

		if session.Values[consts.SessionGeneratorJob] != nil {
			// Append new job to the list
			jobs = session.Values[consts.SessionGeneratorJob].(map[string]*GeneratorJob)
		} else {
			jobs = make(map[string]*GeneratorJob)
		}

		jobs[id] = &GeneratorJob{id, time.Now().Unix(), 0, 0}

		session.Values[consts.SessionGeneratorJob] = jobs

		scaleFactor, err := strconv.Atoi(r.FormValue(consts.FormScaleFactor))
		checkError(err)
		//modelThickness, err := strconv.Atoi(r.FormValue(consts.FormModelThickness))
		//checkError(err)
		travelSpeed, err := strconv.Atoi(r.FormValue(consts.FormTravelSpeed))
		checkError(err)

		generationParams := GeneratorParams{
			ScaleFactor: scaleFactor,
			ModelThickness: 4,
			TravelSpeed: travelSpeed,
		}

		go StartGeneratorJob(r, jobs[id], &generationParams, jobCompletion)

		session.Save(r, w)

		// Write in response body the id
		w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
		w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(id)))
		w.Write([]byte(id))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ImageRemoveHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodDelete {
		session, err := store.Get(r, consts.SessionName)
		checkError(err)
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
	checkError(err)
again:
	temp := GetRandomString()

	for _, f := range files {
		if strings.EqualFold(f.Name(), temp) {
			goto again
		}
	}
	return temp
}

func GetRandomString() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	gob.Register(map[string]*GeneratorJob{})
}
