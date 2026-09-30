package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 3)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RateLimiter(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {

		w.Header().Set("Retry-After", "1")
		writeJSON(w, http.StatusTooManyRequests, map[string]string{
			"error": "too many requests",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "hello! your request was allowed",
	})
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", RateLimiter)

	fmt.Println("Server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		slog.Error("error found", slog.String("error", err.Error()))
	}
}
