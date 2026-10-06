package grpcserver_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/BenyaChef/micro/infrastructure/grpcserver"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
)

type BuilderShould struct {
	suite.Suite
}

func TestBuilderShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuilderShould))
}

func (s *BuilderShould) TestBuild_WithoutParams_ReturnError() {
	server, err := grpcserver.NewBuilder().Build()

	assert.Nil(s.T(), server)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, grpcserver.ErrLogPublisherIsRequired))
	assert.True(s.T(), errors.Is(err, grpcserver.ErrPortIsRequired))
}

func (s *BuilderShould) TestBuild_WithoutPort_ReturnOnlyPortError() {
	server, err := grpcserver.NewBuilder().
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()

	assert.Nil(s.T(), server)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, grpcserver.ErrPortIsRequired))
	assert.False(s.T(), errors.Is(err, grpcserver.ErrLogPublisherIsRequired))
}

func (s *BuilderShould) TestBuild_AllParams_ReturnServerWithAddr() {
	server, err := grpcserver.NewBuilder().
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Port("9091").
		Build()

	require.NoError(s.T(), err)
	require.NotNil(s.T(), server)
	assert.Equal(s.T(), ":9091", server.Addr())
	assert.NotNil(s.T(), server.Server())
}
