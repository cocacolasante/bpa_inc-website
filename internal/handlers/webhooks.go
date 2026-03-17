package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cocacolasante/bpa-inc-website/internal/services"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

type WebhookHandler struct {
	logger *log.Logger
}

func NewWebhookHandler(logger *log.Logger) *WebhookHandler {
	return &WebhookHandler{logger: logger}
}

// StripeWebhook — POST /api/webhooks/stripe
// Handles: invoice.paid, checkout.session.completed, invoice.payment_failed
func (h *WebhookHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("ERROR: Webhook body read failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, err := webhook.ConstructEvent(
		body,
		r.Header.Get("Stripe-Signature"),
		os.Getenv("STRIPE_WEBHOOK_SECRET"),
	)
	if err != nil {
		h.logger.Printf("ERROR: Webhook signature verification failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Return 200 immediately, process async
	w.WriteHeader(http.StatusOK)

	go func() {
		switch event.Type {
		case "invoice.paid":
			var inv stripe.Invoice
			json.Unmarshal(event.Data.Raw, &inv)
			h.handleInvoicePaid(&inv)

		case "checkout.session.completed":
			var session stripe.CheckoutSession
			json.Unmarshal(event.Data.Raw, &session)
			h.handleCheckoutCompleted(&session)

		case "invoice.payment_failed":
			var inv stripe.Invoice
			json.Unmarshal(event.Data.Raw, &inv)
			h.handlePaymentFailed(&inv)
		}
	}()
}

func (h *WebhookHandler) handleInvoicePaid(inv *stripe.Invoice) {
	h.logger.Printf("INFO: Invoice paid: %s, customer: %s", inv.ID, inv.CustomerEmail)

	meta := inv.Metadata
	clientID := meta["client_id"]
	products := meta["products"]

	if clientID == "" || products == "" {
		h.logger.Printf("INFO: Invoice %s has no client_id metadata — skipping provisioning", inv.ID)
		return
	}

	h.triggerProvisioning(clientID, inv.CustomerEmail, products)
}

func (h *WebhookHandler) handleCheckoutCompleted(session *stripe.CheckoutSession) {
	h.logger.Printf("INFO: Checkout completed: %s, customer: %s", session.ID, session.CustomerEmail)

	productKey := session.Metadata["product_key"]
	if productKey == "" {
		h.logger.Printf("WARN: Checkout session %s missing product_key metadata", session.ID)
		return
	}

	var customerName string
	if session.CustomerDetails != nil {
		customerName = session.CustomerDetails.Name
	}

	h.triggerSingleProductProvisioning(session.CustomerEmail, customerName, productKey, session.ID)
}

func (h *WebhookHandler) handlePaymentFailed(inv *stripe.Invoice) {
	h.logger.Printf("WARN: Payment failed for invoice: %s, customer: %s", inv.ID, inv.CustomerEmail)
	services.SendDiscordAlert(fmt.Sprintf(
		"⚠️ Payment failed: %s — %s ($%.2f)",
		inv.CustomerEmail, inv.ID, float64(inv.AmountDue)/100,
	))
}

func (h *WebhookHandler) triggerProvisioning(clientID, email, products string) {
	commandURL := os.Getenv("BLUEPRINT_COMMAND_API_URL")
	adminKey := os.Getenv("BLUEPRINT_COMMAND_ADMIN_KEY")

	if commandURL == "" {
		h.logger.Println("WARN: BLUEPRINT_COMMAND_API_URL not set — skipping provisioning")
		return
	}

	payload, _ := json.Marshal(map[string]string{
		"client_id": clientID,
		"products":  products,
		"trigger":   "invoice_paid",
	})

	req, _ := http.NewRequest("POST",
		commandURL+"/api/internal/provision-from-payment", bytes.NewReader(payload))
	req.Header.Set("X-BPA-Admin-Key", adminKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		h.logger.Printf("ERROR: Provisioning trigger failed for client %s: %v", clientID, err)
		return
	}
	h.logger.Printf("INFO: Provisioning triggered for client %s", clientID)
}

func (h *WebhookHandler) triggerSingleProductProvisioning(email, name, productKey, sessionID string) {
	commandURL := os.Getenv("BLUEPRINT_COMMAND_API_URL")
	adminKey := os.Getenv("BLUEPRINT_COMMAND_ADMIN_KEY")

	if commandURL == "" {
		h.logger.Println("WARN: BLUEPRINT_COMMAND_API_URL not set")
		return
	}

	payload, _ := json.Marshal(map[string]string{
		"email":             email,
		"name":              name,
		"product_key":       productKey,
		"stripe_session_id": sessionID,
		"trigger":           "self_serve_checkout",
	})

	req, _ := http.NewRequest("POST",
		commandURL+"/api/internal/provision-self-serve", bytes.NewReader(payload))
	req.Header.Set("X-BPA-Admin-Key", adminKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		h.logger.Printf("ERROR: Self-serve provisioning failed for %s: %v", email, err)
		return
	}
	h.logger.Printf("INFO: Self-serve provisioning triggered for %s — %s", email, productKey)
}
