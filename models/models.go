package models

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"sync"
)

type models struct {
	gormDb *gorm.DB
}

var (
	m   *models
	one sync.Once
)

func NewModels(g *gorm.DB) {
	one.Do(func() {
		m = &models{gormDb: g}
		m.registerSerializers()
		m.autoMigrate()
		m.seed()
	})
}

// registerSerializers must be called used to register serializers
func (m *models) registerSerializers() {
}

// autoMigrate migrates models
// This function must be used in order to migrate models
func (m *models) autoMigrate() {
	// Initialize models
	err := m.gormDb.Debug().Transaction(func(tx *gorm.DB) error {
		return tx.AutoMigrate()
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Migrating models")
	}
}

// seed function must be used to seed data into database
func (m *models) seed() {
	// Initialize models
	err := m.gormDb.Debug().Transaction(func(tx *gorm.DB) error {
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Seeding models")
	}
}
