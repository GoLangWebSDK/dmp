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
			Config: *cnf,
		},
	}
}

func (mySql *MySQLAdapter) Init() DBAdapter {
	if mySql.DSN == "" {
		mySql.DSN = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
		mySql.DSN = fmt.Sprintf(mySql.DSN,
			mySql.Config.DBUser,
			mySql.Config.DBPass,
			mySql.Config.DBHost,
			mySql.Config.DBPort,
			mySql.Config.DBName,
		)
	}

	return mySql

}

func (mySql *MySQLAdapter) Setup(log zerolog.Logger, logLevel logger.LogLevel) *gorm.DB {
	mySql.DB, mySql.Err = gorm.Open(mysql.Open(mySql.DSN), &gorm.Config{
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

	if mySql.Err != nil {
		log.Fatal().Err(mySql.Err).Msg("Failed to setup database")
	}

	return mySql.DB
}
