package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v76"
	stripeSession "github.com/stripe/stripe-go/v76/checkout/session"
)

type CheckoutHandler struct {
	logger   *log.Logger
	baseURL  string
	products map[string]checkoutProduct
}

type checkoutProduct struct {
	Name         string
	SetupPriceID string
	SubPriceID   string
	ProductKey   string
}

func NewCheckoutHandler(logger *log.Logger) *CheckoutHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	return &CheckoutHandler{
		logger:  logger,
		baseURL: os.Getenv("BASE_URL"),
		products: map[string]checkoutProduct{
			"webchat": {
				Name:         "AI Webchat Agent",
				SetupPriceID: os.Getenv("STRIPE_PRICE_WEBCHAT_SETUP"),
				SubPriceID:   os.Getenv("STRIPE_PRICE_WEBCHAT_MONTHLY"),
				ProductKey:   "webchatagent",
			},
			"reviews": {
				Name:         "Review Response Automation",
				SetupPriceID: os.Getenv("STRIPE_PRICE_REVIEWS_SETUP"),
				SubPriceID:   os.Getenv("STRIPE_PRICE_REVIEWS_MONTHLY"),
				ProductKey:   "reviewsagent",
			},
		},
	}
}

// StartCheckout — GET /checkout/:product
// Redirects to Stripe Checkout for the given product.
func (h *CheckoutHandler) StartCheckout(w http.ResponseWriter, r *http.Request) {
	productSlug := chi.URLParam(r, "product")
	product, ok := h.products[productSlug]

	if !ok {
		http.NotFound(w, r)
		return
	}

	if product.SetupPriceID == "" || product.SubPriceID == "" {
		h.logger.Printf("ERROR: Stripe price IDs not configured for product: %s", productSlug)
		http.Error(w, "Product not available", http.StatusServiceUnavailable)
		return
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(product.SetupPriceID),
				Quantity: stripe.Int64(1),
			},
			{
				Price:    stripe.String(product.SubPriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(fmt.Sprintf(
			"%s/checkout/success?session_id={CHECKOUT_SESSION_ID}&product=%s",
			h.baseURL, productSlug,
		)),
		CancelURL: stripe.String(fmt.Sprintf("%s/products", h.baseURL)),
		Metadata: map[string]string{
			"product_key":  product.ProductKey,
			"product_name": product.Name,
		},
		BillingAddressCollection: stripe.String("required"),
		AllowPromotionCodes:      stripe.Bool(true),
	}

	s, err := stripeSession.New(params)
	if err != nil {
		h.logger.Printf("ERROR: Stripe session creation failed: %v", err)
		http.Error(w, "Unable to process checkout", http.StatusInternalServerError)
		return
	}

	h.logger.Printf("INFO: Stripe checkout session created for %s: %s", product.Name, s.ID)
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
}
