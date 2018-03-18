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

	"fmt"
		"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

var (
	templates     *template.Template
	jobComplPerc  = 0
	jobCompletion = make(chan *GeneratorJob, 100)
	store         = sessions.NewCookieStore([]byte(config.CookieStoreKey))
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	/** ROUTES **/
	router.GET("/", IndexHandler)
	router.GET("/generator", GeneratorHandler)
	router.GET("/uploads/:"+consts.RequestFilename, ImageGetHandler)
	router.GET("/generator/queue/:"+consts.RequestJobId, QueueHandler)
	router.GET("/generator/outputs/:"+consts.RequestJobId, OutputHandler)
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

		// Get the request id
		id := ps.ByName(consts.RequestJobId)
		// Get the job map from session
		jobs := session.Values[consts.SessionGeneratorJob].(map[string]*GeneratorJob)

		// check if id provided by the HTTP request isn't missing and that the job[id] exists
		if id != "" && jobs[id] != nil {
			// Get the last selected job supplied by the jobCompletion channel
			var selectedJob *GeneratorJob
			percCount := len(jobCompletion)
			fmt.Println("c: ", percCount)

			if percCount > 0 {
				for i := 0; i < percCount; i++ {
					selectedJob = <-jobCompletion
					fmt.Println("\nSelectedjob completion: " + strconv.Itoa(selectedJob.Completion))
				}
				jobComplPerc = selectedJob.Completion

				// Reassign the new fresh selected job
				jobs[id] = selectedJob
				sessions.Save(r, w)

				// Return % of completion
				temp := strconv.Itoa(jobComplPerc)
				w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
				w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(temp)))
				w.Write([]byte(temp))
			} else {
				temp := strconv.Itoa(jobComplPerc)
				w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
				w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(temp)))
				w.Write([]byte(temp))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func OutputHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodGet {
		id := ps.ByName(consts.RequestJobId)
		if  id != "" {
			session, err := store.Get(r, consts.SessionName)
			checkError(err)

			jobs := session.Values[consts.SessionGeneratorJob].(map[string]*GeneratorJob)
			if jobs != nil && jobs[id] != nil {
				job := jobs[id]
				if job.Completion == 100 {
					fmt.Println("okei")
					gcode, err := ioutil.ReadFile("/outputs/"+id+".gcode")
					checkError(err)
					w.Header().Set(consts.HttpContentDisposition, "attachment; filename=file.gcode")
					w.Header().Set(consts.HttpContentType, r.Header.Get("Content-Type"))
					w.Write(gcode)
				} else {
					w.WriteHeader(http.StatusConflict)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
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
		session.Save(r, w)

		if j := session.Values[consts.SessionGeneratorJob].(map[string]*GeneratorJob)[id]; j != nil {
			fmt.Println("La mappa della sessione in StartGeneratorJob esiste all'id: " + id)
			fmt.Println("Il job all'id " + id + "è: " + j.Id)
		}

		scaleFactor, err := strconv.Atoi(r.FormValue(consts.FormScaleFactor))
		checkError(err)
		modelThickness, err := strconv.Atoi(r.FormValue(consts.FormModelThickness))
		checkError(err)
		travelSpeed, err := strconv.Atoi(r.FormValue(consts.FormTravelSpeed))
		checkError(err)

		generationParams := GeneratorParams{
			ScaleFactor:    scaleFactor,
			ModelThickness: modelThickness,
			TravelSpeed:    travelSpeed,
		}

		go StartGeneratorJob(r, jobs[id], &generationParams, jobCompletion)

		fmt.Println("L'id generato in StartGeneratorJob come chiave del Job è: " + id)

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
