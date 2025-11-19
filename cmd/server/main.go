package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cocacolasante/bpa-inc-website/internal/handlers"
	"github.com/cocacolasante/bpa-inc-website/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := log.New(os.Stdout, "[BPA-Website] ", log.LstdFlags|log.Lshortfile)
	logger.Println("INFO: Starting Blueprint Automations website server...")

	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Initialize handlers
	pageHandler := handlers.NewPageHandler(logger)
	contactHandler := handlers.NewContactHandler(logger)

	// Static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Routes
	r.Get("/", pageHandler.Home)
	r.Get("/about", pageHandler.About)
	r.Get("/services", pageHandler.Services)
	r.Get("/portfolio", pageHandler.Portfolio)
	r.Get("/contact", pageHandler.Contact)

	// API routes
	r.Post("/api/contact", contactHandler.Submit)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Printf("INFO: Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Fatalf("ERROR: Server failed to start: %v", err)
	}
}
