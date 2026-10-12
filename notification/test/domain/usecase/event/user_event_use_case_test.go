package event_test

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
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
	eventusecase "github.com/BenyaChef/micro/notification/domain/usecase/event"
	eventusecasehandler "github.com/BenyaChef/micro/notification/domain/usecase/event/handler"
	eventusecasemock "github.com/BenyaChef/micro/notification/test/domain/usecase/event/mock"
)

type UserEventUseCaseShould struct {
	suite.Suite

	repository *eventusecasemock.NotificationRepositoryMock
	useCase    *eventusecase.UserEventUseCase
}

func TestUserEventUseCaseShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserEventUseCaseShould))
}

func (s *UserEventUseCaseShould) SetupTest() {
	logPublisher := loggerteststub.NewLogPublisherStub()
	s.repository = eventusecasemock.NewNotificationRepositoryMock()

	handler, err := eventusecasehandler.NewBuilder().
		NotificationRepository(s.repository).
		LogPublisher(logPublisher).
		BuildRegisteredHandler()
	require.NoError(s.T(), err)

	provider, err := eventusecasehandler.NewProvider(
		map[commonautheventmodel.UserEventAction]eventusecasehandler.Handler{
			commonautheventmodel.UserEventActionRegistered: handler,
		},
	)
	require.NoError(s.T(), err)

	useCase, err := eventusecase.NewBuilder().
		HandlerProvider(provider).
		LogPublisher(logPublisher).
		Build()
	require.NoError(s.T(), err)

	s.useCase = useCase
}

func (s *UserEventUseCaseShould) validEvent() *commonautheventmodel.UserEvent {
	return commonautheventmodel.NewUserRegistered(
		idprimitive.NewEntityID().String(),
		commonuserentity.NewUserID().String(),
		"a@b.c",
		time.Now(),
	)
}

func (s *UserEventUseCaseShould) TestBuild_WithoutParams_ReturnError() {
	useCase, err := eventusecase.NewBuilder().Build()

	assert.Nil(s.T(), useCase)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, eventusecase.ErrHandlerProviderIsRequired))
	assert.True(s.T(), errors.Is(err, eventusecase.ErrLogPublisherIsRequired))
}

func (s *UserEventUseCaseShould) TestNewProvider_Empty_ReturnError() {
	provider, err := eventusecasehandler.NewProvider(nil)

	assert.Nil(s.T(), provider)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, eventusecasehandler.ErrHandlersAreRequired))
}

func (s *UserEventUseCaseShould) TestHandle_NilEvent_ReturnError() {
	err := s.useCase.Handle(s.T().Context(), nil)

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, eventusecase.ErrEventIsRequired))
}

func (s *UserEventUseCaseShould) TestHandle_InvalidEvent_ReturnError() {
	event := s.validEvent()
	event.Data.UserID = ""

	err := s.useCase.Handle(s.T().Context(), event)

	require.Error(s.T(), err)
	assert.Empty(s.T(), s.repository.Inserted)
}

func (s *UserEventUseCaseShould) TestHandle_ActionUnknownToSchema_ReturnError() {
	event := s.validEvent()
	event.Action = "DELETED"

	err := s.useCase.Handle(s.T().Context(), event)

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, commonautheventmodel.ErrCodeActionUnsupported))
	assert.Empty(s.T(), s.repository.Inserted)
}

func (s *UserEventUseCaseShould) TestHandle_KnownActionWithoutHandler_SkipsWithoutError() {
	provider, err := eventusecasehandler.NewProvider(
		map[commonautheventmodel.UserEventAction]eventusecasehandler.Handler{
			commonautheventmodel.UserEventAction("SOMETHING_ELSE"): stubHandler{},
		},
	)
	require.NoError(s.T(), err)

	useCase, err := eventusecase.NewBuilder().
		HandlerProvider(provider).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	err = useCase.Handle(s.T().Context(), s.validEvent())

	require.NoError(s.T(), err, "действие без обработчика пропускаем, а не блокируем партицию")
	assert.Empty(s.T(), s.repository.Inserted)
}

type stubHandler struct{}

func (stubHandler) Handle(context.Context, *commonautheventmodel.UserEvent) error {
	return nil
}

func (s *UserEventUseCaseShould) TestHandle_ValidEvent_StoresNotification() {
	event := s.validEvent()

	err := s.useCase.Handle(s.T().Context(), event)

	require.NoError(s.T(), err)
	require.Len(s.T(), s.repository.Inserted, 1)

	stored := s.repository.Inserted[0]
	assert.Equal(s.T(), event.EventID, stored.EventID().String())
	assert.Equal(s.T(), event.Data.UserID, stored.UserID().String())
	assert.Equal(s.T(), notificationentity.KindWelcome, stored.Kind())
	assert.Contains(s.T(), stored.Payload().String(), event.Data.Email)
}

func (s *UserEventUseCaseShould) TestHandle_Duplicate_IsNotAnError() {
	s.repository.InsertFunc = func(context.Context, *notificationentity.Notification) (bool, error) {
		return false, nil
	}

	err := s.useCase.Handle(s.T().Context(), s.validEvent())

	require.NoError(s.T(), err)
}

func (s *UserEventUseCaseShould) TestHandle_RepositoryFails_ReturnErrorSoOffsetIsNotCommitted() {
	cause := errors.New("database is down")
	s.repository.InsertFunc = func(context.Context, *notificationentity.Notification) (bool, error) {
		return false, cause
	}

	err := s.useCase.Handle(s.T().Context(), s.validEvent())

	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, cause)
	assert.True(s.T(), apperrors.EqualByCode(err, eventusecasehandler.ErrCodeStoreFailed))
}
