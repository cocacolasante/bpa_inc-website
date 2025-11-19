package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type PageHandler struct {
	logger    *log.Logger
	templates map[string]*template.Template
}

func NewPageHandler(logger *log.Logger) *PageHandler {
	templates := make(map[string]*template.Template)

	pages := []string{"home", "about", "services", "portfolio", "contact"}

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

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "Blueprint Automations - AI Workflow Automation for Growing Businesses",
		"Description": "We build clean, automated AI systems that let your business run with clarity and ease. Specializing in N8N, AI automations, and API integrations.",
		"Page":        "home",
	}
	h.renderTemplate(w, "home", data)
}

func (h *PageHandler) About(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "About Us - Blueprint Automations",
		"Description": "Meet the team behind Blueprint Automations. Founded in 2025 by Scott Henry and Anthony Colasante, we help small to medium businesses leverage AI automation.",
		"Page":        "about",
	}
	h.renderTemplate(w, "about", data)
}

func (h *PageHandler) Services(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "Our Services - AI Automation Solutions",
		"Description": "Professional AI workflow automation services. Custom workflows for $2,500 with ongoing support starting at $495/month.",
		"Page":        "services",
	}
	h.renderTemplate(w, "services", data)
}

func (h *PageHandler) Portfolio(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "Portfolio - Our Work",
		"Description": "Explore our AI automation workflows and projects. See how we've helped businesses automate and scale.",
		"Page":        "portfolio",
	}
	h.renderTemplate(w, "portfolio", data)
}

func (h *PageHandler) Contact(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "Contact Us - Let's Automate Your Business",
		"Description": "Ready to automate your business? Get in touch with Blueprint Automations for a consultation.",
		"Page":        "contact",
	}
	h.renderTemplate(w, "contact", data)
}
