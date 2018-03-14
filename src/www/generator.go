package main

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ggen/utils/consts"
)

type GeneratorJob struct {
	Id         string
	StartTime  int64
	FinishTime int64
	Completion int
}

type GeneratorParams struct {
	ScaleFactor    int
	ModelThickness int
	TravelSpeed    int
}

/*
 * Do not attempt to write to the 'r' parameter. This would NOT be a thread-safe operation.
 */
func StartGeneratorJob(r *http.Request, job *GeneratorJob, params *GeneratorParams, c chan *GeneratorJob) {
	time.Sleep(time.Second)
	// Create copy of request
	newRequestAlloc := *r
	session, err := store.Get(&newRequestAlloc, consts.SessionName)
	checkError(err)

	job.updateCompletion(10, c)

	imageFilename := session.Values[consts.SessionImageFilename].(string)
	imageFilenameNoExt := strings.TrimSuffix(imageFilename, filepath.Ext(imageFilename))

	// Get current path
	pwd := cmd.NewCmd("pwd")
	s := <-pwd.Start()
	pwdOutput := s.Stdout[0]

	imageFilepath := pwdOutput + "/uploads/" + imageFilename
	tmpPath := pwdOutput + "/uploads/tmp/" // Remove newline byte
	scadFilepath := tmpPath + imageFilenameNoExt + ".scad"

	job.updateCompletion(40, c)

	// image to .scad
	trace2scad := cmd.NewCmd("trace2scad", "-f", "0", "-e", "10", "-o", scadFilepath, imageFilepath)
	s = <-trace2scad.Start()
	fmt.Println(s.Stdout)

	scadFile, err := os.OpenFile(scadFilepath, os.O_APPEND|os.O_WRONLY, 0600)
	checkError(err)
	defer scadFile.Close()

	job.updateCompletion(60, c)

	// add scaling and render .scad to .stl
	scalingParams := strconv.Itoa(params.ScaleFactor)
	modelThickness := strconv.Itoa(params.ModelThickness)
	scadFile.WriteString("\nscale([" + scalingParams + ", " + scalingParams + ", " + modelThickness + "]) {\n\t" + imageFilenameNoExt + "();\n}")

	stlFilepath := tmpPath + imageFilenameNoExt + ".stl"

	openscad := cmd.NewCmd("openscad", "-o", stlFilepath, scadFilepath)
	s = <-openscad.Start()
	fmt.Println(s.Stdout)

	job.updateCompletion(80, c)

	gcodeFilepath := pwdOutput + "/outputs/" + imageFilenameNoExt + ".gcode"
	// slice .stl to .gcode
	slic3r := cmd.NewCmd("slic3r", stlFilepath, "--output", gcodeFilepath)
	s = <-slic3r.Start()
	fmt.Print(s.Stdout)

	job.updateCompletion(100, c)

}

func (job *GeneratorJob) updateCompletion(completion int, c chan *GeneratorJob) {
	job.Completion = completion
	c <- job
}
