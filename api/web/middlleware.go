package web

import (
	"net/http"
	"slices"
)

type MiddlewareFunc func(next http.Handler) http.Handler

func Chain(middlewareChain ...MiddlewareFunc) MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		var next http.Handler = handler

		// need to iterate backward since next should always be the Handler returned
		// by the middleware after it.
		for _, mw := range slices.Backward(middlewareChain) {
			next = mw(next)
		}

		return next
	}
}

func (mf MiddlewareFunc) Complete(end http.Handler) http.Handler {
	return mf(end)
}
