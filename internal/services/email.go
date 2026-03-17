package services

import (
	"fmt"
	"os"

	"github.com/cocacolasante/bpa-inc-website/internal/models"
	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client   *resend.Client
	from     string
	notifyTo string
}

func NewEmailService() *EmailService {
	return &EmailService{
		client:   resend.NewClient(os.Getenv("RESEND_API_KEY")),
		from:     os.Getenv("FROM_EMAIL"),
		notifyTo: os.Getenv("NOTIFY_EMAIL"),
	}
}

func (s *EmailService) SendLeadNotification(lead *models.ContactSubmission) error {
	if os.Getenv("RESEND_API_KEY") == "" {
		return nil
	}

	body := fmt.Sprintf(`
		<h2>New Lead: %s at %s</h2>
		<p><strong>Email:</strong> %s</p>
		<p><strong>Phone:</strong> %s</p>
		<p><strong>Website:</strong> %s</p>
		<p><strong>Source:</strong> %s</p>
		<hr>
		<p><strong>Message:</strong><br>%s</p>
	`, lead.FullName, lead.BusinessName, lead.Email, lead.Phone, lead.Website, lead.Source, lead.Message)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{s.notifyTo},
		Subject: fmt.Sprintf("New Lead: %s — %s", lead.FullName, lead.BusinessName),
		Html:    body,
	})
	return err
}

func (s *EmailService) SendAutoReply(lead *models.ContactSubmission) error {
	if os.Getenv("RESEND_API_KEY") == "" {
		return nil
	}

	body := fmt.Sprintf(`
		<h2>Thanks for reaching out, %s!</h2>
		<p>We received your message and will get back to you within 24 hours.</p>
		<p>In the meantime, you can book a free strategy call directly on our calendar:</p>
		<p><a href="%s" style="background:#6C63FF;color:white;padding:12px 24px;
		   border-radius:4px;text-decoration:none;display:inline-block;">
		   Book a Free Strategy Call</a></p>
		<p>— Blueprint Automations Team</p>
	`, lead.FullName, os.Getenv("CAL_COM_URL"))

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{lead.Email},
		Subject: "We received your message — Blueprint Automations",
		Html:    body,
	})
	return err
}

func (s *EmailService) SendPartnerApplicationAlert(app *models.PartnerApplication) error {
	if os.Getenv("RESEND_API_KEY") == "" {
		return nil
	}

	body := fmt.Sprintf(`
		<h2>New Partner Application</h2>
		<p><strong>Company:</strong> %s</p>
		<p><strong>Contact:</strong> %s (%s)</p>
		<p><strong>Phone:</strong> %s</p>
		<p><strong>Website:</strong> %s</p>
		<p><strong>Years in Business:</strong> %s</p>
		<p><strong>Active Clients:</strong> %s</p>
		<p><strong>Expected Volume:</strong> %s/month</p>
		<hr>
		<p><strong>Why Partner:</strong><br>%s</p>
		<hr>
		<p>Review and approve at: <a href="%s/portals/super/partners/applications">Blueprint Command</a></p>
	`, app.CompanyName, app.ContactName, app.ContactEmail,
		app.ContactPhone, app.Website, app.YearsInBusiness,
		app.ClientCount, app.ExpectedVolume, app.WhyPartner,
		os.Getenv("BLUEPRINT_COMMAND_API_URL"),
	)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{os.Getenv("PARTNER_NOTIFY_EMAIL")},
		Subject: fmt.Sprintf("Partner Application: %s", app.CompanyName),
		Html:    body,
	})
	return err
}

func (s *EmailService) SendPartnerApplicationConfirmation(app *models.PartnerApplication) error {
	if os.Getenv("RESEND_API_KEY") == "" {
		return nil
	}

	body := fmt.Sprintf(`
		<h2>Thanks for applying, %s!</h2>
		<p>We received your partner application for <strong>%s</strong> and will
		review it within 2 business days.</p>
		<p>Once approved, you'll receive access to your partner portal where you can:</p>
		<ul>
			<li>Provision clients with your white-labeled branding</li>
			<li>Manage all your clients in one place</li>
			<li>Track usage, billing, and earnings</li>
		</ul>
		<p>Questions? Just reply to this email.</p>
		<p>— Blueprint Automations Team</p>
	`, app.ContactName, app.CompanyName)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{app.ContactEmail},
		Subject: "Partner Application Received — Blueprint Automations",
		Html:    body,
	})
	return err
}
