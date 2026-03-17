package models

import "errors"

type PartnerApplication struct {
	CompanyName     string `json:"company_name"`
	ContactName     string `json:"contact_name"`
	ContactEmail    string `json:"contact_email"`
	ContactPhone    string `json:"contact_phone"`
	Website         string `json:"website"`
	YearsInBusiness string `json:"years_in_business"`
	ClientCount     string `json:"client_count"`
	ExpectedVolume  string `json:"expected_volume"`
	WhyPartner      string `json:"why_partner"`
}

func (p *PartnerApplication) Validate() error {
	if p.CompanyName == "" {
		return errors.New("company name is required")
	}
	if p.ContactName == "" {
		return errors.New("contact name is required")
	}
	if p.ContactEmail == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(p.ContactEmail) {
		return errors.New("invalid email")
	}
	if p.ContactPhone == "" {
		return errors.New("phone is required")
	}
	if len(p.WhyPartner) < 20 {
		return errors.New("please tell us more about why you want to partner")
	}
	return nil
}
