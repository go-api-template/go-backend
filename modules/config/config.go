package config

import (
	"fmt"
	myfeeder "github.com/go-api-template/go-backend/modules/config/feeder"
	"github.com/go-playground/validator/v10"
	configLoader "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/mail"
	"strings"
	"time"
)

// Config is the global config
// It is initialized in the init() function
var Config *AppConfig

// Default values for config
var version = "dev"
var buildDate = "not defined"

// AppConfig holds the configiration of the application
// It is loaded from environment variables and/or app.env file
type AppConfig struct {
	// todo : add more comments
	App struct {
		Name          string `env:"APP_NAME" validate:"required"`
		env           string `env:"APP_ENV" default:"production" validate:"required,oneof=production staging development"`
		Environnement Environnement
		Debug         bool `env:"APP_DEBUG" default:"false"`
		Version       string
		BuildDate     string
	}

	Cors struct {
		Origins []string `env:"CORS_ALLOWED_ORIGINS" default:"*"`
	}

	Tokens struct {
		Access struct {
			PrivateKey string        `env:"ACCESS_TOKEN_PRIVATE_KEY" validate:"required"`
			PublicKey  string        `env:"ACCESS_TOKEN_PUBLIC_KEY" validate:"required"`
			ExpiresIn  time.Duration `env:"ACCESS_TOKEN_EXPIRED_IN" default:"15m"`
			MaxAge     int           `env:"ACCESS_TOKEN_MAX_AGE" default:"15" validate:"required"`
			// todo : voir la différence entre MaxAge et ExpiresIn
		}
		Refresh struct {
			PrivateKey string        `env:"REFRESH_TOKEN_PRIVATE_KEY" validate:"required"`
			PublicKey  string        `env:"REFRESH_TOKEN_PUBLIC_KEY" validate:"required"`
			ExpiresIn  time.Duration `env:"REFRESH_TOKEN_EXPIRED_IN" default:"60m"`
			MaxAge     int           `env:"REFRESH_TOKEN_MAX_AGE" default:"60" validate:"required"`
			// todo : voir la différence entre MaxAge et ExpiresIn
		}
	}

	Logs struct {
		// Configure the access log file, it's size, age, and how many backups to retain.
		Access struct {
			// Filename is the access to write logs to. Backup log files will be retained
			// in the same directory. It uses <processname>-lumberjack.log in
			// os.TempDir() if empty.
			Filename string `env:"ACCESS_LOG_FILE_NAME" default:"access.log" validate:"required"`

			// MaxSize is the maximum size in megabytes of the log access before it gets
			// rotated. It defaults to 100 megabytes.
			MaxSize int `env:"ACCESS_LOG_MAX_SIZE" default:"100"`

			// MaxAge is the maximum number of days to retain old log files based on the
			// timestamp encoded in their filename.  Note that a day is defined as 24
			// hours and may not exactly correspond to calendar days due to daylight
			// savings, leap seconds, etc. The default is not to remove old log files
			// based on age.
			MaxAge int `env:"ACCESS_LOG_MAX_AGE" default:"30"`

			// MaxBackups is the maximum number of old log files to retain.  The default
			// is to retain all old log files (though MaxAge may still cause them to get
			// deleted.)
			MaxBackups int `env:"ACCESS_LOG_MAX_BACKUPS" default:"30"`

			// LocalTime determines if the time used for formatting the timestamps in
			// backup files is the computer's local time.  The default is to use UTC
			// time.
			LocalTime bool `env:"ACCESS_LOG_LOCAL_TIME" default:"true"`

			// Compress determines if the rotated log files should be compressed
			// using gzip. The default is not to perform compression.
			Compress bool `env:"ACCESS_LOG_COMPRESS" default:"false"`
		}

		// Configure the database log file, it's size, age, and how many backups to retain.
		Database struct {
			// Filename is the access to write logs to. Backup log files will be retained
			// in the same directory. It uses <processname>-lumberjack.log in
			// os.TempDir() if empty.
			Filename string `env:"DATABASE_LOG_FILE_NAME" default:"database.log" validate:"required"`

			// MaxSize is the maximum size in megabytes of the log access before it gets
			// rotated. It defaults to 100 megabytes.
			MaxSize int `env:"DATABASE_LOG_MAX_SIZE" default:"100"`

			// MaxAge is the maximum number of days to retain old log files based on the
			// timestamp encoded in their filename.  Note that a day is defined as 24
			// hours and may not exactly correspond to calendar days due to daylight
			// savings, leap seconds, etc. The default is not to remove old log files
			// based on age.
			MaxAge int `env:"DATABASE_LOG_MAX_AGE" default:"30"`

			// MaxBackups is the maximum number of old log files to retain.  The default
			// is to retain all old log files (though MaxAge may still cause them to get
			// deleted.)
			MaxBackups int `env:"DATABASE_LOG_MAX_BACKUPS" default:"30"`

			// LocalTime determines if the time used for formatting the timestamps in
			// backup files is the computer's local time.  The default is to use UTC
			// time.
			LocalTime bool `env:"DATABASE_LOG_LOCAL_TIME" default:"true"`

			// Compress determines if the rotated log files should be compressed
			// using gzip. The default is not to perform compression.
			Compress bool `env:"DATABASE_LOG_COMPRESS" default:"false"`
		}
	}

	Server struct {
		Scheme   string
		Host     string `env:"APP_HOST" validate:"required"`
		Port     string `env:"APP_PORT" default:"8080" validate:"required"`
		BasePath string `env:"APP_BASE_PATH" default:"/api" validate:"required"`
		Url      string
	}

	Client struct {
		Url string `env:"CLIENT_ORIGIN" validate:"required"`
	}

	Database struct {
		Host             string `env:"POSTGRES_HOST" validate:"required"`
		Port             string `env:"POSTGRES_PORT" validate:"required"`
		User             string `env:"POSTGRES_USER" validate:"required"`
		Pass             string `env:"POSTGRES_PASSWORD" validate:"required"`
		Name             string `env:"POSTGRES_NAME" validate:"required"`
		Charset          string `env:"POSTGRES_CHARSET" default:"utf8mb4"`
		ConnectionString string `env:"POSTGRES_CONNECTION_STRING"`
	}

	Redis struct {
		Host string `env:"REDIS_HOST" validate:"required"`
		Port string `env:"REDIS_PORT" validate:"required"`
	}

	Mailer Mailer
}

