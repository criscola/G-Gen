package main

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"ggen/utils/config"
	"ggen/utils/consts"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	router.ServeFiles("/viewer/*filepath", http.Dir("viewer"))

	/** ROUTES **/
	router.GET("/", IndexHandler)
	router.GET("/generator", GeneratorHandler)
	router.GET("/uploads/:"+consts.RequestFilename, ImageGetHandler)
	router.GET("/generator/queue/:"+consts.RequestJobId, QueueHandler)
	router.GET("/generator/outputs/:"+consts.RequestJobId, OutputHandler)
	router.POST("/generator/imageUpload", ImagePostHandler)
	router.POST("/generator/generate", StartGeneratorJobHandler)
	router.DELETE("/generator/imageRemove/:"+consts.RequestJobId, ImageRemoveHandler)

	http.ListenAndServe(":"+config.ServerPort, context.ClearHandler(router))

}

// IndexHandler handles the index page route
func IndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templates = template.Must(template.ParseFiles("templates/home/index.tmpl", "templates/base.tmpl"))
	err := templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GeneratorHandler handles the generator page route
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

// ImageGetHandler handles get requests for images inside the upload directory
func ImageGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodGet {
		fmt.Println("ID FROM CLIENT IS: " + ps.ByName(consts.RequestFilename))
		job := getGeneratorJobById(r, ps.ByName(consts.RequestFilename))

		fileName := job.FileNames + "." + consts.DefaultImageExtension
		fmt.Println("image fileName: " + fileName)
		// If the extension is jpg convert to jpeg for the content-type response
		/*
		extension := strings.TrimPrefix(filepath.Ext(fileName), ".")
		if strings.EqualFold(extension, "jpg") {
			extension = "jpeg"
		}
		w.Header().Set(consts.HttpContentType, "image/"+extension)
		*/

		w.Header().Set(consts.HttpContentType, "image/"+consts.DefaultImageExtension)
		// Read image from disk
		data, err := ioutil.ReadFile(filepath.Join("./uploads/", fileName))
		checkError(err)

		w.Write(data)
	}
}

