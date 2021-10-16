package server

import (
	"net/http"
	"os"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		if os.Getenv("APP_ENV") == "production" {
			w.Header().Set("Access-Control-Allow-Origin", "https://bulb-audiy.uc.r.appspot.com")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
