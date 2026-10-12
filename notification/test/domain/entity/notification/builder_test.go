package notification_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
	notificationentitystub "github.com/BenyaChef/micro/notification/test/domain/entity/notification/stub"
)

type BuilderShould struct {
	suite.Suite
}

func TestBuilderShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuilderShould))
}

func (s *BuilderShould) TestBuild_WithoutParams_ReturnAllErrors() {
	notification, err := notificationentity.NewBuilder().Build()

	assert.Nil(s.T(), notification)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, notificationentity.ErrNotificationIDIsRequired))
	assert.True(s.T(), errors.Is(err, notificationentity.ErrEventIDIsRequired))
	assert.True(s.T(), errors.Is(err, notificationentity.ErrUserIDIsRequired))
	assert.True(s.T(), errors.Is(err, notificationentity.ErrKindIsRequired))
	assert.True(s.T(), errors.Is(err, notificationentity.ErrPayloadIsRequired))
}

func (s *BuilderShould) TestBuild_AllParams_FillsSentAt() {
	notification := notificationentitystub.GetValidNotification()

	assert.False(s.T(), notification.SentAt().IsZero())
	assert.Equal(s.T(), notificationentity.KindWelcome, notification.Kind())
}

func (s *BuilderShould) TestKindFrom_Unsupported_ReturnError() {
	kind, err := notificationentity.KindFrom("telegram")

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, notificationentity.ErrCodeKindUnsupported))
	assert.True(s.T(), kind.IsZero())
}

func (s *BuilderShould) TestKindFrom_Empty_ReturnRequiredError() {
	kind, err := notificationentity.KindFrom("")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, notificationentity.ErrKindIsRequired))
	assert.True(s.T(), kind.IsZero())
}

func (s *BuilderShould) TestPayloadFrom_NotJSON_ReturnError() {
	for _, raw := range [][]byte{nil, {}, []byte("not json")} {
		payload, err := notificationentity.PayloadFrom(raw)

		require.Error(s.T(), err)
		assert.True(s.T(), errors.Is(err, notificationentity.ErrPayloadIsEmpty))
		assert.True(s.T(), payload.IsZero())
	}
}

func (s *BuilderShould) TestEventIDFrom_NotUUID_ReturnError() {
	eventID, err := notificationentity.EventIDFrom("not-a-uuid")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, notificationentity.ErrEventIDInvalidFormat))
	assert.True(s.T(), eventID.IsZero())
}
