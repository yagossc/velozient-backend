package api

import "net/http"

// CORSMiddleware implements an
// http middleware for setting headers
// for every http request received.
type CORSMiddleware struct {
	handler http.Handler
}

// WithCORSMiddleware wraps a given http handler with the header middleware.
func WithCORSMiddleware(handler http.Handler) *CORSMiddleware {
	return &CORSMiddleware{
		handler: handler,
	}
}

// ServeHTTP implements the needed http handler interface.
func (middleware *CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*") // ideally this wouldn't allow *
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

	// This exists because of the browser fetching
	// for CORS related headers and such.
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	middleware.handler.ServeHTTP(w, r)
}
