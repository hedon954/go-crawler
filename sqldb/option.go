package sqldb

import "go.uber.org/zap"

type options struct {
	logger *zap.Logger
	sqlUrl string
}

var defaultOption = options{
	logger: zap.NewNop(),
}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithSqlUrl(sqlUrl string) Option {
	return func(opts *options) {
		opts.sqlUrl = sqlUrl
	}
}
