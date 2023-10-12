package postgres_db

import (
	"context"
	"database/sql"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/config"
	database_logger "github.com/go-api-template/go-backend/modules/logger/database"
	zerologGorm "github.com/go-mods/zerolog-gorm"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

type postgresDb struct {
	gormDb *gorm.DB
	sqlDb  *sql.DB
}

var (
	// p is the singleton instance of the postgresDb struct
	p *postgresDb
	// Prevent multiple initialization
	once sync.Once
)

func NewPostgres(ctx context.Context) (*gorm.DB, *sql.DB) {
	once.Do(func() {
		p = &postgresDb{}
		p.initialize(ctx)
		p.initializeModels()
	})
	return p.gormDb, p.sqlDb
}

func (p *postgresDb) initialize(ctx context.Context) {
	log.Debug().Msg("Initializing Postgres...")

	// Initialize gormDb
	gormDb, err := gorm.Open(postgres.Open(config.Config.Database.ConnectionString), &gorm.Config{
		Logger: &zerologGorm.GormLogger{},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Initializing Postgres")
	}

	// Add logger to gorm
	gormDb = gormDb.WithContext(database_logger.DatabaseLogger.WithContext(ctx))

	// Get sqlDb from gorm
	sqlDb, err := gormDb.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Getting Postgres instance")
	}

	// Validate connection to database by pinging it
	err = sqlDb.PingContext(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Pinging Postgres instance")
	}

	// Set gorm in debug mode if debug mode is enabled
	if config.Config.App.Debug {
		gormDb = gormDb.Debug()
	}

	// Set the singleton instance
	p.gormDb = gormDb
	p.sqlDb = sqlDb
}

func (p *postgresDb) initializeModels() {
	models.NewModels(p.gormDb)
}
