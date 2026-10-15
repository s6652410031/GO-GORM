package main

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Valid order statuses
var validOrderStatuses = map[string]bool{
	"Pending":   true,
	"Completed": true,
	"Cancelled": true,
}

// createOrderDB creates an order within a database transaction.
// It verifies the customer and product, validates available stock,
// calculates TotalPrice = Price * Quantity, deducts stock from the product,
// and ensures both the order creation and stock reduction succeed or roll back together.
func createOrderDB(db *gorm.DB, customerID uint, productID uint, quantity int) (*Order, error) {
	if quantity <= 0 {
		return nil, errors.New("order quantity must be greater than zero")
	}

	var createdOrder Order

	// Execute inside a database transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Verify customer exists
		var customer Customer
		if err := tx.First(&customer, customerID).Error; err != nil {
			return fmt.Errorf("customer ID %d not found: %w", customerID, err)
		}

		// 2. Verify product exists
		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			return fmt.Errorf("product ID %d not found: %w", productID, err)
		}

		// 3. Verify stock availability
		if product.Quantity < quantity {
			return fmt.Errorf("insufficient stock for product '%s' (ID %d): requested %d, in stock %d",
				product.Name, product.ID, quantity, product.Quantity)
		}

		// 4. Calculate total price
		totalPrice := product.Price * float64(quantity)

		// 5. Deduct stock from product
		product.Quantity -= quantity
		if err := tx.Save(&product).Error; err != nil {
			return fmt.Errorf("failed to deduct product stock: %w", err)
		}

		// 6. Create the order record
		createdOrder = Order{
			CustomerID: customerID,
			ProductID:  productID,
			Quantity:   quantity,
			TotalPrice: totalPrice,
			Status:     "Pending",
		}
		if err := tx.Create(&createdOrder).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Preload associations for convenience
	_ = db.Preload("Customer").Preload("Product").First(&createdOrder, createdOrder.ID)
	return &createdOrder, nil
}

// getOrdersByCustomerDB retrieves all orders placed by a specific customer
func getOrdersByCustomerDB(db *gorm.DB, customerID uint) ([]Order, error) {
	var orders []Order
	err := db.Preload("Customer").Preload("Product").
		Where("customer_id = ?", customerID).
		Find(&orders).Error
	return orders, err
}

// getOrdersByStatusDB retrieves orders filtered by status ("Pending", "Completed", "Cancelled")
func getOrdersByStatusDB(db *gorm.DB, status string) ([]Order, error) {
	var orders []Order
	err := db.Preload("Customer").Preload("Product").
		Where("status = ?", status).
		Find(&orders).Error
	return orders, err
}

// getOrdersByDateRangeDB retrieves orders created within a specific date range
func getOrdersByDateRangeDB(db *gorm.DB, startDate, endDate time.Time) ([]Order, error) {
	var orders []Order
	err := db.Preload("Customer").Preload("Product").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Find(&orders).Error
	return orders, err
}

// getOrdersWithDetailsDB retrieves all orders with full Customer and Product details
func getOrdersWithDetailsDB(db *gorm.DB) ([]Order, error) {
	var orders []Order
	err := db.Preload("Customer").Preload("Product").Preload("Product.Category").
		Find(&orders).Error
	return orders, err
}

// updateOrderStatusDB updates the status of an order.
// Status must strictly be "Pending", "Completed", or "Cancelled".
// If status is changed to "Cancelled", the reserved stock is restored to the product in a transaction.
func updateOrderStatusDB(db *gorm.DB, id uint, newStatus string) error {
	if !validOrderStatuses[newStatus] {
		return fmt.Errorf("invalid status '%s': status must be 'Pending', 'Completed', or 'Cancelled'", newStatus)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var order Order
		if err := tx.First(&order, id).Error; err != nil {
			return fmt.Errorf("order ID %d not found: %w", id, err)
		}

		oldStatus := order.Status
		if oldStatus == newStatus {
			return nil // no change
		}

		// If cancelling an order, restore product stock
		if newStatus == "Cancelled" && oldStatus != "Cancelled" {
			var product Product
			if err := tx.First(&product, order.ProductID).Error; err == nil {
				product.Quantity += order.Quantity
				if err := tx.Save(&product).Error; err != nil {
					return fmt.Errorf("failed to restore stock for cancelled order: %w", err)
				}
			}
		}

		// If reactivating a cancelled order, check and deduct stock
		if oldStatus == "Cancelled" && newStatus != "Cancelled" {
			var product Product
			if err := tx.First(&product, order.ProductID).Error; err != nil {
				return fmt.Errorf("cannot reactivate order, product not found: %w", err)
			}
			if product.Quantity < order.Quantity {
				return fmt.Errorf("cannot reactivate order: insufficient stock (available %d, required %d)", product.Quantity, order.Quantity)
			}
			product.Quantity -= order.Quantity
			if err := tx.Save(&product).Error; err != nil {
				return fmt.Errorf("failed to deduct stock for reactivated order: %w", err)
			}
		}

		order.Status = newStatus
		return tx.Save(&order).Error
	})
}
