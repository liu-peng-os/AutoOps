package db

import (
	"dodevops-api/common/config"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

// NewGormLogger creates a GORM logger that writes to both stdout and logs/app.log.
func NewGormLogger() logger.Interface {
	os.MkdirAll("logs", os.ModePerm)

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("unable to open GORM log file: %v", err)
		return logger.Default.LogMode(logger.Info)
	}

	mw := io.MultiWriter(os.Stdout, file)
	return logger.New(
		log.New(mw, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
}

// SetupDBLink initializes the primary application database connection.
func SetupDBLink() error {
	dbConfig := config.Config.Db

	database, err := openGormDB(dbConfig)
	if err != nil {
		return err
	}
	if database.Error != nil {
		return database.Error
	}

	Db = database

	if err := RunMigrations(Db, dbConfig.MigrationPath); err != nil {
		return err
	}

	// The config flag is still accepted, but phase1 schema bootstrapping is now goose-first.
	if dbConfig.AutoMigrate {
		if err := AutoMigrate(Db); err != nil {
			return err
		}
	}

	sqlDB, err := Db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpen)
	return nil
}

func openGormDB(dbConfig config.Db) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(dbConfig.Dialects) {
	case "", "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&sql_mode=''",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Db,
			dbConfig.Charset,
		)
		dialector = mysql.Open(dsn)
	case "postgres", "postgresql":
		sslMode := dbConfig.SSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Db,
			sslMode,
		)
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", dbConfig.Dialects)
	}

	const maxAttempts = 3
	const retryDelay = 2 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		database, err := gorm.Open(dialector, &gorm.Config{
			Logger:                                   NewGormLogger(),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			if database == nil {
				return nil, errors.New("gorm returned nil database instance")
			}
			return database, nil
		}

		lastErr = err
		if !shouldRetryDBConnect(err) || attempt == maxAttempts {
			return nil, err
		}

		log.Printf("database connection attempt %d/%d failed: %v; retrying in %s", attempt, maxAttempts, err, retryDelay)
		time.Sleep(retryDelay)
	}

	return nil, lastErr
}

func shouldRetryDBConnect(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	retryHints := []string{
		"unexpected eof",
		"connection reset by peer",
		"broken pipe",
		"i/o timeout",
		"timeout",
		"temporary failure",
		"connection refused",
		"no route to host",
	}

	for _, hint := range retryHints {
		if strings.Contains(message, hint) {
			return true
		}
	}

	return false
}
