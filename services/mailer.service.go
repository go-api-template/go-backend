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
	SendVerificationCode(user *models.User) error
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

// SendVerificationCode sends a verification code to the user
func (s *MailerServiceImpl) SendVerificationCode(user *models.User) error {
	// data to be passed to the template
	data := map[string]any{
		"Title":     fmt.Sprintf("Your verification code %s", config.Config.App.Name),
		"AppUrl":    config.Config.Server.Url,
		"AppName":   config.Config.App.Name,
		"VerifyUrl": config.Config.Client.Url + "/auth/verify?key=" + user.VerificationCode,
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

//// SendVerificationCode sends a verification code to the user
//func (s *MailService) SendVerificationCode(user *models.User, code string) error {
//
//	// emailData is the data that will be passed to the template
//	type emailData struct {
//		Subject   string
//		FirstName string
//		URL       string
//	}
//
//	//
//	data := emailData{
//		Subject: "Votre code de vérification GoVoit",
//		//FirstName: user.FirstName(),
//		URL: config.Config.Server.FrontEndUrl + "/verify?key=" + code,
//	}
//
//	return s.sendTemplateEmail(user.Email, data.Subject, "auth_verify.gohtml", data)
//}
//
//// sendTemplateEmail sends an email where the body is a template.
//func (s *MailService) sendTemplateEmail(to string, subject string, templateName string, data any) error {
//
//	body, err := s.generateEmailBody(templateName, data)
//	if err != nil {
//		return err
//	}
//
//	return s.sendEmail(to, subject, body.String())
//}
//

//// sendEmail sends an email.
//// It can be used to send emails with SMTP or Mailer.
//func (s *MailService) sendEmail(recipient string, subject string, body string) error {
//	if config.Config.Mail.Type == "smtp" {
//		return s.sendSmtpEmail(recipient, subject, body)
//	} else {
//		return s.sendMailerEmail(recipient, subject, body)
//	}
//}
//
//// sendSmtpEmail sends an email with SMTP.
//func (s *MailService) sendSmtpEmail(recipient string, subject string, body string) error {
//
//	smtpHost := config.Config.Mail.Smtp.Host
//	smtpPort := config.Config.Mail.Smtp.Port
//	smtpUser := config.Config.Mail.Smtp.User
//	smtpPass := config.Config.Mail.Smtp.Pass
//	from := config.Config.Mail.Smtp.Sender
//	to := recipient
//
//	m := gomail.NewMessage()
//
//	m.SetHeader("From", from)
//	m.SetHeader("To", to)
//	m.SetHeader("Subject", subject)
//	m.SetBody("text/html", body)
//	m.AddAlternative("text/plain", html2text.HTML2Text(body))
//
//	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
//	// #nosec G402
//	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
//
//	// Send Email
//	if err := d.DialAndSend(m); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (s *MailService) sendMailerEmail(recipient string, subject string, body string) error {
//
//	mailer := mailerClient.Mailer{Url: config.Config.Mail.Mailer.Url}
//	mailerUser := mailerClient.User{Username: config.Config.Mail.Mailer.User, Password: config.Config.Mail.Mailer.Pass}
//
//	message := mailerClient.Message{
//		From:    config.Config.Mail.Mailer.Sender,
//		To:      recipient,
//		Subject: subject,
//	}
//
//	message.HtmlBody = body
//	message.PlainBody = html2text.HTML2Text(body)
//
//	// Send Email
//	err := mailer.SendSecureMail(mailerUser, message, nil)
//	if err != nil {
//		return err
//	}
//	return nil
//}
