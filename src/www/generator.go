package main

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"ggen/utils/consts"
)

type GCodeDialect int

const (
	RepRap GCodeDialect = iota
	Ultimaker
)

type GeneratorJob struct {
	FileNames  string
	StartTime  int64
	FinishTime int64
	Completion int
	Params     GeneratorParams
}

type GeneratorParams struct {
	ScaleFactor    int
	ModelThickness int
	TravelSpeed    int
	Dialect			GCodeDialect
}


// StartGeneratorJob is used to start a generator job associated with a specific session
func StartGeneratorJob(job *GeneratorJob, c chan *GeneratorJob) {
	time.Sleep(time.Second)

	job.StartTime = time.Now().Unix()
	job.updateCompletion(10, c)

	imageFilename := job.FileNames
	imageFilenameNoExt := strings.TrimSuffix(imageFilename, filepath.Ext(imageFilename))
	fmt.Println("imageFilename for gcode is: ", imageFilename)
	// Get current path
	pwd := cmd.NewCmd("pwd")
	s := <-pwd.Start()
	pwdOutput := s.Stdout[0]

	imageFilepath := pwdOutput + "/uploads/" + imageFilename + "." + consts.DefaultImageExtension
	tmpPath := pwdOutput + "/uploads/tmp/" // Remove newline byte
	scadFilepath := tmpPath + imageFilenameNoExt + ".scad"

	job.updateCompletion(40, c)

	// image to .scad
	trace2scad := cmd.NewCmd("trace2scad", "-f", "0", "-e", "10", "-o", scadFilepath, imageFilepath)
	fmt.Println(trace2scad.Args)
	s = <-trace2scad.Start()
	fmt.Println(s.Stdout)
	fmt.Println(s.Stderr)

	scadFile, err := os.OpenFile(scadFilepath, os.O_APPEND|os.O_WRONLY, 0600)
	checkError(err)
	defer scadFile.Close()

	job.updateCompletion(60, c)

	// add scaling and render .scad to .stl
	scalingParams := strconv.Itoa(job.Params.ScaleFactor)
	modelThickness := strconv.Itoa(job.Params.ModelThickness)
	scadFile.WriteString("\nscale([" + scalingParams + ", " + scalingParams + ", " + modelThickness + "]) {\n\t" + imageFilenameNoExt + "();\n}")

	stlFilepath := tmpPath + imageFilenameNoExt + ".stl"

	openscad := cmd.NewCmd("openscad", "-o", stlFilepath, scadFilepath)
	s = <-openscad.Start()
	fmt.Println(s.Stdout)

	job.updateCompletion(80, c)

	// slice .stl to .gcode
	gcodeFilepath := pwdOutput + "/outputs/" + imageFilenameNoExt + ".gcode"
	if job.Params.Dialect == RepRap {
		slic3r := cmd.NewCmd("slic3r", stlFilepath, "--output", gcodeFilepath)
		s = <-slic3r.Start()
		fmt.Print(s.Stdout)
	} else if job.Params.Dialect == Ultimaker {
		cura := cmd.NewCmd("cura", "")
	}

	job.FinishTime = time.Now().Unix()
	job.updateCompletion(100, c)
}

func (job *GeneratorJob) updateCompletion(completion int, c chan *GeneratorJob) {
	job.Completion = completion
	c <- job
}