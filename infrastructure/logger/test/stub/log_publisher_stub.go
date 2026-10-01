package loggerteststub

import (
	"context"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

var _ loggerinterface.LogPublisher = (*LogPublisherStub)(nil)

type LogPublisherStub struct{}

func NewLogPublisherStub() *LogPublisherStub {
	return &LogPublisherStub{}
}

func (s *LogPublisherStub) LogError(_ context.Context, _ ...error)         {}
func (s *LogPublisherStub) LogWarn(_ context.Context, _ string, _ ...any)  {}
func (s *LogPublisherStub) LogInfo(_ context.Context, _ string, _ ...any)  {}
func (s *LogPublisherStub) LogDebug(_ context.Context, _ string, _ ...any) {}
