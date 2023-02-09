package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
	"go.uber.org/zap"
)

type Scheduler struct {
	requestCh chan *fetcher.Request
	workerCh  chan *fetcher.Request
	out       chan fetcher.ParseResult
	options
}

func NewScheduler(opts ...Option) *Scheduler {
	dopts := defaultOptions
	for _, opt := range opts {
		opt(&dopts)
	}
	s := &Scheduler{}
	s.options = dopts
	return s
}

func (s *Scheduler) Run() {
	requestCh := make(chan *fetcher.Request)
	workerCh := make(chan *fetcher.Request)
	out := make(chan fetcher.ParseResult)
	s.requestCh = requestCh
	s.workerCh = workerCh
	s.out = out
	go s.Schedule()
	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}
	s.HandleResult()
}

func (s *Scheduler) Schedule() {
	for {
		var req *fetcher.Request
		var ch chan *fetcher.Request
		if len(s.Seeds) > 0 {
			req = s.Seeds[0]
			s.Seeds = s.Seeds[1:]
			ch = s.workerCh
		}
		select {
		case r := <-s.requestCh:
			s.Seeds = append(s.Seeds, r)
		case ch <- req:
		}
	}
}

func (s *Scheduler) CreateWork() {
	for {
		r := <-s.workerCh
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.Logger.Error("can not fetch ", zap.Error(err))
			continue
		}
		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

func (s *Scheduler) HandleResult() {
	for {
		select {
		case result := <-s.out:
			for _, req := range result.Requests {
				s.requestCh <- req
			}
			for _, item := range result.Items {
				s.Logger.Sugar().Info("get result", item)
			}
		}
	}
}
