package middlewares

import (
	"coffee/helpers"
	"net/http"
)

func SentryErrorLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
		customWriter := w.(*helpers.ResponseWriterWithContent)
		helpers.GetSentryEvent(customWriter, r, customWriter.InputBodyCopy)
	})
}
