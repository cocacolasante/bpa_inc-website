package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cocacolasante/bpa-inc-website/internal/models"
)

func SendDiscordLeadAlert(lead *models.ContactSubmission) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		return
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       fmt.Sprintf("🔥 New Lead: %s", lead.BusinessName),
				"color":       0x6C63FF,
				"description": fmt.Sprintf("**%s** (%s)", lead.FullName, lead.Email),
				"fields": []map[string]interface{}{
					{"name": "Phone", "value": lead.Phone, "inline": true},
					{"name": "Website", "value": lead.Website, "inline": true},
					{"name": "Source", "value": lead.Source, "inline": true},
					{"name": "Message", "value": lead.Message, "inline": false},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	http.Post(webhookURL, "application/json", bytes.NewReader(body))
}

func SendDiscordAlert(message string) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		return
	}

	payload := map[string]interface{}{
		"content": message,
	}

	body, _ := json.Marshal(payload)
	http.Post(webhookURL, "application/json", bytes.NewReader(body))
}
