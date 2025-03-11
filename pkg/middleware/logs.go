package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode: 0,
		}
		next.ServeHTTP(wrapper, r)
		log.Println(r.Method, wrapper.StatusCode, r.URL.Path, time.Since(start))
	})
}
