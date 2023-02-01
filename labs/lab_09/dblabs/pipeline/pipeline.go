package main

import (
	"sync"
	"time"
)

type Line func(interface{}) (interface{}, error)

type Pipeline struct {
	Lines []Line
}

// RunSequential TODO: add context with timeout
func (p Pipeline) RunSequential(requests []Request) error {
	queues := make([]chan Request, len(p.Lines)+1)
	for i := range queues {
		queues[i] = make(chan Request, len(requests))
	}

	for _, r := range requests {
		r.LineTimes = append(r.LineTimes, LineTime{StartQueue: time.Now().UnixNano()})
		queues[0] <- r
	}

	for i := range p.Lines {
		for range requests {
			request := <-queues[i]
			request.LineTimes[i].EndQueue = time.Now().UnixNano()

			newData, err := p.Lines[i](request.Data)
			if err != nil {
				return err
			}
			request.Data = newData
			request.LineTimes[i].EndProcessing = time.Now().UnixNano()

			request.LineTimes = append(request.LineTimes, LineTime{StartQueue: time.Now().UnixNano()})
			queues[i+1] <- request
		}
	}

	for i := range requests {
		requests[i] = <-queues[len(queues)-1]
	}

	return nil
}

func (p Pipeline) RunParallel(requests []Request) error {
	queues := make([]chan Request, len(p.Lines)+1)
	for i := range queues {
		queues[i] = make(chan Request, len(requests))
	}

	for _, r := range requests {
		r.LineTimes = append(r.LineTimes, LineTime{StartQueue: time.Now().UnixNano()})
		queues[0] <- r
	}

	wg := sync.WaitGroup{}

	for i := range p.Lines {
		wg.Add(1)
		go func(lineIndex int, line func(interface{}) (interface{}, error), qin chan Request, qout chan Request) {
			for range requests {
				request := <-qin
				request.LineTimes[lineIndex].EndQueue = time.Now().UnixNano()

				newData, err := line(request.Data)
				if err != nil {
					panic(err)
				}
				request.Data = newData
				request.LineTimes[lineIndex].EndProcessing = time.Now().UnixNano()

				request.LineTimes = append(request.LineTimes, LineTime{StartQueue: time.Now().UnixNano()})
				qout <- request
			}
			wg.Done()
		}(i, p.Lines[i], queues[i], queues[i+1])
	}

	wg.Wait()

	for i := range requests {
		requests[i] = <-queues[len(queues)-1]
	}

	return nil
}
