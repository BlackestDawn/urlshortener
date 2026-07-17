package main

import (
	"net/http"
	"time"
)

// newTimeoutHandler bounds request handling time. If h does not finish
// within dt, the client receives a 503 and further writes by h fail with
// http.ErrHandlerTimeout, but h keeps running until it respects context
// cancellation itself.
func newTimeoutHandler(h http.Handler, dt time.Duration) http.Handler {
	return http.TimeoutHandler(h, dt, "request timed out")
}
