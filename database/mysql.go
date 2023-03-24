package database

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLAdapter struct {
	Adapter
}

func NewMySQLAdapter(cnf *DBConfig) *MySQLAdapter {
	return &MySQLAdapter{
		Adapter: Adapter{
			config: *cnf,
		},
	}
}

func (mySql *MySQLAdapter) Init(log zerolog.Logger, logLevel logger.LogLevel) *gorm.DB {
	mySql.dsn = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	mySql.dsn = fmt.Sprintf(mySql.dsn,
		mySql.config.DBUser,
		mySql.config.DBPass,
		mySql.config.DBHost,
		mySql.config.DBPort,
		mySql.config.DBName,
	)

	mySql.db, mySql.err = gorm.Open(mysql.Open(mySql.dsn), &gorm.Config{
		Logger: logger.New(
			&log, // IO.writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logLevel,    // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		),
	})

	if mySql.err != nil {
		log.Fatal().Err(mySql.err).Msg("Failed to setup database")
	}

	return mySql.db
}

func (mySql *MySQLAdapter) Setup(cnf *DBConfig) DBAdapter {
	return NewMySQLAdapter(cnf)
}
