package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/invoice"
	"github.com/stripe/stripe-go/v76/invoiceitem"
)

type BillingHandler struct {
	logger *log.Logger
}

func NewBillingHandler(logger *log.Logger) *BillingHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &BillingHandler{logger: logger}
}

type CreateInvoiceRequest struct {
	ClientName    string   `json:"client_name"`
	ClientEmail   string   `json:"client_email"`
	Products      []string `json:"products"`
	SetupFeeTotal int64    `json:"setup_fee_total"`
	MonthlyTotal  int64    `json:"monthly_total"`
	DueDate       int64    `json:"due_date"`
	Description   string   `json:"description"`
}

// CreateInvoice — POST /api/billing/create-invoice
// Called by Blueprint Command after a proposal is signed.
// Protected by X-BPA-Admin-Key header.
func (h *BillingHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-BPA-Admin-Key") != os.Getenv("BPA_INTERNAL_KEY") {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	// Create Stripe customer
	cust, err := customer.New(&stripe.CustomerParams{
		Name:  stripe.String(req.ClientName),
		Email: stripe.String(req.ClientEmail),
	})
	if err != nil {
		h.logger.Printf("ERROR: Stripe customer creation failed: %v", err)
		http.Error(w, `{"error":"stripe error"}`, http.StatusInternalServerError)
		return
	}

	// Create invoice item for setup fee
	_, err = invoiceitem.New(&stripe.InvoiceItemParams{
		Customer:    stripe.String(cust.ID),
		Amount:      stripe.Int64(req.SetupFeeTotal),
		Currency:    stripe.String("usd"),
		Description: stripe.String(fmt.Sprintf("Setup fee — %s", req.Description)),
	})
	if err != nil {
		h.logger.Printf("ERROR: Invoice item creation failed: %v", err)
		http.Error(w, `{"error":"stripe error"}`, http.StatusInternalServerError)
		return
	}

	// Create and finalize the invoice
	inv, err := invoice.New(&stripe.InvoiceParams{
		Customer:         stripe.String(cust.ID),
		CollectionMethod: stripe.String("send_invoice"),
		DaysUntilDue:     stripe.Int64(7),
		AutoAdvance:      stripe.Bool(true),
	})
	if err != nil {
		h.logger.Printf("ERROR: Invoice creation failed: %v", err)
		http.Error(w, `{"error":"stripe error"}`, http.StatusInternalServerError)
		return
	}

	inv, err = invoice.FinalizeInvoice(inv.ID, nil)
	if err != nil {
		h.logger.Printf("ERROR: Invoice finalization failed: %v", err)
		http.Error(w, `{"error":"stripe error"}`, http.StatusInternalServerError)
		return
	}
	invoice.SendInvoice(inv.ID, nil)

	h.logger.Printf("INFO: Invoice created and sent: %s to %s", inv.ID, req.ClientEmail)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"invoice_id":  inv.ID,
		"customer_id": cust.ID,
		"status":      string(inv.Status),
		"hosted_url":  inv.HostedInvoiceURL,
	})
}
