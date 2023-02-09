package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
	"go.uber.org/zap"
)

type Schedule struct {
	requestCh chan *fetcher.Request
	workerCh  chan *fetcher.Request
	reqQueue  []*fetcher.Request
	Logger    *zap.Logger
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	s.requestCh = make(chan *fetcher.Request)
	s.workerCh = make(chan *fetcher.Request)
	return s
}

func (s *Schedule) Push(reqs ...*fetcher.Request) {
	for _, req := range reqs {
		s.requestCh <- req
	}
}

func (s *Schedule) Pull() *fetcher.Request {
	r := <-s.workerCh
	return r
}

func (s *Schedule) Schedule() {
	for {
		var req *fetcher.Request
		var ch chan *fetcher.Request

		if len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerCh
		}

		select {
		case r := <-s.requestCh:
			s.reqQueue = append(s.reqQueue, r)
		case ch <- req:
		}
	}
}
