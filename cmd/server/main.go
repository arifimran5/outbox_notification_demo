package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/arifimran5/outbox_notification_demo/internal/database"
	"github.com/arifimran5/outbox_notification_demo/internal/events"
	"github.com/arifimran5/outbox_notification_demo/internal/handlers"
	"github.com/arifimran5/outbox_notification_demo/internal/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("couldn't load env variables")
		return
	}

	// 1. Initialize DB
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer database.CloseDB()

	// 2. Initialize SSE
	events.InitSSE()

	// 3. Start Outbox Processor (Background Context)
	ctx, cancel := context.WithCancel(context.Background())
	events.StartOutboxProcessor(ctx)

	// CORS Setup
	webUrl := os.Getenv("APP_URL")
	webUrlArr := strings.Split(webUrl, ",")
	corsObj := gohandlers.CORS(
		gohandlers.AllowedOrigins(webUrlArr), // In prod, be specific
		gohandlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}),
		gohandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// 4. Setup Router
	r := mux.NewRouter()

	// Public Routes
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")

	// Protected Routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/topics", handlers.GetTopics).Methods("GET")
	api.HandleFunc("/topics/{topicId}/subscribe", handlers.Subscribe).Methods("POST")
	api.HandleFunc("/topics/{topicId}/unsubscribe", handlers.Unsubscribe).Methods("POST")
	api.HandleFunc("/subscriptions", handlers.GetUserSubscriptions).Methods("GET")
	api.HandleFunc("/topics/{topicId}/posts", handlers.CreatePost).Methods("POST")
	api.HandleFunc("/topics/{topicId}/posts", handlers.GetTopicPosts).Methods("GET")

	// SSE Endpoint (Protected)
	api.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(string)
		events.Manager.Subscribe(userID, w, r)
	})

	port := os.Getenv("SERVER_PORT")

	// 5. Server Setup with Graceful Shutdown
	srv := &http.Server{
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: 0,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 120,
		Handler:      corsObj(r),
	}

	// Run server in goroutine
	go func() {
		log.Printf("Server starting on :%v", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shutdown logic
	cancel() // Stop the Outbox processor

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer shutdownCancel()

	srv.Shutdown(shutdownCtx)
	log.Println("Server shutting down gracefully")
}
