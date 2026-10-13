package otel_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	infraotel "github.com/BenyaChef/micro/infrastructure/otel"
)

type ProviderShould struct {
	suite.Suite

	provider *infraotel.Provider
}

func TestProviderShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProviderShould))
}

func (s *ProviderShould) SetupTest() {
	provider, err := infraotel.NewBuilder().
		ServiceName("auth").
		Environment("test").
		MetricsPort("9101").
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build(s.T().Context())
	require.NoError(s.T(), err)

	s.provider = provider
	s.T().Cleanup(func() { provider.Close(context.Background()) })
}

func (s *ProviderShould) scrape() string {
	s.T().Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(s.T().Context(), "GET", "/metrics", nil)

	s.provider.MetricsServer().Handler().ServeHTTP(recorder, request)
	require.Equal(s.T(), 200, recorder.Code)

	return recorder.Body.String()
}

func (s *ProviderShould) TestBuild_WithoutParams_ReturnError() {
	provider, err := infraotel.NewBuilder().MetricsPort("").Build(s.T().Context())

	assert.Nil(s.T(), provider)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, infraotel.ErrServiceNameIsRequired))
	assert.True(s.T(), errors.Is(err, infraotel.ErrLogPublisherIsRequired))
	assert.True(s.T(), errors.Is(err, infraotel.ErrPortIsRequired))
}

func (s *ProviderShould) TestHTTPMonitor_RecordsDurationWithRouteAndStatus() {
	monitor, err := infraotel.NewHTTPMonitor(s.provider.MeterProvider())
	require.NoError(s.T(), err)

	monitor.RecordRequest(s.T().Context(), "GET", "/posts/{postID}", 404, 0.012)

	body := s.scrape()

	assert.Contains(s.T(), body, "http_server_request_duration_seconds")
	assert.Contains(s.T(), body, `http_route="/posts/{postID}"`)
	assert.Contains(s.T(), body, `http_response_status_code="404"`)
	assert.Contains(s.T(), body, `http_request_method="GET"`)
}

func (s *ProviderShould) TestMonitors_AreNilSafe() {
	var (
		httpMonitor  *infraotel.HTTPMonitor
		grpcMonitor  *infraotel.GRPCMonitor
		kafkaMonitor *infraotel.KafkaMonitor
	)

	assert.NotPanics(s.T(), func() {
		httpMonitor.RecordRequest(s.T().Context(), "GET", "/", 200, 0.1)
		grpcMonitor.RecordRequest(s.T().Context(), "/Svc/Method", "OK", 0.1)
		kafkaMonitor.RecordPublish(s.T().Context(), "user", infraotel.ResultSuccess)
		kafkaMonitor.RecordProcess(s.T().Context(), "user", infraotel.ResultSuccess, 0.1)
	})
}

func (s *ProviderShould) TestGRPCMonitor_RecordsCode() {
	monitor, err := infraotel.NewGRPCMonitor(s.provider.MeterProvider())
	require.NoError(s.T(), err)

	monitor.RecordRequest(s.T().Context(), "/micro.auth.schema.AuthService/ValidateToken", "OK", 0.003)

	body := s.scrape()

	assert.Contains(s.T(), body, "rpc_server_duration_seconds")
	assert.Contains(s.T(), body, `rpc_grpc_status_code="OK"`)
}

func (s *ProviderShould) TestKafkaMonitor_RecordsPublishAndProcess() {
	monitor, err := infraotel.NewKafkaMonitor(s.provider.MeterProvider())
	require.NoError(s.T(), err)

	monitor.RecordPublish(s.T().Context(), "user", infraotel.ResultSuccess)
	monitor.RecordProcess(s.T().Context(), "user", infraotel.ResultError, 0.5)

	body := s.scrape()

	assert.Contains(s.T(), body, "messaging_publish_count")
	assert.Contains(s.T(), body, "messaging_process_duration_seconds")
	assert.Contains(s.T(), body, `messaging_destination_name="user"`)
	assert.Contains(s.T(), body, `messaging_operation_result="error"`)
}

// Имя сервиса на сами метрики не навешивается: в дашбордах сервис различается
// меткой job из конфига скрейпа. Но target_info оставлен, чтобы по эндпоинту
// было видно, чей он.
func (s *ProviderShould) TestScrape_ExposesServiceNameViaTargetInfo() {
	body := s.scrape()

	assert.Contains(s.T(), body, "target_info")
	assert.Contains(s.T(), body, `service_name="auth"`)
	assert.Contains(s.T(), body, `deployment_environment_name="test"`)
}
