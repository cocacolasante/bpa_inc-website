package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cocacolasante/bpa-inc-website/internal/db"
	"github.com/cocacolasante/bpa-inc-website/internal/handlers"
	"github.com/cocacolasante/bpa-inc-website/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignored if not present)
	godotenv.Load()

	logger := log.New(os.Stdout, "[BPA-Website] ", log.LstdFlags|log.Lshortfile)
	logger.Println("INFO: Starting Blueprint Automations website server...")

	// Initialize database (optional — runs without DB if DATABASE_URL not set)
	db.Init(logger)

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
	billingHandler := handlers.NewBillingHandler(logger)
	webhookHandler := handlers.NewWebhookHandler(logger)
	checkoutHandler := handlers.NewCheckoutHandler(logger)
	partnerHandler := handlers.NewPartnerHandler(logger)

	// Static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Page routes
	r.Get("/", pageHandler.Home)
	r.Get("/about", pageHandler.About)
	r.Get("/services", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/products", http.StatusMovedPermanently)
	})
	r.Get("/portfolio", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/products", http.StatusMovedPermanently)
	})
	r.Get("/contact", pageHandler.Contact)
	r.Get("/audit", pageHandler.Audit)
	r.Get("/free-audit", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/audit", http.StatusMovedPermanently)
	})
	r.Get("/products", pageHandler.Products)
	r.Get("/pricing", pageHandler.Pricing)
	r.Get("/partners", pageHandler.Partners)
	r.Get("/checkout/success", pageHandler.CheckoutSuccess)

	// API routes
	r.Post("/api/contact", contactHandler.Submit)
	r.Post("/api/partner/apply", partnerHandler.Apply)

	// Checkout
	r.Get("/checkout/{product}", checkoutHandler.StartCheckout)

	// Internal + webhooks
	r.Post("/api/billing/create-invoice", billingHandler.CreateInvoice)
	r.Post("/api/webhooks/stripe", webhookHandler.StripeWebhook)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"bpa-website"}`))
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
