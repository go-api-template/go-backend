package mailer

import (
	"context"
	"crypto/tls"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
	"sync"
)

// Mailer is used to send emails
// It use a buffered channel to queue mails and a context to stop the service
var mailer *Mailer
var once sync.Once

// Mailer is used to send emails
type Mailer struct {
	ctx      context.Context
	finished context.CancelFunc
	mutex    sync.Mutex

	mailChannel chan *Message
	isStarted   bool
	isRunning   bool
}

// NewMailer creates a new Mailer
func NewMailer(ctx context.Context) *Mailer {
	once.Do(func() {
		// create a cancelable context from the original context
		ctx, cancel := context.WithCancel(ctx)

		// create the mailer
		mailer = &Mailer{
			ctx:         ctx,
			finished:    cancel,
			mailChannel: make(chan *Message, 100),
		}
	})
	return mailer
}

// Start starts the mailer
func (m *Mailer) Start() {
	// prevent multiple start
	if m.isStarted {
		return
	}
	// start the goroutine which listens to the mail channel
	go func() {
		for {
			select {
			case message := <-m.mailChannel:
				m.sendMail(message)
				if len(m.mailChannel) == 0 {
					m.isRunning = false
				}
			}
		}
	}()
	// set isStarted to true
	m.isStarted = true
}

// Stop stops the mail service
func (m *Mailer) Stop() {
	m.finished()
}

// Count returns the number of messages in the mail channel
func (m *Mailer) Count() int {
	return len(m.mailChannel)
}

// IsRunning returns true if the mailer is running
func (m *Mailer) IsRunning() bool {
	return m.isRunning
}

// SendMail adds a message to the mail channel
func (m *Mailer) SendMail(message *Message) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.isRunning = true
	m.mailChannel <- message
}

// sendMail sends the mail
func (m *Mailer) sendMail(message *Message) {
	// create a new gomail.Message
	mailMessage := ToMessage(message)

	// get the mailer config
	config := config.Config.Mailer.Smtp

	// create a new dialer
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	// #nosec G402
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// send the email
	if err := dialer.DialAndSend(mailMessage); err != nil {
		log.Error().Err(err).Msg("Failed to send verification code")
	}
}
