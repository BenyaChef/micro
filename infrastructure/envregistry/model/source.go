package envregistrymodel

import (
	"os"
	"strings"
)

type Source func(key EnvKey) string

func OSSource(key EnvKey) string {
	return strings.TrimSpace(os.Getenv(key.String()))
}
