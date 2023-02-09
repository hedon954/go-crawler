package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
	"go.uber.org/zap"
)

type Option func(opts *options)

type options struct {
	WorkCount int
	Fetcher   fetcher.Fetcher
	Logger    *zap.Logger
	Seeds     []*fetcher.Task
}

var defaultOptions = options{
	Logger: zap.NewNop(),
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
