package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
//DB.AutoMigrate(&User{}, &Task{})
func InitDB() error {
	dsn := "host=localhost user=postgres password=yourpassword dbname=taskmgr port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//var err error
	
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}


	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	err = db.AutoMigrate(&User{}, &Task{})
    if err != nil {
        return err
    }

    DB = db
    fmt.Println("Database connected and migrated!")
    return nil
}
