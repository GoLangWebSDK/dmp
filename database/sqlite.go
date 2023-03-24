package database

import (
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLiteAdapter struct {
	Adapter
}

func NewSQLiteAdapter(cnf *DBConfig) *SQLiteAdapter {
	return &SQLiteAdapter{
		Adapter: Adapter{
			config: *cnf,
		},
	}
}

// Setup initializes the database instance
func (sql *SQLiteAdapter) Init(log zerolog.Logger, logLevel logger.LogLevel) *gorm.DB {
	sql.db, sql.err = gorm.Open(sqlite.Open(sql.config.DBName), &gorm.Config{})

	if sql.err != nil {
		log.Fatal().Err(sql.err).Msg("Failed to setup database")
	}

	return sql.db
}

func (sql *SQLiteAdapter) Setup(cnf *DBConfig) DBAdapter {
	return NewSQLiteAdapter(cnf)
}
