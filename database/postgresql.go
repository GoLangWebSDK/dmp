package database

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreAdapter struct {
	Adapter
}

func NewPostgreAdapter(cnf *DBConfig) *PostgreAdapter {
	return &PostgreAdapter{
		Adapter: Adapter{
			config: *cnf,
		},
	}
}

func (postgre *PostgreAdapter) Init(log zerolog.Logger, logLevel logger.LogLevel) *gorm.DB {
	postgre.dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		postgre.config.DBHost,
		postgre.config.DBUser,
		postgre.config.DBPass,
		postgre.config.DBName,
		postgre.config.DBPort,
	)

	postgre.db, postgre.err = gorm.Open(
		postgres.Open(postgre.dsn),
		&gorm.Config{
			Logger: logger.New(
				&log, // IO.writer
				logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logLevel,    // Log level, https://gorm.io/docs/logger.html
					IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
					Colorful:                  false,       // Disable color
				},
			),
		},
	)

	if postgre.err != nil {
		log.Fatal().Err(postgre.err).Msg("Failed to setup database")
	}

	return postgre.db
}

func (postgre *PostgreAdapter) Setup(cnf *DBConfig) DBAdapter {
	return NewPostgreAdapter(cnf)
}
