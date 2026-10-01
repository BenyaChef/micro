package envregistry

import (
	"errors"
	"strconv"
	"time"

	envregistryinterface "github.com/BenyaChef/micro/infrastructure/envregistry/interface"
	envregistrymodel "github.com/BenyaChef/micro/infrastructure/envregistry/model"
)

var _ envregistryinterface.EnvRegistry = (*Registry)(nil)

type Registry struct {
	source   envregistrymodel.Source
	problems []error
}

func (r *Registry) String(key envregistrymodel.EnvKey) string {
	value := r.source(key)
	if value == "" {
		r.problems = append(r.problems, ErrEnvNotFound(key))
	}

	return value
}

func (r *Registry) StringDefault(key envregistrymodel.EnvKey, fallback string) string {
	value := r.source(key)
	if value == "" {
		return fallback
	}

	return value
}

func (r *Registry) Int(key envregistrymodel.EnvKey, fallback int) int {
	raw := r.source(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		r.problems = append(r.problems, ErrInvalidEnvValue(key, raw, "int"))

		return fallback
	}

	return value
}

func (r *Registry) Duration(key envregistrymodel.EnvKey, fallback time.Duration) time.Duration {
	raw := r.source(key)
	if raw == "" {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		r.problems = append(r.problems, ErrInvalidEnvValue(key, raw, "duration"))

		return fallback
	}

	return value
}

func (r *Registry) Bool(key envregistrymodel.EnvKey, fallback bool) bool {
	raw := r.source(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		r.problems = append(r.problems, ErrInvalidEnvValue(key, raw, "bool"))

		return fallback
	}

	return value
}

func (r *Registry) Err() error {
	if len(r.problems) == 0 {
		return nil
	}

	return errors.Join(r.problems...)
}
