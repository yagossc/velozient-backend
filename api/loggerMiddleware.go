package api

import (
	"fmt"
	"net/http"
	"time"
)

// format definitions for the log line
const logFormat = "[%s] - [%s] - [%s]  - [status=%d]\n"
const timeFormat = "2006-01-02 15:04:05"

// LoggerMiddleware implements a
// VERY NAIVE approach of an
// http middleware for logging
// http requests receveid.
type LoggerMiddleware struct {
	handler        http.Handler
	rw             http.ResponseWriter
	responseStatus int
}

// Wrap http's ResponseWrite Write method.
func (l *LoggerMiddleware) Write(b []byte) (int, error) {
	return l.rw.Write(b)
}

// Wrap http's ResponseWrite WriteHeader method.
func (l *LoggerMiddleware) WriteHeader(statusCode int) {
	l.responseStatus = statusCode
	l.rw.WriteHeader(statusCode)
}

// Wrap http's ResponseWrite Header method.
func (l *LoggerMiddleware) Header() http.Header {
	return l.rw.Header()
}

// WithLoggerMiddleware wraps a given http handler with the logger middleware.
func WithLoggerMiddleware(handler http.Handler) *LoggerMiddleware {
	return &LoggerMiddleware{
		handler: handler,
	}
}

// ServeHTTP implements the needed http handler interface.
func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// very naive approach
	l.rw = w
	l.handler.ServeHTTP(l, r)

	// log the desired fields
	fmt.Printf(logFormat,
		time.Now().Format(timeFormat),
		r.Method, r.RequestURI, l.responseStatus)
}
