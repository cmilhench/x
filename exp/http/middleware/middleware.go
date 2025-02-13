package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func CreateStack(in ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, v := range in {
			next = v(next)
		}
		return next
	}
}
