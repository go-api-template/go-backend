package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	configLoader "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var Config *AppConfig
var version = "dev"
var buildDate = "not defined"

type AppConfig struct {
	App struct {
		Name      string `env:"APP_NAME" validate:"required"`
		Env       string `env:"APP_ENV" validate:"required,oneof=development production"`
		Debug     bool   `env:"APP_DEBUG"`
		Port      int    `env:"APP_PORT" validate:"required"`
		Version   string
		BuildDate string
	}

	Cors struct {
		Origins []string `env:"CORS_ALLOWED_ORIGINS"`
	}

	Tokens struct {
		Access struct {
			PrivateKey string        `env:"ACCESS_TOKEN_PRIVATE_KEY" validate:"required"`
			PublicKey  string        `env:"ACCESS_TOKEN_PUBLIC_KEY" validate:"required"`
			ExpiresIn  time.Duration `env:"ACCESS_TOKEN_EXPIRED_IN"`
			MaxAge     int           `env:"ACCESS_TOKEN_MAX_AGE" validate:"required"`
		}
		Refresh struct {
			PrivateKey string        `env:"REFRESH_TOKEN_PRIVATE_KEY" validate:"required"`
			PublicKey  string        `env:"REFRESH_TOKEN_PUBLIC_KEY" validate:"required"`
			ExpiresIn  time.Duration `env:"REFRESH_TOKEN_EXPIRED_IN"`
			MaxAge     int           `env:"REFRESH_TOKEN_MAX_AGE" validate:"required"`
		}
	}

	Logs struct {
		Access struct {
			// Filename is the access to write logs to. Backup log files will be retained
			// in the same directory. It uses <processname>-lumberjack.log in
			// os.TempDir() if empty.
			Filename string `env:"ACCESS_LOG_FILE_NAME" validate:"required"`

			// MaxSize is the maximum size in megabytes of the log access before it gets
			// rotated. It defaults to 100 megabytes.
			MaxSize int `env:"ACCESS_LOG_MAX_SIZE"`

			// MaxAge is the maximum number of days to retain old log files based on the
			// timestamp encoded in their filename.  Note that a day is defined as 24
			// hours and may not exactly correspond to calendar days due to daylight
			// savings, leap seconds, etc. The default is not to remove old log files
			// based on age.
			MaxAge int `env:"ACCESS_LOG_MAX_AGE"`

			// MaxBackups is the maximum number of old log files to retain.  The default
			// is to retain all old log files (though MaxAge may still cause them to get
			// deleted.)
			MaxBackups int `env:"ACCESS_LOG_MAX_BACKUPS"`

			// LocalTime determines if the time used for formatting the timestamps in
			// backup files is the computer's local time.  The default is to use UTC
			// time.
			LocalTime bool `env:"ACCESS_LOG_LOCAL_TIME"`

			// Compress determines if the rotated log files should be compressed
			// using gzip. The default is not to perform compression.
			Compress bool `env:"ACCESS_LOG_COMPRESS"`
		}

		Database struct {
			// Filename is the access to write logs to. Backup log files will be retained
			// in the same directory. It uses <processname>-lumberjack.log in
			// os.TempDir() if empty.
			Filename string `env:"DATABASE_LOG_FILE_NAME" validate:"required"`

			// MaxSize is the maximum size in megabytes of the log access before it gets
			// rotated. It defaults to 100 megabytes.
			MaxSize int `env:"DATABASE_LOG_MAX_SIZE"`

			// MaxAge is the maximum number of days to retain old log files based on the
			// timestamp encoded in their filename.  Note that a day is defined as 24
			// hours and may not exactly correspond to calendar days due to daylight
			// savings, leap seconds, etc. The default is not to remove old log files
			// based on age.
			MaxAge int `env:"DATABASE_LOG_MAX_AGE"`

			// MaxBackups is the maximum number of old log files to retain.  The default
			// is to retain all old log files (though MaxAge may still cause them to get
			// deleted.)
			MaxBackups int `env:"DATABASE_LOG_MAX_BACKUPS"`

			// LocalTime determines if the time used for formatting the timestamps in
			// backup files is the computer's local time.  The default is to use UTC
			// time.
			LocalTime bool `env:"DATABASE_LOG_LOCAL_TIME"`

			// Compress determines if the rotated log files should be compressed
			// using gzip. The default is not to perform compression.
			Compress bool `env:"DATABASE_LOG_COMPRESS"`
		}
	}

	Server struct {
		Scheme      string
		Host        string `env:"APP_HOST" validate:"required"`
		Port        string `env:"APP_PORT" validate:"required"`
		Url         string
		BasePath    string `env:"APP_BASE_PATH"`
		FrontEndUrl string `env:"APP_FRONT_END_URL"`
	}

	Client struct {
		Url string `env:"CLIENT_ORIGIN" validate:"required"`
	}

	Database struct {
		Host    string `env:"POSTGRES_HOST" validate:"required"`
		Port    string `env:"POSTGRES_PORT" validate:"required"`
		User    string `env:"POSTGRES_USER" validate:"required"`
		Pass    string `env:"POSTGRES_PASSWORD" validate:"required"`
		Name    string `env:"POSTGRES_NAME" validate:"required"`
		Charset string `env:"POSTGRES_CHARSET"`

		ConnectionString string
	}

	Redis struct {
		Host string `env:"REDIS_HOST" validate:"required"`
		Port string `env:"REDIS_PORT" validate:"required"`
	}

	Mail struct {
		Type string `env:"MAIL_DRIVER" validate:"required,oneof=smtp mailer"`
		Smtp struct {
			Sender string `env:"SMTP_FROM" validate:"required"`
			Host   string `env:"SMTP_HOST" validate:"required"`
			Port   int    `env:"SMTP_PORT" validate:"required"`
			User   string `env:"SMTP_USER" validate:"required"`
			Pass   string `env:"SMTP_PASS" validate:"required"`
		}
		Mailer struct {
			Sender string `env:"MAILER_FROM" validate:"required"`
			Url    string `env:"MAILER_URL" validate:"required"`
			User   string `env:"MAILER_USER" validate:"required"`
			Pass   string `env:"MAILER_PASS" validate:"required"`
		}
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

	// Add feeder to load from app.env file
	if _, err := os.Stat("app.env"); err == nil {
		loader.AddFeeder(feeder.DotEnv{Path: "app.env"})
	}

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

	// Validate config data
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		log.Fatal().Err(err).Msg("loading config access")
	}

	// Debug mode ?
	debug := c.App.Debug

	// Set the server scheme to http in debug mode
	if debug {
		c.Server.Scheme = "http"
	} else {
		c.Server.Scheme = "https"
	}

	// Update server url
	if len(c.Server.Url) == 0 {
		c.Server.Url = fmt.Sprintf("%s://%s:%s", c.Server.Scheme, c.Server.Host, c.Server.Port)
	}

	// Version
	c.App.Version = version

	// Build date
	c.App.BuildDate = buildDate

	// Set the log level to trace in debug mode
	if debug {
		log.Logger = log.Level(zerolog.TraceLevel)
	}

	// Token expiration
	c.Tokens.Access.ExpiresIn = time.Duration(c.Tokens.Access.MaxAge) * time.Minute
	c.Tokens.Refresh.ExpiresIn = time.Duration(c.Tokens.Refresh.MaxAge) * time.Minute

	// Create database connection string
	c.Database.ConnectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Pass,
		c.Database.Name)
}
