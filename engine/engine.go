// Package engine
// @description implements the schedule engine of cralwer
package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
	"go.uber.org/zap"
)

type ScheduleEngine struct {
	requestCh chan *fetcher.Request
	workerCh  chan *fetcher.Request
	WorkCount int
	Fetcher   fetcher.Fetcher
	Logger    *zap.Logger
	out       chan fetcher.ParseResult
	Seeds     []*fetcher.Request
}

func (s *ScheduleEngine) Run() {
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

func (s *ScheduleEngine) Schedule() {
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

func (s *ScheduleEngine) CreateWork() {
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

func (s *ScheduleEngine) HandleResult() {
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
