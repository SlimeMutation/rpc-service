package httputil

import "net/http"

type WrappedResponseWriter struct {
	StatusCode  int
	ResponseLen int

	w           http.ResponseWriter
	wroteHeader bool
}

func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	return &WrappedResponseWriter{
		StatusCode: http.StatusOK,
		w:          w,
	}
}

func (w *WrappedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *WrappedResponseWriter) Write(bytes []byte) (int, error) {
	n, err := w.w.Write(bytes)
	w.ResponseLen += n
	return n, err
}

func (w *WrappedResponseWriter) WriteHeader(StatusCode int) {
	if !w.wroteHeader {
		w.wroteHeader = true
		w.StatusCode = StatusCode
		w.w.WriteHeader(StatusCode)
	}
}