// Mailer holds the configuration of the mailer
type Mailer struct {
	Enable bool `env:"MAILER_ENABLE" default:"false"`
	From   struct {
		Name    string `env:"MAILER_FROM_NAME"`
		Address string `env:"MAILER_FROM_ADDRESS" validate:"required_if=Enable true,email"`
	}
	EnvelopeFrom         string `env:"MAILER_ENVELOPE_FROM" default:""`
	OverrideEnvelopeFrom bool
	PlainText            bool   `env:"MAILER_SEND_AS_PLAIN_TEXT" default:"false"`
	Prefix               string `env:"MAILER_SUBJECT_PREFIX"`
	Smtp                 struct {
		Protocol             string `env:"SMTP_PROTOCOL" validate:"required_if=Enable true,oneof=smtp smtps smtp+starttls dummy"`
		Host                 string `env:"SMTP_HOST" validate:"required_if=Enable true"`
		Port                 int    `env:"SMTP_PORT" validate:"required_if=Enable true"`
		Username             string `env:"SMTP_USERNAME" validate:"required_if=Enable true"`
		Password             string `env:"SMTP_PASSWORD" validate:"required_if=Enable true"`
		EnableHelo           bool   `env:"SMTP_ENABLE_HELO" default:"true"`
		HeloHostname         string `env:"SMTP_HELO_HOSTNAME"`
		ForceTrustServerCert bool   `env:"SMTP_FORCE_TRUST_SERVER_CERT" default:"false"`
		UseClientCert        bool   `env:"SMTP_USE_CLIENT_CERT" default:"false"`
		ClientCertFile       string `env:"SMTP_CLIENT_CERT_FILE"`
		ClientKeyFile        string `env:"SMTP_CLIENT_KEY_FILE"`
		ClientKeyPassphrase  string `env:"SMTP_CLIENT_KEY_PASSPHRASE"`
	}
}

