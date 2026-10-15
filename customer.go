package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// createCustomerDB adds a new customer to the database
func createCustomerDB(db *gorm.DB, customer *Customer) error {
	if customer == nil {
		return errors.New("customer cannot be nil")
	}
	if strings.TrimSpace(customer.Name) == "" {
		return errors.New("customer name cannot be empty")
	}
	if strings.TrimSpace(customer.Email) == "" {
		return errors.New("customer email cannot be empty")
	}

	return db.Create(customer).Error
}

// getCustomerByEmailDB finds a customer by their email address
func getCustomerByEmailDB(db *gorm.DB, email string) (*Customer, error) {
	var customer Customer
	err := db.Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// getCustomerByIDDB finds a customer by ID
func getCustomerByIDDB(db *gorm.DB, id uint) (*Customer, error) {
	var customer Customer
	err := db.First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// getAllCustomersDB retrieves all customers
func getAllCustomersDB(db *gorm.DB) ([]Customer, error) {
	var customers []Customer
	err := db.Find(&customers).Error
	return customers, err
}

// updateCustomerByIDDB updates fields of a customer by ID
func updateCustomerByIDDB(db *gorm.DB, id uint, updates map[string]interface{}) error {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return fmt.Errorf("customer ID %d not found: %w", id, err)
	}

	return db.Model(&customer).Updates(updates).Error
}

// deleteCustomerByIDDB deletes a customer by ID.
// Rule: Cannot delete a customer who has existing orders in order history.
func deleteCustomerByIDDB(db *gorm.DB, id uint) error {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return fmt.Errorf("customer ID %d not found: %w", id, err)
	}

	// Check if this customer has existing orders
	var orderCount int64
	if err := db.Model(&Order{}).Where("customer_id = ?", id).Count(&orderCount).Error; err != nil {
		return fmt.Errorf("failed to check customer orders: %w", err)
	}
	if orderCount > 0 {
		return fmt.Errorf("cannot delete customer '%s' (ID %d): customer has %d order(s) in history", customer.Name, id, orderCount)
	}

	return db.Delete(&customer).Error
}
