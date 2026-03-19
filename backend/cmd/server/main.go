package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"web-go-template/internal/db"
	"web-go-template/internal/handler"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/wordle?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	queries := db.New(pool)
	h := handler.New(queries)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api/wordle", func(r chi.Router) {
		r.Get("/today", h.GetTodayWord)
		r.Post("/guess", h.GuessWord)
		r.Post("/scores", h.SubmitScore)
		r.Get("/scores", h.GetScores)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server starting on :%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
