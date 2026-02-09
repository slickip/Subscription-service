package main

import (
	"log"
	"net/http"
	"time"

	"github.com/slickip/Subscription-service/internal/config"
	"github.com/slickip/Subscription-service/internal/db"
	"github.com/slickip/Subscription-service/internal/handlers"
)

func main() {
	cfg := config.Load()
	dbConn := db.New(cfg)
	mux := http.NewServeMux()

	subscriptionHandler := &handlers.SubscriptionHandler{
		DB: dbConn,
	}

	mux.Handle("/subscriptions", subscriptionHandler)
	mux.Handle("/subscriptions/", subscriptionHandler)
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("Starting subscription service on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("completed %s %s in %s", r.Method, r.URL.Path, time.Since(start))
	})
}
