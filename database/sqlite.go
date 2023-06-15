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
			Config: *cnf,
		},
	}
}

// Setup initializes the database instance
func (sql *SQLiteAdapter) Init(log zerolog.Logger, logLevel logger.LogLevel) *gorm.DB {
	sql.DB, sql.Err = gorm.Open(sqlite.Open(sql.Config.DBName), &gorm.Config{})

	if sql.Err != nil {
		log.Fatal().Err(sql.Err).Msg("Failed to setup database")
	}

	return sql.DB
}

func (sql *SQLiteAdapter) Setup(cnf *DBConfig) DBAdapter {
	return NewSQLiteAdapter(cnf)
}
