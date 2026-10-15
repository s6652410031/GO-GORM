package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Category represents a product category (1:N with Product)
type Category struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"products,omitempty"`
}

// Product represents an item in the store (1:N with Order, N:1 with Category)
type Product struct {
	gorm.Model
	Name       string   `gorm:"type:varchar(150);not null" json:"name"`
	Price      float64  `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	CategoryID uint     `gorm:"not null;index" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Orders     []Order  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"orders,omitempty"`
}

// Customer represents a store client (1:N with Order)
type Customer struct {
	gorm.Model
	Name   string  `gorm:"type:varchar(100);not null" json:"name"`
	Email  string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone  string  `gorm:"type:varchar(20)" json:"phone"`
	Orders []Order `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"orders,omitempty"`
}

// Order represents an order placed by a customer for a single product (N:1 Customer, N:1 Product)
type Order struct {
	gorm.Model
	CustomerID uint     `gorm:"not null;index" json:"customer_id"`
	ProductID  uint     `gorm:"not null;index" json:"product_id"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	TotalPrice float64  `gorm:"type:decimal(12,2);not null" json:"total_price"`
	Status     string   `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"` // Pending, Completed, Cancelled
	Customer   Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Product    Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// InitDB connects to MySQL, creates database store_db if not exists, and performs AutoMigrate.
func InitDB(user, pass, host, port, dbName string) (*gorm.DB, error) {
	// 1. Connect without dbName to ensure store_db exists
	var rootDSN string
	if pass == "" {
		rootDSN = fmt.Sprintf("%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", user, host, port)
	} else {
		rootDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port)
	}

	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
	}

	// Create database if not exists
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName)
	if err := rootDB.Exec(createDBSQL).Error; err != nil {
		return nil, fmt.Errorf("failed to create database %s: %w", dbName, err)
	}

	// Close root connection
	sqlDB, err := rootDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	// 2. Connect to the specific database
	var dbDSN string
	if pass == "" {
		dbDSN = fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, host, port, dbName)
	} else {
		dbDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)
	}

	db, err := gorm.Open(mysql.Open(dbDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", dbName, err)
	}

	// 3. AutoMigrate all 4 tables
	log.Printf("[Database] Running AutoMigrate for Category, Product, Customer, Order...")
	err = db.AutoMigrate(&Category{}, &Product{}, &Customer{}, &Order{})
	if err != nil {
		return nil, fmt.Errorf("auto-migration failed: %w", err)
	}

	log.Printf("[Database] Connected and migrated store_db successfully.")
	return db, nil
}
