package envregistry

import (
	"fmt"

	envregistrymodel "github.com/BenyaChef/micro/infrastructure/envregistry/model"
	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeEnvNotFound     = errors.ErrorCode("c2af70f6-001")
	ErrCodeInvalidEnvValue = errors.ErrorCode("c2af70f6-002")
)

func ErrEnvNotFound(envKey envregistrymodel.EnvKey) error {
	msg := fmt.Sprintf("Environment variable is required: name=%s", envKey)

	return errors.NewErrorWithLevel(ErrCodeEnvNotFound, msg, errors.LevelCritical)
}

func ErrInvalidEnvValue(envKey envregistrymodel.EnvKey, value, expectedType string) error {
	msg := fmt.Sprintf(
		"Environment variable has invalid value: name=%s, value=%s, expected=%s",
		envKey, value, expectedType,
	)

	return errors.NewErrorWithLevel(ErrCodeInvalidEnvValue, msg, errors.LevelCritical)
}
