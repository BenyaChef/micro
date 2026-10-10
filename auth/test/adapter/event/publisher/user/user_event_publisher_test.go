package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	usereventpublisher "github.com/BenyaChef/micro/auth/adapter/event/publisher/user"
	usereventpublishermock "github.com/BenyaChef/micro/auth/test/adapter/event/publisher/user/mock"
	userentitystub "github.com/BenyaChef/micro/auth/test/domain/entity/user/stub"
	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

type UserEventPublisherShould struct {
	suite.Suite

	queuePublisher *usereventpublishermock.QueuePublisherMock
	publisher      *usereventpublisher.UserEventPublisher
}

func TestUserEventPublisherShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserEventPublisherShould))
}

func (s *UserEventPublisherShould) SetupTest() {
	s.queuePublisher = usereventpublishermock.NewQueuePublisherMock()

	publisher, err := usereventpublisher.NewBuilder().
		QueuePublisher(s.queuePublisher).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.publisher = publisher
}

func (s *UserEventPublisherShould) TestBuild_WithoutParams_ReturnError() {
	publisher, err := usereventpublisher.NewBuilder().Build()

	assert.Nil(s.T(), publisher)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, usereventpublisher.ErrQueuePublisherIsRequired))
	assert.True(s.T(), errors.Is(err, usereventpublisher.ErrLogPublisherIsRequired))
}

func (s *UserEventPublisherShould) TestBuild_TopicDefaultsToUser() {
	published := s.publishOne()

	assert.Equal(s.T(), commonautheventmodel.TopicUser, published.Topic)
	assert.Equal(s.T(), "user", published.Topic)
}

func (s *UserEventPublisherShould) publishOne() *queueservermodel.ProducerMessage {
	s.T().Helper()

	userID := commonuserentity.NewUserID()

	require.NoError(s.T(), s.publisher.Registered(s.T().Context(), userID, userentitystub.GetEmail()))
	require.Len(s.T(), s.queuePublisher.Sent, 1)

	return s.queuePublisher.Sent[0]
}

func (s *UserEventPublisherShould) TestRegistered_ZeroUserID_DoesNotPublish() {
	err := s.publisher.Registered(s.T().Context(), commonuserentity.UserID{}, userentitystub.GetEmail())

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, usereventpublisher.ErrUserIDIsRequired))
	assert.Empty(s.T(), s.queuePublisher.Sent)
}

func (s *UserEventPublisherShould) TestRegistered_KeyIsUserID() {
	userID := commonuserentity.NewUserID()

	require.NoError(s.T(), s.publisher.Registered(s.T().Context(), userID, userentitystub.GetEmail()))
	require.Len(s.T(), s.queuePublisher.Sent, 1)

	assert.Equal(s.T(), userID.String(), s.queuePublisher.Sent[0].PartitionKey)
}

func (s *UserEventPublisherShould) TestRegistered_PayloadMatchesContract() {
	userID := commonuserentity.NewUserID()
	before := time.Now().Add(-time.Second)

	require.NoError(s.T(), s.publisher.Registered(s.T().Context(), userID, userentitystub.GetEmail()))
	require.Len(s.T(), s.queuePublisher.Sent, 1)

	var raw map[string]any
	require.NoError(s.T(), json.Unmarshal(s.queuePublisher.Sent[0].Value, &raw))

	assert.Equal(s.T(), "REGISTERED", raw["action"])
	assert.NotEmpty(s.T(), raw["event_id"])
	assert.NotEmpty(s.T(), raw["occurred_at"])

	data, ok := raw["data"].(map[string]any)
	require.True(s.T(), ok)
	assert.Equal(s.T(), userID.String(), data["user_id"])
	assert.Equal(s.T(), userentitystub.RawEmail, data["email"])

	event, err := commonautheventmodel.UserEventFromBytes(s.queuePublisher.Sent[0].Value)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), commonautheventmodel.UserEventActionRegistered, event.Action)
	assert.Equal(s.T(), userID.String(), event.PartitionKey())
	assert.True(s.T(), event.OccurredAt.After(before))
}

func (s *UserEventPublisherShould) TestRegistered_QueueFails_WrapError() {
	cause := errors.New("kafka is down")
	s.queuePublisher.SendMessageFunc = func(context.Context, *queueservermodel.ProducerMessage) error {
		return cause
	}

	err := s.publisher.Registered(s.T().Context(), commonuserentity.NewUserID(), userentitystub.GetEmail())

	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, cause)
	assert.True(s.T(), apperrors.EqualByCode(err, usereventpublisher.ErrCodePublishFailed))
}

func (s *UserEventPublisherShould) TestRegistered_TwoEvents_HaveDifferentEventIDs() {
	email := userentitystub.GetEmail()

	require.NoError(s.T(), s.publisher.Registered(s.T().Context(), commonuserentity.NewUserID(), email))
	require.NoError(s.T(), s.publisher.Registered(s.T().Context(), commonuserentity.NewUserID(), email))
	require.Len(s.T(), s.queuePublisher.Sent, 2)

	first, err := commonautheventmodel.UserEventFromBytes(s.queuePublisher.Sent[0].Value)
	require.NoError(s.T(), err)

	second, err := commonautheventmodel.UserEventFromBytes(s.queuePublisher.Sent[1].Value)
	require.NoError(s.T(), err)

	assert.NotEqual(s.T(), first.EventID, second.EventID)
}
