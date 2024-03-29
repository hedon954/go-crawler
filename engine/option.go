package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
	"go.uber.org/zap"
)

type Option func(opts *options)

type options struct {
	WorkCount     int
	ChannelBuffer int
	Fetcher       fetcher.Fetcher
	Logger        *zap.Logger
	Seeds         []*fetcher.Task
	scheduler     Scheduler
}

var defaultOptions = options{
	Logger:        zap.NewNop(),
	ChannelBuffer: 1024,
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.Logger = logger
	}
}

func WithFetcher(f fetcher.Fetcher) Option {
	return func(opts *options) {
		opts.Fetcher = f
	}
}

func WithWorkCount(workCount int) Option {
	return func(opts *options) {
		opts.WorkCount = workCount
	}
}

func WithSeeds(seeds []*fetcher.Task) Option {
	return func(opts *options) {
		opts.Seeds = seeds
	}
}

func WithScheduler(scheduler Scheduler) Option {
	return func(opts *options) {
		opts.scheduler = scheduler
	}
}

func WithChannelBuffer(bufferSize int) Option {
	return func(opts *options) {
		opts.ChannelBuffer = bufferSize
	}
}
