package database

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type DBMigrator interface {
	Models() []interface{}
	Migrations() []*gormigrate.Migration
}

type Migration struct {
	DBAManager *Database
	Migrations []*gormigrate.Migration
	Models     []interface{}
	db         *gorm.DB
	log        zerolog.Logger
}

func NewMigration(dbm *Database) *Migration {
	return &Migration{
		DBAManager: dbm,
		db:         dbm.Engine,
	}
}

func (m *Migration) Run() error {
	if err := m.migrateModels(); err != nil {
		m.log.Fatal().Err(err).Msg("Failed to migrate models")
		return err
	}

	if m.Migrations == nil {
		return errors.New("No migrations to run!")
	}

	migration := gormigrate.New(m.db, gormigrate.DefaultOptions, m.Migrations)
	migration.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(m.Models...)
		if err != nil {
			m.log.Fatal().Err(err).Msg("Init Schema failed")
			return err
		}
		return nil
	})

	if err := migration.Migrate(); err != nil && err != gorm.ErrInvalidField {
		m.log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	return nil
}

func (m *Migration) migrateModels() error {
	stmt := &gorm.Statement{DB: m.db}

	if len(m.Models) == 0 {
		return errors.New("No models to migrate!")
	}

	for _, model := range m.Models {
		if err := stmt.Parse(&model); err != nil {
			return err
		}

		id := fmt.Sprintf("create_%v", stmt.Schema.Table)
		migration := &gormigrate.Migration{
			ID: id,
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&model)
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(stmt.Schema.Table)
			},
		}
		m.Migrations = append(m.Migrations, migration)
	}

	return nil
}

func (m *Migration) AddModels(models ...interface{}) *Migration {
	m.Models = append(m.Models, models...)
	return m
}

func (m *Migration) AddMigrations(migrations ...*gormigrate.Migration) *Migration {
	m.Migrations = append(m.Migrations, migrations...)
	return m
}

func (m *Migration) AddMigrator(migrator DBMigrator) *Migration {
	m.Migrations = append(m.Migrations, migrator.Migrations()...)
	m.Models = append(m.Models, migrator.Models()...)
	return m
}
