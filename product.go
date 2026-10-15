package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// createProductDB adds a new product to the database
func createProductDB(db *gorm.DB, product *Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	if product.Name == "" {
		return errors.New("product name cannot be empty")
	}
	if product.Price < 0 {
		return errors.New("product price cannot be negative")
	}
	if product.Quantity < 0 {
		return errors.New("product quantity cannot be negative")
	}

	// Verify that the Category exists
	var category Category
	if err := db.First(&category, product.CategoryID).Error; err != nil {
		return fmt.Errorf("category ID %d does not exist: %w", product.CategoryID, err)
	}

	return db.Create(product).Error
}

// getAllProductsDB retrieves all active products with their category preloaded
func getAllProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Find(&products).Error
	return products, err
}

// getProductByIDDB retrieves a product by its ID
func getProductByIDDB(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	err := db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// getProductByNameDB searches for products whose name contains the search term (LIKE)
func getProductByNameDB(db *gorm.DB, name string) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Where("name LIKE ?", "%"+name+"%").Find(&products).Error
	return products, err
}

// getProductsByCategoryDB retrieves products by category ID
func getProductsByCategoryDB(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// getProductsByPriceRangeDB retrieves products within a min and max price range
func getProductsByPriceRangeDB(db *gorm.DB, minPrice, maxPrice float64) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products).Error
	return products, err
}

// getProductsLowStockDB retrieves products where quantity is strictly less than threshold
func getProductsLowStockDB(db *gorm.DB, threshold int) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Where("quantity < ?", threshold).Find(&products).Error
	return products, err
}

// getProductsByCategoryAndPriceDB retrieves products matching both category ID and price range
func getProductsByCategoryAndPriceDB(db *gorm.DB, categoryID uint, minPrice, maxPrice float64) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").
		Where("category_id = ? AND price BETWEEN ? AND ?", categoryID, minPrice, maxPrice).
		Find(&products).Error
	return products, err
}

// getProductsSortedByPriceDB retrieves products sorted by price ("ASC" or "DESC")
func getProductsSortedByPriceDB(db *gorm.DB, order string) ([]Product, error) {
	var products []Product
	orderClause := "price ASC"
	if strings.ToUpper(strings.TrimSpace(order)) == "DESC" {
		orderClause = "price DESC"
	}
	err := db.Preload("Category").Order(orderClause).Find(&products).Error
	return products, err
}

// getTopExpensiveProductsDB retrieves the top N most expensive products
func getTopExpensiveProductsDB(db *gorm.DB, limit int) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Order("price DESC").Limit(limit).Find(&products).Error
	return products, err
}

// updateProductByIDDB updates fields of a product by ID using a map or struct
func updateProductByIDDB(db *gorm.DB, id uint, updates map[string]interface{}) error {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	// Validate price if present in updates
	if priceVal, ok := updates["price"]; ok {
		switch p := priceVal.(type) {
		case float64:
			if p < 0 {
				return errors.New("price cannot be negative")
			}
		case int:
			if p < 0 {
				return errors.New("price cannot be negative")
			}
		}
	}

	// Validate quantity if present in updates
	if qtyVal, ok := updates["quantity"]; ok {
		switch q := qtyVal.(type) {
		case int:
			if q < 0 {
				return errors.New("quantity cannot be negative")
			}
		case float64:
			if q < 0 {
				return errors.New("quantity cannot be negative")
			}
		}
	}

	return db.Model(&product).Updates(updates).Error
}

// updateProductPriceDB updates the price of a product (must not be negative)
func updateProductPriceDB(db *gorm.DB, id uint, newPrice float64) error {
	if newPrice < 0 {
		return errors.New("price cannot be negative")
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	return db.Model(&product).Update("price", newPrice).Error
}

// updateProductStockDB updates the stock quantity of a product (must not be negative)
func updateProductStockDB(db *gorm.DB, id uint, newQuantity int) error {
	if newQuantity < 0 {
		return errors.New("stock quantity cannot be negative")
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	return db.Model(&product).Update("quantity", newQuantity).Error
}

// softDeleteProductDB performs a soft delete on a product
func softDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	return db.Delete(&product).Error
}

// hardDeleteProductDB permanently deletes a product from the database.
// Rule: Cannot hard delete if the product has orders in order history.
func hardDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.Unscoped().First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	// Check if this product has existing orders
	var orderCount int64
	if err := db.Model(&Order{}).Where("product_id = ?", id).Count(&orderCount).Error; err != nil {
		return fmt.Errorf("failed to verify order history: %w", err)
	}
	if orderCount > 0 {
		return fmt.Errorf("cannot permanently delete product '%s' (ID %d): referenced in %d order(s)", product.Name, id, orderCount)
	}

	return db.Unscoped().Delete(&product).Error
}

// getDeletedProductsDB retrieves all soft-deleted products
func getDeletedProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Unscoped().Where("deleted_at IS NOT NULL").Find(&products).Error
	return products, err
}

// restoreProductDB restores a soft-deleted product
func restoreProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.Unscoped().First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	if !product.DeletedAt.Valid {
		return fmt.Errorf("product ID %d is not deleted", id)
	}

	return db.Unscoped().Model(&product).Update("deleted_at", nil).Error
}
