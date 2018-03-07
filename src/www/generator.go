package main

import (
	"time"
)

type GeneratorJob struct {
	Id string
	StartTime  int64
	FinishTime int64
	Completion int
}

type GeneratorParams struct {
	ScaleFactor float64
	TravelSpeed float64
}

func StartGeneratorJob(job *GeneratorJob, params *GeneratorParams, c chan *GeneratorJob) {
	time.Sleep(5 * time.Second)
	job.Completion = 60
	c <- job
	time.Sleep(3 * time.Second)
	job.Completion = 78
	time.Sleep(2 * time.Second)
	job.Completion = 100
	job.FinishTime = time.Now().Unix()
	c <- job
}
