package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	once     sync.Once
	instance *gorm.DB
)

// GetDB tries to connect to the database with a retry mechanism
func GetDB() *gorm.DB {
	once.Do(func() {
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")

		// Ensure all necessary environment variables are set
		if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" {
			log.Fatal("One or more required database environment variables are not set")
		}

		// DSN (Data Source Name)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
		const maxRetries = 5
		const retryInterval = 5 * time.Second // Adjust the retry interval as needed

		for i := 0; i < maxRetries; i++ {
			var err error
			instance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Printf("Failed to connect to database, attempt %d/%d: %v", i+1, maxRetries, err)
				time.Sleep(retryInterval)
				continue
			}
			return
		}
		log.Fatalf("Failed to connect to database after %d attempts", maxRetries)
	})

	return instance
}
