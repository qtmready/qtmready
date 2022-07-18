package middlewares

import "net/http"

func ContentTypeJSON(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
