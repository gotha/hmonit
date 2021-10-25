package main

import (
	"net/http"

	"github.com/Financial-Times/gourmet/log"
)

type RequestLoggingMiddleware struct {
	logger *log.StructuredLogger
}

func NewRequestLoggingMiddleware(l *log.StructuredLogger) *RequestLoggingMiddleware {
	return &RequestLoggingMiddleware{logger: l}
}

func (m *RequestLoggingMiddleware) Log(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info("request received",
			log.WithField("method", r.Method),
			log.WithField("path", r.URL.Path),
		)
		next(w, r)
	}
}
