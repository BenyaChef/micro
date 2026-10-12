package notification_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	"github.com/BenyaChef/micro/infrastructure/postgres"
	notificationrepository "github.com/BenyaChef/micro/notification/adapter/repository/notification"
	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
	notificationentitystub "github.com/BenyaChef/micro/notification/test/domain/entity/notification/stub"
)

const envTestDSN = "POSTGRES_TEST_DSN"

type NotificationRepositoryShould struct {
	suite.Suite

	client     *postgres.Client
	repository *notificationrepository.NotificationRepository
}

func TestNotificationRepositoryShould(t *testing.T) {
	suite.Run(t, new(NotificationRepositoryShould))
}

func (s *NotificationRepositoryShould) SetupSuite() {
	dsn := os.Getenv(envTestDSN)
	if dsn == "" {
		s.T().Skipf("%s is not set, skipping test against live database", envTestDSN)
	}

	client, err := postgres.NewBuilder().DSN(dsn).Build(s.T().Context())
	require.NoError(s.T(), err)

	s.client = client

	repository, err := notificationrepository.NewBuilder().
		Executor(client).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.repository = repository
}

func (s *NotificationRepositoryShould) TearDownSuite() {
	if s.client != nil {
		s.client.Close()
	}
}

func (s *NotificationRepositoryShould) cleanup(eventID notificationentity.EventID) {
	s.T().Helper()

	s.T().Cleanup(func() {
		_, err := s.client.Exec(context.Background(), "DELETE FROM notifications WHERE event_id = $1", eventID.String())
		require.NoError(s.T(), err)
	})
}

func (s *NotificationRepositoryShould) TestBuild_WithoutParams_ReturnError() {
	repository, err := notificationrepository.NewBuilder().Build()

	assert.Nil(s.T(), repository)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, notificationrepository.ErrExecutorIsRequired))
	assert.True(s.T(), errors.Is(err, notificationrepository.ErrLogPublisherIsRequired))
}

func (s *NotificationRepositoryShould) TestInsert_NilNotification_ReturnError() {
	inserted, err := s.repository.Insert(s.T().Context(), nil)

	require.Error(s.T(), err)
	assert.False(s.T(), inserted)
	assert.True(s.T(), errors.Is(err, notificationrepository.ErrNotificationIsRequired))
}

func (s *NotificationRepositoryShould) TestInsert_FirstDelivery_ReturnsInserted() {
	notification := notificationentitystub.GetValidNotification()
	s.cleanup(notification.EventID())

	inserted, err := s.repository.Insert(s.T().Context(), notification)

	require.NoError(s.T(), err)
	assert.True(s.T(), inserted)
}

func (s *NotificationRepositoryShould) TestInsert_RedeliveredEvent_DoesNotCreateDuplicate() {
	notification := notificationentitystub.GetValidNotification()
	s.cleanup(notification.EventID())

	first, err := s.repository.Insert(s.T().Context(), notification)
	require.NoError(s.T(), err)
	require.True(s.T(), first)

	second, err := s.repository.Insert(s.T().Context(), notification)

	require.NoError(s.T(), err, "повторная доставка не ошибка")
	assert.False(s.T(), second, "вторая вставка не должна создавать запись")

	var count int
	row := s.client.QueryRow(s.T().Context(),
		"SELECT count(*) FROM notifications WHERE event_id = $1", notification.EventID().String())
	require.NoError(s.T(), row.Scan(&count))
	assert.Equal(s.T(), 1, count)
}

func (s *NotificationRepositoryShould) TestInsert_SameEventDifferentID_StillDeduplicated() {
	first := notificationentitystub.GetValidNotification()
	s.cleanup(first.EventID())

	inserted, err := s.repository.Insert(s.T().Context(), first)
	require.NoError(s.T(), err)
	require.True(s.T(), inserted)

	retry, err := notificationentity.NewBuilder().
		ID(notificationentity.NewNotificationID()).
		EventID(first.EventID()).
		UserID(commonuserentity.NewUserID()).
		Kind(notificationentity.KindWelcome).
		Payload(first.Payload()).
		Build()
	require.NoError(s.T(), err)

	inserted, err = s.repository.Insert(s.T().Context(), retry)

	require.NoError(s.T(), err)
	assert.False(s.T(), inserted, "дедупликация идёт по event_id, а не по id записи")
}

func (s *NotificationRepositoryShould) TestGetByEventID_Inserted_ReturnEqualEntity() {
	notification := notificationentitystub.GetValidNotification()
	s.cleanup(notification.EventID())

	_, err := s.repository.Insert(s.T().Context(), notification)
	require.NoError(s.T(), err)

	found, err := s.repository.GetByEventID(s.T().Context(), notification.EventID())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), notification.ID(), found.ID())
	assert.Equal(s.T(), notification.EventID(), found.EventID())
	assert.Equal(s.T(), notification.UserID(), found.UserID())
	assert.Equal(s.T(), notification.Kind(), found.Kind())
	assert.JSONEq(s.T(), notification.Payload().String(), found.Payload().String())
}

func (s *NotificationRepositoryShould) TestGetByEventID_Unknown_ReturnNilNil() {
	found, err := s.repository.GetByEventID(s.T().Context(), notificationentitystub.GetEventID())

	require.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}
