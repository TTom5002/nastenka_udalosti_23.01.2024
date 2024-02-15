package models

import (
	"nastenka_udalosti/internal/forms"
)

// TemplateData drží data k odeslání z handlers do templates

type TemplateData struct {
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
	User            User
}
