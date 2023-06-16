package database

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global variable for database connection
var DBManager *Database

type DBAdapter interface {
	Init() DBAdapter
	Setup(log zerolog.Logger, logLevel logger.LogLevel) *gorm.DB
	// NewMock() (sqlmock.Sqlmock, *gorm.DB)
}

type Adapter struct {
	DB     *gorm.DB
	Config DBConfig
	DSN    string
	Err    error
}

type Database struct {
	Engine    *gorm.DB
	DBadapter DBAdapter
	DBconfig  DBConfig
	LogLvl    logger.LogLevel
	Log       zerolog.Logger
}

type DBConfig struct {
	DBName string
	DBUser string
	DBPass string
	DBHost string
	DBPort int
}

func NewDatabase(cnf *DBConfig) *Database {

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return &Database{
		DBconfig: *cnf,
		Log:      zerolog.New(output).With().Timestamp().Logger(),
		LogLvl:   logger.Error,
	}
}

func (db *Database) Init() *Database {
	db.Engine = db.DBadapter.Setup(db.Log, db.LogLvl)
	return db
}

func (db *Database) Config(dbname string, dbuser string, dbpass string, dbhost string, dbport int) *Database {
	db.DBconfig = DBConfig{
		DBName: dbname,
		DBUser: dbuser,
		DBPass: dbpass,
		DBHost: dbhost,
		DBPort: dbport,
	}
	return db
}

func (db *Database) Adapter(adapter DBAdapter) *Database {
	db.DBadapter = adapter.Init()
	return db
}

func (db *Database) LogLevel(logLevel logger.LogLevel) *Database {
	db.LogLvl = logLevel
	return db
}

func (db *Database) Logger(log zerolog.Logger) *Database {
	db.Log = log
	return db
}
