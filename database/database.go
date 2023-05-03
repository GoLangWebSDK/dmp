package database

import (
	//
	// I made a little poopoo here. The dmp/database package is now dependent on another package of mine: logger.
	// Firstly, I created logger as a temporary package just to handle the initial dependencies that came from
	// the initial stack. Secondly, I don't like the idea of our packages being codependent on each other.
	// Maybe we can make them dependent on the logger package, but I think that's not a good idea. Thirdly,
	// the package will be rewritten or expanded anyway into the erlog package. So, maybe we can create an internal
	// logging system for each package that can be configured from the app package, or make all packages
	// dependent on our custom error and logging package.
	//
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// How to handle database connection in a controller?
// One way is to create a global variable and assign it in main.go
// then use directly inside a repostry.

// Global variable for database connection
var DBManager *Database

type DBAdapter interface {
	Init(log zerolog.Logger, logLevel logger.LogLevel) *gorm.DB
	Setup(config *DBConfig) DBAdapter
	// NewMock() (sqlmock.Sqlmock, *gorm.DB)
}

type Adapter struct {
	db     *gorm.DB
	config DBConfig
	dsn    string
	err    error
}

type Database struct {
	Engine    *gorm.DB
	DBadapter DBAdapter
	DBconfig  DBConfig
	LogLvl    logger.LogLevel
	log       zerolog.Logger
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
		log:      zerolog.New(output).With().Timestamp().Logger(),
		LogLvl:   logger.Error,
	}
}

func (db *Database) Init() *Database {
	db.Engine = db.DBadapter.Init(db.log, db.LogLvl)
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
	// what is this? if db.config == (DBConfig{}) ??
	if db.DBconfig == (DBConfig{}) {
		db.log.Error().Msg("Database config not set.")
		return db
	}

	db.DBadapter = adapter.Setup(&db.DBconfig)
	return db
}

func (db *Database) LogLevel(logLevel logger.LogLevel) *Database {
	db.LogLvl = logLevel
	return db
}

func (db *Database) Logger(log zerolog.Logger) *Database {
	db.log = log
	return db
}

func (db *Database) Log(log zerolog.Logger, logLevel logger.LogLevel) *Database {
	db.log = log
	db.LogLvl = logLevel
	return db
}
