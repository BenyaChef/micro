package envregistryinterface

import (
	"time"

	envregistrymodel "micro/infrastructure/envregistry/model"
)

type EnvRegistry interface {
	String(key envregistrymodel.EnvKey) string
	StringDefault(key envregistrymodel.EnvKey, fallback string) string
	Int(key envregistrymodel.EnvKey, fallback int) int
	Duration(key envregistrymodel.EnvKey, fallback time.Duration) time.Duration
	Bool(key envregistrymodel.EnvKey, fallback bool) bool
	Err() error
}
