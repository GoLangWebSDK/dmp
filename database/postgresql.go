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
	*Adapter
}

func NewPostgreAdapter(cnf *DBConfig) *PostgreAdapter {
	return &PostgreAdapter{
		Adapter: &Adapter{
			Config: *cnf,
		},
	}
}

func (postgre *PostgreAdapter) Init() DBAdapter {
	if postgre.DSN == "" {
		postgre.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			postgre.Config.DBHost,
			postgre.Config.DBUser,
			postgre.Config.DBPass,
			postgre.Config.DBName,
			postgre.Config.DBPort,
		)
	}

	return postgre
}

func (postgre *PostgreAdapter) Setup(log zerolog.Logger, logLevel logger.LogLevel) *gorm.DB {
	postgre.DB, postgre.Err = gorm.Open(
		postgres.Open(postgre.DSN),
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

	if postgre.Err != nil {
		log.Fatal().Err(postgre.Err).Msg("Failed to setup database")
	}

	return postgre.DB
}
