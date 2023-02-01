package main

type LineTime struct {
	StartQueue    int64 `json:"start_queue"`
	EndQueue      int64 `json:"end_queue"`
	EndProcessing int64 `json:"end_processing"`
}

type Request struct {
	ID        int         `json:"id"`
	Data      interface{} `json:"-"`
	LineTimes []LineTime  `json:"line_times"`
}

func (r Request) GetQueueTime(qind int) int64 {
	return r.LineTimes[qind].EndQueue - r.LineTimes[qind].StartQueue
}

func GetAvgQueueTime(requests []Request) []float32 {
	if len(requests) == 0 {
		return nil
	}

	res := make([]float32, len(requests[0].LineTimes))
	for i := range res {
		for _, r := range requests {
			res[i] += float32(r.GetQueueTime(i))
		}
		res[i] /= float32(len(requests))
	}

	return res
}

func GetMinQueueTime(requests []Request) []int64 {
	if len(requests) == 0 {
		return nil
	}

	res := make([]int64, len(requests[0].LineTimes))
	for i := range res {
		res[i] = requests[0].GetQueueTime(i)
		for _, r := range requests {
			if t := r.GetQueueTime(i); t < res[i] {
				res[i] = t
			}
		}
	}

	return res
}

func GetMaxQueueTime(requests []Request) []int64 {
	if len(requests) == 0 {
		return nil
	}

	res := make([]int64, len(requests[0].LineTimes))
	for i := range res {
		res[i] = requests[0].GetQueueTime(i)
		for _, r := range requests {
			if t := r.GetQueueTime(i); t > res[i] {
				res[i] = t
			}
		}
	}

	return res
}

func (r Request) GetProcessingTime(pind int) int64 {
	return r.LineTimes[pind].EndProcessing - r.LineTimes[pind].EndQueue
}

func GetAvgProcessingTime(requests []Request) []float32 {
	if len(requests) == 0 {
		return nil
	}

	sums := make([]float32, len(requests[0].LineTimes))
	for i := range sums {
		for _, r := range requests {
			sums[i] += float32(r.GetProcessingTime(i))
		}
		sums[i] /= float32(len(requests))
	}

	return sums
}

func GetMinProcessingTime(requests []Request) []int64 {
	if len(requests) == 0 {
		return nil
	}

	res := make([]int64, len(requests[0].LineTimes))
	for i := range res {
		res[i] = requests[0].GetProcessingTime(i)
		for _, r := range requests {
			if t := r.GetQueueTime(i); t < res[i] {
				res[i] = t
			}
		}
	}

	return res
}

func GetMaxProcessingTime(requests []Request) []int64 {
	if len(requests) == 0 {
		return nil
	}

	res := make([]int64, len(requests[0].LineTimes))
	for i := range res {
		res[i] = requests[0].GetProcessingTime(i)
		for _, r := range requests {
			if t := r.GetQueueTime(i); t > res[i] {
				res[i] = t
			}
		}
	}

	return res
}
