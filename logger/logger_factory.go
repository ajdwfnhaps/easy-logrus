package logger

import (
	"context"
)

//LogFactory logger工厂
type LogFactory struct {
}

//NewLoggerFactory new工厂
func NewLoggerFactory() *LogFactory {
	return new(LogFactory)
}

// CreateLogger 创建logger
func (f *LogFactory) CreateLogger(ctx context.Context) *Entry {
	return CreateLoggerWithContext(ctx)
}
