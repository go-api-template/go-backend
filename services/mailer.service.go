package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/mailer"
	"html/template"
	"os"
	"path/filepath"
)

// MailerService is an interface for the MailerServiceImpl
// It declares the methods that the service must implement
type MailerService interface {
	SendVerificationToken(user *models.User) error
	SendResetToken(user *models.User) error
}

// MailerServiceImpl is the service for mail
// It implements the IMailService interface
type MailerServiceImpl struct {
	ctx        context.Context
	mailer     *mailer.Mailer
	config     config.Mailer
	layoutTmpl *template.Template
}

type TplName string

const (
	mailAuthVerify TplName = "auth/verify"
	mailAuthReset  TplName = "auth/reset"
)

const (
	mailTemplatesPath = "templates/mail"
)

// MailService implements the IMailService interface
var _ MailerService = &MailerServiceImpl{}

// NewMailerService creates a new MailerServiceImpl
func NewMailerService(ctx context.Context) (MailerService, error) {
	// Check mailer setting.
	if !config.Config.Mailer.Enable {
		return nil, nil
	}

	// Create a new mailer service
	s := &MailerServiceImpl{
		ctx:    ctx,
		config: config.Config.Mailer,
		mailer: mailer.NewMailer(ctx),
	}

	// Load templates from the file system
	err := s.loadTemplates()
	if err != nil {
		return nil, err
	}

	// start the mailer service
	s.mailer.Start()

	return s, nil
}

// loadTemplates loads the templates.
func (s *MailerServiceImpl) loadTemplates() error {
	if s.layoutTmpl != nil {
		return nil
	}

	// Create the main template which is the base template for all other templates
	layoutTmpl, err := s.loadLayoutTemplate()
	if err != nil {
		return err
	}

	// Assign the layout template
	s.layoutTmpl = layoutTmpl

	return nil
}

// loadLayoutTemplate loads the layout template which is the base template for all other templates
func (s *MailerServiceImpl) loadLayoutTemplate() (*template.Template, error) {
	// Create the layout template
	layoutTmpl, err := template.ParseGlob(filepath.Clean(mailTemplatesPath + "/layout/*.gohtml"))
	if err != nil {
		return nil, err
	}

	return layoutTmpl, nil
}

// generateFromTemplate generates a body from a template name and data
func (s *MailerServiceImpl) generateFromTemplate(tplName TplName, data any) (body string, err error) {
	// Check if the layout template has been loaded
	if s.layoutTmpl == nil {
		return "", errors.New("layout template not loaded")
	}

	// Check if template file exists
	tplFile := filepath.Clean(mailTemplatesPath + "/" + string(tplName) + ".gohtml")
	if _, err := os.Stat(tplFile); os.IsNotExist(err) {
		return "", err
	}

	// Create the template
	tpl, err := template.Must(s.layoutTmpl.Clone()).ParseFiles(tplFile)
	if err != nil {
		return "", err
	}

	// Execute template
	var content bytes.Buffer
	err = tpl.ExecuteTemplate(&content, "layout/main", data)

	return content.String(), err
}

// SendVerificationToken sends a verification token to the user
func (s *MailerServiceImpl) SendVerificationToken(user *models.User) error {
	if user.VerificationToken == "" {
		return errors.New("verification token is empty")
	}

	// data to be passed to the template
	data := map[string]any{
		"Title":     fmt.Sprintf("Your verification code %s", config.Config.App.Name),
		"AppUrl":    config.Config.Server.Url,
		"AppName":   config.Config.App.Name,
		"VerifyUrl": config.Config.Client.Url + "/auth/verify?key=" + user.VerificationToken,
	}

	// Generate the email body from the template
	content, err := s.generateFromTemplate(mailAuthVerify, data)
	if err != nil {
		return err
	}

	// Prepare the email
	message := mailer.NewMessage(user.Email, fmt.Sprintf("Your verification code %s", config.Config.App.Name), content)

	// Send the email
	s.mailer.SendMail(&message)

	return nil
}

// SendResetToken sends a reset token to the user
func (s *MailerServiceImpl) SendResetToken(user *models.User) error {
	if user.ResetToken == "" {
		return errors.New("reset token is empty")
	}

	// data to be passed to the template
	data := map[string]any{
		"Title":    fmt.Sprintf("Reset your %s password", config.Config.App.Name),
		"AppUrl":   config.Config.Server.Url,
		"AppName":  config.Config.App.Name,
		"ResetUrl": config.Config.Client.Url + "/auth/reset-password?key=" + user.ResetToken,
	}

	// Generate the email body from the template
	content, err := s.generateFromTemplate(mailAuthReset, data)
	if err != nil {
		return err
	}

	// Prepare the email
	message := mailer.NewMessage(user.Email, fmt.Sprintf("Reset your %s password", config.Config.App.Name), content)

	// Send the email
	s.mailer.SendMail(&message)

	return nil
}
