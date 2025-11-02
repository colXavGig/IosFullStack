package web

import (
	"net/http"
)

func RecoveryFromPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				err, ok := v.(error)
				if !ok {
					return
				}
				_ = JsonResponseWriter(w, &HttpError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
					Details: err.Error(),
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
