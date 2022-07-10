package middleware

import (
	"net/http"
)

func OS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(SetOSCtx(r.Context(), r.UserAgent()))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
