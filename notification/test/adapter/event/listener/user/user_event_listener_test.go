package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	idprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/id"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
	usereventlistener "github.com/BenyaChef/micro/notification/adapter/event/listener/user"
	usereventlistenermock "github.com/BenyaChef/micro/notification/test/adapter/event/listener/user/mock"
)

type UserEventListenerShould struct {
	suite.Suite

	useCase  *usereventlistenermock.UserEventUseCaseMock
	listener *usereventlistener.UserEventListener
}

func TestUserEventListenerShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserEventListenerShould))
}

func (s *UserEventListenerShould) SetupTest() {
	s.useCase = usereventlistenermock.NewUserEventUseCaseMock()

	listener, err := usereventlistener.NewBuilder().
		UserEventUseCase(s.useCase).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.listener = listener
}

func (s *UserEventListenerShould) message(value []byte) *queueservermodel.ConsumerMessage {
	return &queueservermodel.ConsumerMessage{
		Topic:     commonautheventmodel.TopicUser,
		Key:       "key",
		Value:     value,
		Partition: 0,
		Offset:    42,
	}
}

func (s *UserEventListenerShould) validEventBytes() []byte {
	s.T().Helper()

	event := commonautheventmodel.NewUserRegistered(
		idprimitive.NewEntityID().String(),
		commonuserentity.NewUserID().String(),
		"a@b.c",
		time.Now(),
	)

	raw, err := event.ToBytes()
	require.NoError(s.T(), err)

	return raw
}

func (s *UserEventListenerShould) TestBuild_WithoutParams_ReturnError() {
	listener, err := usereventlistener.NewBuilder().Build()

	assert.Nil(s.T(), listener)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, usereventlistener.ErrUserEventUseCaseIsRequired))
	assert.True(s.T(), errors.Is(err, usereventlistener.ErrLogPublisherIsRequired))
}

func (s *UserEventListenerShould) TestHandle_ValidMessage_PassesEventToUseCase() {
	err := s.listener.Handle(s.T().Context(), s.message(s.validEventBytes()))

	require.NoError(s.T(), err)
	require.Len(s.T(), s.useCase.Handled, 1)
	assert.Equal(s.T(), commonautheventmodel.UserEventActionRegistered, s.useCase.Handled[0].Action)
}

func (s *UserEventListenerShould) TestHandle_GarbageMessage_SkippedWithoutError() {
	err := s.listener.Handle(s.T().Context(), s.message([]byte("{не json")))

	require.NoError(s.T(), err, "нечитаемое сообщение коммитим, иначе оно заблокирует партицию навсегда")
	assert.Empty(s.T(), s.useCase.Handled)
}

func (s *UserEventListenerShould) TestHandle_ValidJSONButBrokenSchema_SkippedWithoutError() {
	err := s.listener.Handle(s.T().Context(), s.message([]byte(`{"event_id":"","action":"","data":{}}`)))

	require.NoError(s.T(), err)
	assert.Empty(s.T(), s.useCase.Handled)
}

func (s *UserEventListenerShould) TestHandle_NilMessage_SkippedWithoutError() {
	err := s.listener.Handle(s.T().Context(), nil)

	require.NoError(s.T(), err)
	assert.Empty(s.T(), s.useCase.Handled)
}

func (s *UserEventListenerShould) TestHandle_UseCaseFails_ErrorIsReturned() {
	cause := errors.New("database is down")
	s.useCase.HandleFunc = func(context.Context, *commonautheventmodel.UserEvent) error { return cause }

	err := s.listener.Handle(s.T().Context(), s.message(s.validEventBytes()))

	require.Error(s.T(), err, "сбой обработки не коммитим, сообщение придёт снова")
	assert.ErrorIs(s.T(), err, cause)
}
