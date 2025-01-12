package database

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bahramiofficial/watchtower/src/utilities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

type Configdb struct {
	Host            string
	Port            int
	User            string
	Password        string
	DbName          string
	SslMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

func InitDb() error {
	port, _ := strconv.Atoi(utilities.GetValueEnv("DB_HOST"))
	cfg := &Configdb{
		Host:            utilities.GetValueEnv("DB_HOST"),
		Port:            port,
		User:            utilities.GetValueEnv("DB_USER"),
		Password:        utilities.GetValueEnv("DB_PASS"),
		DbName:          utilities.GetValueEnv("DB_NAME"),
		SslMode:         "disable",
		MaxIdleConns:    15,
		MaxOpenConns:    100,
		ConnMaxLifetime: 5,
	}

	var err error
	cnn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SslMode)

	dbClient, err = gorm.Open(postgres.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDb, _ := dbClient.DB()
	err = sqlDb.Ping()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	log.Println("Db connection established")
	return nil
}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	con, _ := dbClient.DB()
	con.Close()
}

// GetDbAfterInit ensures the database is initialized before being used.
// The connection will be closed after the function finishes executing.
func GetDbAfterInit() (*gorm.DB, func(), error) {
	// Initialize database connection if it's not already initialized
	if dbClient == nil {
		err := InitDb()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to initialize database: %w", err)
		}
	}

	// Return the database connection and a function to close the connection
	return dbClient, func() {
		CloseDb() // This will be deferred and called after the job is completed
	}, nil
}
