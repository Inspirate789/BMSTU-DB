package main

import (
	"encoding/json"
	"flag"
	"os"
)

type BenchResult struct {
	Requests           []Request `json:"requests"`
	LinesCount         int       `json:"lines_count"`
	MinQueueTimes      []int64   `json:"min_queue_times"`
	AvgQueueTimes      []float32 `json:"avg_queue_times"`
	MaxQueueTimes      []int64   `json:"max_queue_times"`
	MinProcessingTimes []int64   `json:"min_processing_times"`
	AvgProcessingTimes []float32 `json:"avg_processing_times"`
	MaxProcessingTimes []int64   `json:"max_processing_times"`
}

var requestCount int

func init() {
	flag.IntVar(&requestCount, "rcnt", 0, "Request Count")
}

func main() {
	flag.Parse()

	lines := []Line{connectDB, insertTerm, selectTerm, deleteTerm, disconnectDB}
	p := Pipeline{lines}

	funcs := []func(Pipeline, []Request) error{Pipeline.RunSequential, Pipeline.RunParallel}
	logfiles := []string{"log/sync.json", "log/parallel.json"}

	for i := range funcs {
		requests := make([]Request, requestCount)
		for j := 0; j < requestCount; j++ {
			requests[j].ID = j
		}

		err := funcs[i](p, requests)
		if err != nil {
			panic(err)
		}

		b := BenchResult{
			Requests:           requests,
			LinesCount:         len(lines),
			MinQueueTimes:      GetMinQueueTime(requests)[:len(lines)],
			AvgQueueTimes:      GetAvgQueueTime(requests)[:len(lines)],
			MaxQueueTimes:      GetMaxQueueTime(requests)[:len(lines)],
			MinProcessingTimes: GetMinProcessingTime(requests)[:len(lines)],
			AvgProcessingTimes: GetAvgProcessingTime(requests)[:len(lines)],
			MaxProcessingTimes: GetMaxProcessingTime(requests)[:len(lines)],
		}

		file, _ := os.OpenFile(logfiles[i], os.O_CREATE|os.O_WRONLY, os.ModePerm)

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		err = encoder.Encode(b)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
}
