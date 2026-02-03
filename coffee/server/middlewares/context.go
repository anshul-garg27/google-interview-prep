package middlewares

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"context"
	"net/http"
)

func ApplicationContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		appCtx, _ := appcontext.CreateApplicationContext(ctx, r)
		newCtx := context.WithValue(ctx, constants.AppContextKey, &appCtx)
		r = r.WithContext(newCtx)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
