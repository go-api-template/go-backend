package mailer

import (
	"github.com/go-api-template/go-backend/modules/config"
	"gopkg.in/gomail.v2"
	"jaytaylor.com/html2text"
	"net/mail"
	"time"
)

// Message is the struct that contains the data to send an email
type Message struct {
	from    mail.Address
	To      []*mail.Address
	ReplyTo mail.Address
	Date    time.Time
	Subject string
	Body    string
	Headers map[string][]string
}

// NewMessage creates a new Message
func NewMessage(to, subject, body string) Message {
	// convert the to string to a mail.Address
	toAddress := &mail.Address{Address: to}

	// create a new Message
	message := Message{
		To:      []*mail.Address{toAddress},
		Subject: subject,
		Date:    time.Now(),
		Body:    body,
	}

	// get the from address from the config
	message.from.Name = config.Config.Mailer.From.Name
	message.from.Address = config.Config.Mailer.From.Address

	return message
}

// From sets the from address
func (m *Message) From(name string, address string) {
	m.from.Name = name
	m.from.Address = address
}

// SetHeader adds additional headers to a message
func (m *Message) SetHeader(field string, value ...string) {
	m.Headers[field] = value
}

// ToMessage converts a Message to gomail.Message
func ToMessage(message *Message) *gomail.Message {
	// create a new gomail.Message
	m := gomail.NewMessage()
	// set the sender
	m.SetAddressHeader("From", message.from.Address, message.from.Name)
	// convert the list of recipients to a list of strings
	// because gomail.Message.SetHeader() accepts a list of strings
	tos := make([]string, len(message.To))
	for i, v := range message.To {
		tos[i] = m.FormatAddress(v.Address, v.Name)
	}
	// set the recipients
	m.SetHeader("To", tos...)
	// set the reply to
	m.SetAddressHeader("Reply-To", message.ReplyTo.Address, message.ReplyTo.Name)
	// set the subject
	m.SetHeader("Subject", message.Subject)
	// set the date
	m.SetDateHeader("Date", message.Date)
	// prevent Outlook from adding a "Reply" header
	m.SetHeader("X-Auto-Response-Suppress", "All")
	// set the body
	plainBody, err := html2text.FromString(message.Body)
	if err != nil {
		m.SetBody("text/plain", plainBody)
	} else {
		m.SetBody("text/plain", plainBody)
		m.AddAlternative("text/html", message.Body)
	}
	// set additional headers
	for header := range message.Headers {
		m.SetHeader(header, message.Headers[header]...)
	}

	return m
}
