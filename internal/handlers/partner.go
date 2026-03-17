package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cocacolasante/bpa-inc-website/internal/db"
	"github.com/cocacolasante/bpa-inc-website/internal/models"
	"github.com/cocacolasante/bpa-inc-website/internal/services"
)

type PartnerHandler struct {
	logger *log.Logger
	email  *services.EmailService
}

func NewPartnerHandler(logger *log.Logger) *PartnerHandler {
	return &PartnerHandler{
		logger: logger,
		email:  services.NewEmailService(),
	}
}

func (h *PartnerHandler) Apply(w http.ResponseWriter, r *http.Request) {
	var app models.PartnerApplication
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		h.sendJSON(w, http.StatusBadRequest, false, "Invalid request")
		return
	}

	if err := app.Validate(); err != nil {
		h.sendJSON(w, http.StatusBadRequest, false, err.Error())
		return
	}

	// Save to DB
	if db.DB != nil {
		_, err := db.DB.Exec(`
			INSERT INTO partner_applications
			(company_name, contact_name, contact_email, contact_phone, website,
			 years_in_business, client_count, expected_volume, why_partner)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			app.CompanyName, app.ContactName, app.ContactEmail, app.ContactPhone,
			app.Website, app.YearsInBusiness, app.ClientCount,
			app.ExpectedVolume, app.WhyPartner,
		)
		if err != nil {
			h.logger.Printf("WARN: Failed to save partner application to DB: %v", err)
		}
	}

	// Notify BPA team
	go func() {
		if err := h.email.SendPartnerApplicationAlert(&app); err != nil {
			h.logger.Printf("WARN: Partner application alert email failed: %v", err)
		}
		services.SendDiscordAlert(fmt.Sprintf(
			"🤝 New partner application: **%s** (%s) — %s/month expected",
			app.CompanyName, app.ContactEmail, app.ExpectedVolume,
		))
	}()

	// Auto-reply to applicant
	go func() {
		if err := h.email.SendPartnerApplicationConfirmation(&app); err != nil {
			h.logger.Printf("WARN: Partner application confirmation email failed: %v", err)
		}
	}()

	h.logger.Printf("INFO: Partner application received from %s — %s",
		app.ContactName, app.CompanyName)
	h.sendJSON(w, http.StatusOK, true,
		"Application received! We'll review it and get back to you within 2 business days.")
}

func (h *PartnerHandler) sendJSON(w http.ResponseWriter, status int, success bool, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"message": msg,
	})
}
