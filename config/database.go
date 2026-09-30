package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dbName := "minitube"

	// Connect without database name first to create it if missing
	rootDSN := "root:@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"

	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to MySQL server:", err)
	}

	sql := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;",
		dbName,
	)

	if err := rootDB.Exec(sql).Error; err != nil {
		log.Fatal("Could not create database:", err)
	}

	// Now connect to the actual database
	dsn := fmt.Sprintf(
		"root:@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Database connected successfully")

	DB = db
}