func init() {
	// Create an instance of AppConfig
	Config = &AppConfig{}

	// Load configuration access
	Config.load()
}

func (c *AppConfig) load() {

	// Create the config loader
	loader := configLoader.New()

	// Add feeder to load from default values
	loader.AddFeeder(myfeeder.Default{})

	// Add feeder to load from *.env file
	loader.AddFeeder(myfeeder.GlobEnvs{Patterns: []string{"app.env", "app.*.env"}})

	// Add feeder to load from environment variables
	loader.AddFeeder(feeder.Env{})

	// Read config access
	err := loader.AddStruct(c).Feed()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load environment variables")
	}
}

// Setup : this function is called while the config is loaded by golobby/config
func (c *AppConfig) Setup() {

	// update app config
	c.setupApp()

	// update server config
	c.setupServer()

	// update token config
	c.setupTokens()

	// update log config
	c.setupLogs()

	// update database config
	c.setupDatabase()

	// update mailer config
	c.setupMailer()

	// Validate config data
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		log.Fatal().Err(err).Msg("loading config access")
	}
}

// setupApp updates the app config
func (c *AppConfig) setupApp() {
	// Version
	c.App.Version = version
	// Build date
	c.App.BuildDate = buildDate

	// Environment
	c.App.Environnement = Environnement(c.App.env)
}

// setupServer updates the server config
func (c *AppConfig) setupServer() {
	// Set the server scheme to http in debug mode
	if c.App.Debug {
		c.Server.Scheme = "http"
	} else {
		c.Server.Scheme = "https"
	}
	// Clean the base path
	c.Server.BasePath = strings.Trim(c.Server.BasePath, " ")
	c.Server.BasePath = strings.Trim(c.Server.BasePath, "/")
	// Update server url
	if len(c.Server.Url) == 0 {
		c.Server.Url = fmt.Sprintf("%s://%s:%s", c.Server.Scheme, c.Server.Host, c.Server.Port)
		if len(c.Server.BasePath) > 0 {
			c.Server.Url = fmt.Sprintf("%s/%s", c.Server.Url, c.Server.BasePath)
		}
	}
}

// setupTokens updates the token config
func (c *AppConfig) setupTokens() {
	// Token expiration
	if c.Tokens.Access.ExpiresIn == 0 {
		c.Tokens.Access.ExpiresIn = time.Duration(c.Tokens.Access.MaxAge) * time.Minute
	}
	if c.Tokens.Refresh.ExpiresIn == 0 {
		c.Tokens.Refresh.ExpiresIn = time.Duration(c.Tokens.Refresh.MaxAge) * time.Minute
	}
}

// setupLogs updates the log config
func (c *AppConfig) setupLogs() {
	// Set the log level to trace in debug mode
	if c.App.Debug {
		log.Logger = log.Level(zerolog.TraceLevel)
	}
}

// setupDatabase updates the database config
func (c *AppConfig) setupDatabase() {
	// Create database connection string
	if len(c.Database.ConnectionString) == 0 {
		c.Database.ConnectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.Pass,
			c.Database.Name)
	}
}

// setupMailer updates the mailer config
func (c *AppConfig) setupMailer() {
	switch c.Mailer.EnvelopeFrom {
	case "":
		c.Mailer.OverrideEnvelopeFrom = false
	case "<>":
		c.Mailer.EnvelopeFrom = ""
		c.Mailer.OverrideEnvelopeFrom = true
	default:
		parsed, err := mail.ParseAddress(c.Mailer.EnvelopeFrom)
		if err != nil {
			log.Fatal().Msgf("Invalid mailer.ENVELOPE_FROM (%s): %v", c.Mailer.EnvelopeFrom, err)
		}
		c.Mailer.EnvelopeFrom = parsed.Address
		c.Mailer.OverrideEnvelopeFrom = true
	}
}
