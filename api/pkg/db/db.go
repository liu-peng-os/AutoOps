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

	// Phase 1 keeps a compatibility fallback so an empty PostgreSQL database can still boot.
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

	database, err := gorm.Open(dialector, &gorm.Config{
		Logger:                                   NewGormLogger(),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	if database == nil {
		return nil, errors.New("gorm returned nil database instance")
	}

	return database, nil
}