// QueueHandler returns data about a job in progress by their jobId provided through ps parameter
func QueueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodGet {
		session := getSession(r)

		// Get the request id
		id := ps.ByName(consts.RequestJobId)
		// Get the job map from session
		jobs := getGeneratorJobs(session)

		// check if id provided by the HTTP request isn't missing and that the job[id] exists
		if id != "" && jobs[id] != nil {
			// Get the last selected job supplied by the jobCompletion channel
			var selectedJob *GeneratorJob
			percCount := len(jobCompletion)

			if percCount > 0 && percCount < 100 {
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

// OutputHandler handles the download requests by their jobId provided through ps
func OutputHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodGet {
		fmt.Println("outputhandler invoked")
		if id := ps.ByName(consts.RequestJobId); id != "" {
			fmt.Println("id outputhandler is: " + id)
			if job := getGeneratorJobById(r, id); job != nil {
				if job.Completion == 100 {
					gcode, err := ioutil.ReadFile("./outputs/" + job.FileNames + ".gcode")
					fmt.Println("Job filenames is : " + job.FileNames)
					checkError(err)

					w.Header().Set(consts.HttpContentDisposition, "attachment; filename=\"file.gcode\"")
					w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(gcode)))
					w.Header().Set(consts.HttpContentType, consts.HttpMimeApplicationOctetStream)
					w.Write(gcode)
				} else {
					w.WriteHeader(http.StatusConflict)
				}
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

// ImagePostHandler handles image uploads to the server
func ImagePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile(consts.FormImage)
		checkError(err)

		defer file.Close()

		// Generate random filename which will be used widely for all the files used in the current job
		fileNames := TempFileName()
		fmt.Println("FILENAMES GENERATED IS: " + fileNames)
		// Generate random job id
		session := getSession(r)

		// Generate id for job and add to session struct
		var jobs map[string]*GeneratorJob
		var jobId string

		// If the job map is already existing, get job, or else make a new job map
		if session.Values[consts.SessionGeneratorJob] != nil {
			jobs = session.Values[consts.SessionGeneratorJob].(map[string]*GeneratorJob)
			fmt.Println("GOT JOB MAP")
		} else {
			jobs = make(map[string]*GeneratorJob)
			fmt.Println("CREATED JOB MAP")
		}

		// Create new job id
		jobId = GetRandomString()
		fmt.Println("JOBID GENERATED IS: " + jobId)

		// Create a new GeneratorJob
		jobs[jobId] = &GeneratorJob{fileNames, 0, 0, 0, GeneratorParams{}}
		session.Values[consts.SessionGeneratorJob] = jobs
		sessions.Save(r, w)

		// If image is not .png format, convert to .png
		extension := strings.ToLower(filepath.Ext(handler.Filename))
		imagePath := filepath.Join("./uploads/", fileNames+extension)
		/*if extension != ".png" {
			imageToPngCommand := cmd.NewCmd("convert " + imagePath + " ./uploads/" + fileNames + ".png")
			imageToPngCommand.Start()
			fmt.Println("Converting image to png...")
		}*/

		// Open/Create temp file for image using name stored in the session store
		f, err := os.OpenFile(imagePath, os.O_WRONLY|os.O_CREATE, 0666)
		checkError(err)

		defer f.Close()
		_, err = io.Copy(f, file)
		checkError(err)

		// Return job id to client
		w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
		w.Header().Set(consts.HttpContentLength, strconv.Itoa(len(jobId)))
		w.Write([]byte(jobId))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// StartGeneratorJobHandler starts a job using the data provided in the POST request body
func StartGeneratorJobHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("startgenjob")
	if r.Method == http.MethodPost {
		// Start goroutine for the generator job
		// Add goroutine to map holding all jobs with keys as id
		// Send back 202 (accepted) status code
		// Get session
		session := getSession(r)

		id := r.FormValue(consts.RequestJobId)
		fmt.Println("ID GENJOB: " + id)
		if id != "" {
			scaleFactor, err := strconv.Atoi(r.FormValue(consts.FormScaleFactor))
			checkError(err)
			modelThickness, err := strconv.Atoi(r.FormValue(consts.FormModelThickness))
			checkError(err)
			travelSpeed, err := strconv.Atoi(r.FormValue(consts.FormTravelSpeed))
			checkError(err)

			genJob := getGeneratorJobs(session)[id]

			genJob.Params = GeneratorParams{
				ScaleFactor:    scaleFactor,
				ModelThickness: modelThickness,
				TravelSpeed:    travelSpeed,
			}

			go StartGeneratorJob(genJob, jobCompletion)

			// Write in response body the id
			w.Header().Set(consts.HttpContentType, consts.HttpMimeTextPlain)
			w.Header().Set(consts.HttpContentLength, strconv.Itoa(len("OK")))
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ImageRemoveHandler handles DELETE requests for images
func ImageRemoveHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodDelete {
		if id := ps.ByName(consts.RequestJobId); id != "" {
			imageName := getGeneratorJobById(r, id).FileNames + "." + consts.DefaultImageExtension
			fmt.Println("IMAGENAME IN REMOVEHANDLER IS : " + imageName)
			err := os.Remove(filepath.Join("./uploads/", imageName))
			checkError(err)
		} else {
			w.WriteHeader(http.StatusBadRequest)
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

// GetRandomString generates a 16-bytes alphanumeric random string
func GetRandomString() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}

// checkError checks for an error and panics if something went wrong
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// getSession gets the user session through its HTTP request
func getSession(r *http.Request) *sessions.Session {
	session, err := store.Get(r, consts.SessionName)
	if session.IsNew {
		fmt.Println("created new session " + session.Name())
	}
	checkError(err)
	return session
}

// getGeneratorJobs gets the map containing all the job belonging to a session specified by the session parameter
func getGeneratorJobs(session *sessions.Session) map[string]*GeneratorJob {
	return session.Values[consts.SessionGeneratorJob].(map[string]*GeneratorJob)
}

// getGeneratorJobById gets a job associated with an id of a specific session
func getGeneratorJobById(r *http.Request, id string) *GeneratorJob {
	return getGeneratorJobs(getSession(r))[id]
}

func init() {
	fmt.Println("Starting webserver at " + config.ServerURL + ":" + config.ServerPort + "...")
	gob.Register(map[string]*GeneratorJob{})
}
