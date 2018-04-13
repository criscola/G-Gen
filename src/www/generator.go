package main

import (
	"fmt"
	"ggen/utils/config"
	"ggen/utils/consts"
	"github.com/go-cmd/cmd"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	Dialect        string
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

	// Replace transparent pixels with white pixels
	//fillTransparentPixels := cmd.NewCmd("convert", "-flatten", imageFilepath, imageFilepath)
	/*
	fillTransparentPixels := cmd.NewCmd("convert", "-flatten", imageFilepath, imageFilepath)

	fmt.Println("imageFilepath:", imageFilepath)
	s = <-fillTransparentPixels.Start()
	fmt.Println("convert: ", s.Stdout)
	fmt.Println("convert: ", s.Stderr)*/

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
	fmt.Println("openscad: ")
	fmt.Print(s.Stdout)

	job.updateCompletion(80, c)

	// slice .stl to .gcode
	gcodeFilepath := pwdOutput + "/outputs/" + imageFilenameNoExt + ".gcode"
	fmt.Println("DIALECT IS: " + job.Params.Dialect + " CONSTS IS: " + consts.FormRepRap + " " + consts.FormUltimaker)
	if job.Params.Dialect == consts.FormRepRap {
		fmt.Println("Slicing with slic3r...")
		slic3r := cmd.NewCmd("slic3r", stlFilepath, "--output", gcodeFilepath)
		s = <-slic3r.Start()
		fmt.Print(s.Stdout)
		fmt.Print(s.Stderr)
	} else if job.Params.Dialect == consts.FormUltimaker {
		fmt.Println("Slicing with ultimaker...")
		cura := cmd.NewCmd("CuraEngine", "slice", "-v", "-j", config.ResourcesDirPath+"/fdmprinter.def.json",
			"-o", gcodeFilepath, "-l", stlFilepath,
			"-s", "expand_skins_expand_distance=0",
			"-s", "speed_infill="+strconv.Itoa(job.Params.TravelSpeed),
		)
		s = <-cura.Start()
		fmt.Println(cura.Args)
		fmt.Print(s.Stderr)
		fmt.Print(s.Stdout)
	}

	job.FinishTime = time.Now().Unix()
	job.updateCompletion(100, c)
}

func (job *GeneratorJob) updateCompletion(completion int, c chan *GeneratorJob) {
	job.Completion = completion
	c <- job
}
