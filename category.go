package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// createCategoryDB adds a new category to the database
func createCategoryDB(db *gorm.DB, category *Category) error {
	if category == nil {
		return errors.New("category cannot be nil")
	}
	if category.Name == "" {
		return errors.New("category name cannot be empty")
	}
	result := db.Create(category)
	return result.Error
}

// getAllCategoriesDB retrieves all categories
func getAllCategoriesDB(db *gorm.DB) ([]Category, error) {
	var categories []Category
	err := db.Find(&categories).Error
	return categories, err
}

// getCategoryByIDDB retrieves a category by its ID
func getCategoryByIDDB(db *gorm.DB, id uint) (*Category, error) {
	var category Category
	err := db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// deleteCategoryByIDDB deletes a category by its ID.
// Rule: Cannot delete a category if products are still associated with it.
func deleteCategoryByIDDB(db *gorm.DB, id uint) error {
	// 1. Check if category exists
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		return fmt.Errorf("category ID %d not found: %w", id, err)
	}

	// 2. Check if any products exist under this category
	var productCount int64
	if err := db.Model(&Product{}).Where("category_id = ?", id).Count(&productCount).Error; err != nil {
		return fmt.Errorf("failed to check existing products: %w", err)
	}
	if productCount > 0 {
		return fmt.Errorf("cannot delete category '%s' (ID %d): %d product(s) still exist in this category", category.Name, id, productCount)
	}

	// 3. Delete the category
	return db.Delete(&category).Error
}
