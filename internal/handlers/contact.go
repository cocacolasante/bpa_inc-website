package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cocacolasante/bpa-inc-website/internal/db"
	"github.com/cocacolasante/bpa-inc-website/internal/models"
	"github.com/cocacolasante/bpa-inc-website/internal/services"
)

type ContactHandler struct {
	logger *log.Logger
	email  *services.EmailService
}

func NewContactHandler(logger *log.Logger) *ContactHandler {
	return &ContactHandler{
		logger: logger,
		email:  services.NewEmailService(),
	}
}

type ContactResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h *ContactHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var submission models.ContactSubmission

	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		h.logger.Printf("ERROR: Failed to decode contact form: %v", err)
		h.sendJSONResponse(w, http.StatusBadRequest, false, "Invalid request format")
		return
	}

	if err := submission.Validate(); err != nil {
		h.logger.Printf("INFO: Contact form validation failed: %v", err)
		h.sendJSONResponse(w, http.StatusBadRequest, false, err.Error())
		return
	}

	h.logger.Printf("INFO: New contact form submission received")
	h.logger.Printf("INFO: Name: %s", submission.FullName)
	h.logger.Printf("INFO: Business: %s", submission.BusinessName)
	h.logger.Printf("INFO: Email: %s", submission.Email)
	h.logger.Printf("INFO: Phone: %s", submission.Phone)
	h.logger.Printf("INFO: Website: %s", submission.Website)
	h.logger.Printf("INFO: Source: %s", submission.Source)
	h.logger.Printf("INFO: Message: %s", submission.Message)

	// Save to DB (optional)
	if db.DB != nil {
		_, err := db.DB.Exec(
			`INSERT INTO leads (full_name, business_name, email, phone, website, message, source)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			submission.FullName, submission.BusinessName, submission.Email,
			submission.Phone, submission.Website, submission.Message, submission.Source,
		)
		if err != nil {
			h.logger.Printf("WARN: Failed to save lead to DB: %v", err)
		}
	}

	// Notify team and send auto-reply (non-blocking, errors logged not fatal)
	go func() {
		if err := h.email.SendLeadNotification(&submission); err != nil {
			h.logger.Printf("WARN: Lead notification email failed: %v", err)
		}
		services.SendDiscordLeadAlert(&submission)
	}()

	go func() {
		if err := h.email.SendAutoReply(&submission); err != nil {
			h.logger.Printf("WARN: Auto-reply email failed: %v", err)
		}
	}()

	h.sendJSONResponse(w, http.StatusOK, true, "Thank you for your message! We'll get back to you within 24 hours.")
}

func (h *ContactHandler) sendJSONResponse(w http.ResponseWriter, status int, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ContactResponse{
		Success: success,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
