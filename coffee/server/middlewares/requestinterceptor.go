package middlewares

import (
	"bytes"
	"coffee/helpers"
	"io"
	"net/http"
)

func RequestInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer that captures the response content
		rw := &helpers.ResponseWriterWithContent{ResponseWriter: w, ResponseBodyCopy: &bytes.Buffer{}}
		var err error
		rw.InputBodyCopy, err = io.ReadAll(r.Body)
		if err == nil {
			if len(rw.InputBodyCopy) > 0 {
				r.Body = io.NopCloser(bytes.NewBuffer(rw.InputBodyCopy))
			}
		}
		next.ServeHTTP(rw, r)
	})
}
