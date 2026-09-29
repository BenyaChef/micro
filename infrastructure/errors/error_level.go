package errors

type ErrorLevel string

const (
	LevelInfo     ErrorLevel = "info"
	LevelWarn     ErrorLevel = "warn"
	LevelError    ErrorLevel = "error"
	LevelCritical ErrorLevel = "critical"
)

func (l ErrorLevel) String() string {
	return string(l)
}
