package main

import "time"

type GeneratorJob struct {
	StartTime  int64
	FinishTime int64
	Completion int
}

type GeneratorParams struct {
	ScaleFactor float64
	TravelSpeed float64
}

func StartGeneratorJob(job *GeneratorJob, params *GeneratorParams) {
	job.Completion = 100
	job.FinishTime = time.Now().Unix()
}
