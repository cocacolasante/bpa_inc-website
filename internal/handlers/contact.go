package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cocacolasante/bpa-inc-website/internal/models"
)

type ContactHandler struct {
	logger *log.Logger
}

func NewContactHandler(logger *log.Logger) *ContactHandler {
	return &ContactHandler{
		logger: logger,
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

	// Log the successful submission
	h.logger.Printf("INFO: New contact form submission received")
	h.logger.Printf("INFO: Name: %s", submission.FullName)
	h.logger.Printf("INFO: Business: %s", submission.BusinessName)
	h.logger.Printf("INFO: Email: %s", submission.Email)
	h.logger.Printf("INFO: Phone: %s", submission.Phone)
	h.logger.Printf("INFO: Website: %s", submission.Website)
	h.logger.Printf("INFO: Message: %s", submission.Message)

	// TODO: In production, you would:
	// 1. Save to database
	// 2. Send email notification
	// 3. Integrate with CRM
	// 4. Send auto-response email

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
