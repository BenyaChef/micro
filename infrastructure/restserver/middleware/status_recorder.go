package restservermiddleware

import "net/http"

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func newStatusRecorder(writer http.ResponseWriter) *statusRecorder {
	return &statusRecorder{ResponseWriter: writer, status: http.StatusOK}
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Status() int {
	return r.status
}
