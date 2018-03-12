package main

import (

	"ggen/utils/consts"
	"net/http"
	"os/exec"
	"strings"
	"path/filepath"
	"time"
	"os"
	"strconv"
)

type GeneratorJob struct {
	Id 			string
	StartTime  	int64
	FinishTime 	int64
	Completion 	int
}

type GeneratorParams struct {
	ScaleFactor 	int
	ModelThickness 	int
	TravelSpeed 	int
}

/*
 * Do not attempt to write to the 'r' parameter. This would NOT be a thread-safe operation.
 */
func StartGeneratorJob(r *http.Request, job *GeneratorJob, params *GeneratorParams, c chan *GeneratorJob) {
/*	time.Sleep(5 * time.Second)
	job.Completion = 60
	c <- job
	time.Sleep(3 * time.Second)
	job.Completion = 78
	time.Sleep(2 * time.Second)
	job.Completion = 100
	job.FinishTime = time.Now().Unix()
	c <- job*/
	time.Sleep(time.Second)
	// Create copy of request
	newRequestAlloc := *r
	session, err := store.Get(&newRequestAlloc, consts.SessionName)
	checkError(err)

	imageFilename := session.Values[consts.SessionImageFilename].(string)
	imageFilenameNoExt := strings.TrimSuffix(imageFilename, filepath.Ext(imageFilename))
	// Get current path
	pwd := exec.Command("pwd")
	pwdOutput, err := pwd.Output()
	checkError(err)

	imageFilepath := string(pwdOutput)+"/uploads/"+imageFilename
	tmpPath := string(pwdOutput[:len(pwdOutput)-1])+"/uploads/tmp/" // Remove newline byte
	scadFilepath := tmpPath+imageFilenameNoExt+".scad"

	// image to .scad
	test := exec.Command("trace2scad", "-f", "0", "-e", "10", "-o", scadFilepath, imageFilepath)
	test.Wait()

	scadFile, err := os.OpenFile(scadFilepath, os.O_APPEND|os.O_WRONLY, 0600)
	checkError(err)
	defer scadFile.Close()

	// add scaling and render .scad to .stl
	scalingParams := strconv.Itoa(params.ScaleFactor)
	modelThickness := strconv.Itoa(params.ModelThickness)
	scadFile.WriteString("\nscale(["+scalingParams+", "+scalingParams+", "+modelThickness+"]) {\n\t"+imageFilenameNoExt+"();\n}")

	stlFilepath := tmpPath+imageFilenameNoExt+".stl"

	exec.Command("openscad", "-o", stlFilepath, scadFilepath)

	gcodeFilepath := tmpPath+imageFilenameNoExt+".gcode"
	// slice .stl to .gcode
	exec.Command("slic3r", stlFilepath, "--output", gcodeFilepath)
}