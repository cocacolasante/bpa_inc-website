package models

import (
	"errors"
	"regexp"
	"strings"
)

type ContactSubmission struct {
	FullName     string `json:"full_name"`
	BusinessName string `json:"business_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Website      string `json:"website"`
	Message      string `json:"message"`
	Source       string `json:"source"` // "contact_form" | "roi_audit" | "self_serve_checkout"
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var phoneRegex = regexp.MustCompile(`^[\d\s\-\+\(\)]{10,}$`)

func (c *ContactSubmission) Validate() error {
	if c.Source == "" {
		c.Source = "contact_form"
	}

	c.FullName = strings.TrimSpace(c.FullName)
	c.BusinessName = strings.TrimSpace(c.BusinessName)
	c.Phone = strings.TrimSpace(c.Phone)
	c.Email = strings.TrimSpace(c.Email)
	c.Website = strings.TrimSpace(c.Website)
	c.Message = strings.TrimSpace(c.Message)

	if c.FullName == "" {
		return errors.New("full name is required")
	}
	if len(c.FullName) < 2 {
		return errors.New("full name must be at least 2 characters")
	}

	if c.BusinessName == "" {
		return errors.New("business name is required")
	}

	if c.Phone == "" {
		return errors.New("phone number is required")
	}
	if !phoneRegex.MatchString(c.Phone) {
		return errors.New("invalid phone number format")
	}

	if c.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(c.Email) {
		return errors.New("invalid email format")
	}

	if c.Message == "" {
		return errors.New("message is required")
	}
	if len(c.Message) < 10 {
		return errors.New("message must be at least 10 characters")
	}

	return nil
}
