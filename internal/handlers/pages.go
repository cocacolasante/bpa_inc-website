package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type PageHandler struct {
	logger    *log.Logger
	templates map[string]*template.Template
}

func NewPageHandler(logger *log.Logger) *PageHandler {
	templates := make(map[string]*template.Template)

	pages := []string{
		"home", "about", "services", "portfolio", "contact",
		"audit", "products", "partners", "checkout_success", "pricing",
	}

	for _, page := range pages {
		tmpl := template.Must(template.ParseFiles(
			filepath.Join("templates", "layouts", "base.html"),
			filepath.Join("templates", "pages", page+".html"),
		))
		templates[page] = tmpl
	}

	logger.Println("INFO: Page templates loaded successfully")

	return &PageHandler{
		logger:    logger,
		templates: templates,
	}
}

func (h *PageHandler) renderTemplate(w http.ResponseWriter, page string, data interface{}) {
	tmpl, ok := h.templates[page]
	if !ok {
		h.logger.Printf("ERROR: Template not found: %s", page)
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		h.logger.Printf("ERROR: Template execution failed for %s: %v", page, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *PageHandler) baseData(title, description, page string) map[string]interface{} {
	return map[string]interface{}{
		"Title":       title,
		"Description": description,
		"Page":        page,
		"CalComURL":   os.Getenv("CAL_COM_URL"),
	}
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := h.baseData(
		"Blueprint Automations - AI Workflow Automation for Growing Businesses",
		"We build clean, automated AI systems that let your business run with clarity and ease. Specializing in N8N, AI automations, and API integrations.",
		"home",
	)
	h.renderTemplate(w, "home", data)
}

func (h *PageHandler) About(w http.ResponseWriter, r *http.Request) {
	data := h.baseData(
		"About Us - Blueprint Automations",
		"Meet the team behind Blueprint Automations. Founded in 2025 by Scott Henry and Anthony Colasante, we help small to medium businesses leverage AI automation.",
		"about",
	)
	h.renderTemplate(w, "about", data)
}

func (h *PageHandler) Services(w http.ResponseWriter, r *http.Request) {
	data := h.baseData(
		"Our Services - AI Automation Solutions",
		"Professional AI workflow automation services. Custom workflows for $2,500 with ongoing support starting at $495/month.",
		"services",
	)
	h.renderTemplate(w, "services", data)
}

func (h *PageHandler) Portfolio(w http.ResponseWriter, r *http.Request) {
	data := h.baseData(
		"Portfolio - Our Work",
		"Explore our AI automation workflows and projects. See how we've helped businesses automate and scale.",
		"portfolio",
	)
	h.renderTemplate(w, "portfolio", data)
}

func (h *PageHandler) Contact(w http.ResponseWriter, r *http.Request) {
	data := h.baseData(
		"Contact Us - Let's Automate Your Business",
		"Ready to automate your business? Get in touch with Blueprint Automations for a consultation.",
		"contact",
	)
	h.renderTemplate(w, "contact", data)
}

func (h *PageHandler) Audit(w http.ResponseWriter, r *http.Request) {
	data := h.baseData(
		"Free Automation Audit — Blueprint Automations",
		"Find out how many hours per week your team could win back with AI automation. Free personalized report.",
		"audit",
	)
	data["ROIToolURL"] = template.URL(os.Getenv("ROI_TOOL_URL"))
	data["ROITenantID"] = os.Getenv("ROI_TENANT_ID")
	h.renderTemplate(w, "audit", data)
}

func (h *PageHandler) Products(w http.ResponseWriter, r *http.Request) {
	type ManagedProduct struct {
		Icon        string
		Name        string
		Description string
		SetupFee    string
		MonthlyFee  string
	}

	data := h.baseData(
		"AI Automation Products — Blueprint Automations",
		"Purpose-built AI automation products for every part of your business. Self-serve or fully managed.",
		"products",
	)

	data["ManagedProducts"] = []ManagedProduct{
		{"📧", "AI Email Triage", "Automatically sort, draft, and respond to your email — saving 8+ hours/week.", "$299", "$199"},
		{"📅", "Appointment Scheduling", "Full scheduling automation with reminders, no-shows, and waitlists.", "$397", "$149"},
		{"📋", "Lead Capture Pipeline", "Capture every lead, enrich them, and fire automated follow-up sequences.", "$497", "$149"},
		{"💰", "Invoice Automation", "Automated invoicing, payment reminders, and collections.", "$197", "$79"},
		{"📱", "Social Media Autopilot", "AI-generated social content posted across all your platforms.", "$297", "$197"},
		{"🤖", "AI Support Agent", "Handles tier-1 customer support via chat and email.", "$997", "$397"},
		{"📊", "Operations Dashboard", "Real-time business intelligence with daily AI briefings.", "$3000", "$497"},
		{"📞", "Voice AI Agent", "AI phone agent that answers calls and books appointments 24/7.", "$1500", "$397"},
		{"📄", "Proposal Generator", "AI-generated branded proposals delivered minutes after your discovery call.", "$597", "$197"},
		{"🧠", "Internal AI Assistant", "Company-trained AI assistant that knows your SOPs and products.", "$2000", "$597"},
	}

	h.renderTemplate(w, "products", data)
}

func (h *PageHandler) Partners(w http.ResponseWriter, r *http.Request) {
	data := h.baseData(
		"Partner Program — Blueprint Automations",
		"White-label our AI automation suite under your brand. You set the prices. We handle the infrastructure.",
		"partners",
	)
	h.renderTemplate(w, "partners", data)
}

func (h *PageHandler) Pricing(w http.ResponseWriter, r *http.Request) {
	data := h.baseData(
		"Pricing — Blueprint Automations",
		"Simple, transparent pricing for AI automation bundles and individual products. Foundation from $497/mo, Growth from $997/mo, Premium from $2,497/mo.",
		"pricing",
	)
	h.renderTemplate(w, "pricing", data)
}

func (h *PageHandler) CheckoutSuccess(w http.ResponseWriter, r *http.Request) {
	productNames := map[string]string{
		"webchat": "AI Webchat Agent",
		"reviews": "Review Response Automation",
	}

	product := r.URL.Query().Get("product")
	productName, ok := productNames[product]
	if !ok {
		productName = "Blueprint Automations Product"
	}

	data := h.baseData(
		"Welcome! — Blueprint Automations",
		"Your purchase is confirmed.",
		"",
	)
	data["ProductName"] = productName
	h.renderTemplate(w, "checkout_success", data)
}
