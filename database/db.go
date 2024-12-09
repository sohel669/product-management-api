// package database

// import (
//     "fmt"
//     "log"
//     "os"
//     "time"

//     "gorm.io/driver/postgres"
//     "gorm.io/gorm"
//     "gorm.io/gorm/logger"
// )

// var DB *gorm.DB

// // ConnectDB initializes the database connection
// func ConnectDB() {
//     // Load environment variables or hardcode connection string
//     host := os.Getenv("DB_HOST")      // e.g., "localhost"
//     port := os.Getenv("DB_PORT")      // e.g., "5432"
//     user := os.Getenv("DB_USER")      // e.g., "postgres"
//     password := os.Getenv("DB_PASS")  // e.g., "Sohel@123"
//     dbname := os.Getenv("DB_NAME")    // e.g., "Project"
//     sslmode := os.Getenv("DB_SSL")    // e.g., "disable"

//     dsn := fmt.Sprintf(
//         "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Kolkata",
//         host, port, user, password, dbname, sslmode,
//     )

//     var err error
//     DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
//         Logger: logger.Default.LogMode(logger.Info), // Enable detailed logs for debugging
//     })
//     if err != nil {
//         log.Fatalf("Failed to connect to the database: %v", err)
//     }

//     // Optional: Verify database connection with a ping
//     sqlDB, err := DB.DB()
//     if err != nil {
//         log.Fatalf("Error getting SQL DB instance: %v", err)
//     }

//     sqlDB.SetMaxIdleConns(10)               // Set idle connections
//     sqlDB.SetMaxOpenConns(100)              // Set max open connections
//     sqlDB.SetConnMaxLifetime(1 * time.Hour) // Set connection lifetime

//     if err := sqlDB.Ping(); err != nil {
//         log.Fatalf("Database connection is not alive: %v", err)
//     }

//     log.Println("Database connected successfully!")
// }

// // GetDB provides the database instance
// func GetDB() *gorm.DB {
//     return DB
// }
package main

import (
	"log"
)

func InitDB() {
	log.Println("Initializing Database...")


	log.Println("Database initialized successfully")
}

