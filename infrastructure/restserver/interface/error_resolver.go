package restserverinterface

type ErrorResolver interface {
	GetErrorCode(err error) string
	GetErrorText(err error) string
	GetHTTPCode(err error) int
}
