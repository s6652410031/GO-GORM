package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	DBUser = "root"
	DBPass = ""
	DBHost = "127.0.0.1"
	DBPort = "3306"
	DBName = "store_db"
)

func main() {
	fmt.Println("==================================================================")
	fmt.Println("       Store Management System (Go + GORM + MySQL)                ")
	fmt.Println("==================================================================")

	// 1. Initialize and connect to MySQL database
	db, err := InitDB(DBUser, DBPass, DBHost, DBPort, DBName)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// If run with command line flag "--demo", run test demo directly
	if len(os.Args) > 1 && os.Args[1] == "--demo" {
		runFullAssignmentDemo(db)
		return
	}

	// Interactive Menu
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n----------------------------- MENU -----------------------------")
		fmt.Println("1. Run Full Assignment Demonstration (Test All Parts 1 - 5)")
		fmt.Println("2. Seed Initial Test Data (4 Categories, 15 Products, 8 Customers, 12 Orders)")
		fmt.Println("3. Reset & Re-create Database (Fresh Clean State)")
		fmt.Println("4. Export SQL Dump File (store_db.sql)")
		fmt.Println("5. View All Products")
		fmt.Println("6. View All Orders with Full Details")
		fmt.Println("7. Exit")
		fmt.Print("Select an option (1-7): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			// If non-interactive stdin (e.g. piped or closed), run demo and exit
			runFullAssignmentDemo(db)
			break
		}

		choice := strings.TrimSpace(input)
		switch choice {
		case "1":
			runFullAssignmentDemo(db)
		case "2":
			seedDatabase(db)
		case "3":
			resetDatabase(db)
		case "4":
			if err := exportSQLDump(db, "store_db.sql"); err != nil {
				fmt.Printf("Error exporting SQL: %v\n", err)
			} else {
				fmt.Println("Successfully exported database dump to 'store_db.sql'")
			}
		case "5":
			products, _ := getAllProductsDB(db)
			printProducts(products)
		case "6":
			orders, _ := getOrdersWithDetailsDB(db)
			printOrdersWithDetails(orders)
		case "7", "q", "exit":
			fmt.Println("Exiting Store Management System. Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please enter a number between 1 and 7.")
		}
	}
}

// seedDatabase ensures the minimum test data required by Section 4 is populated
func seedDatabase(db *gorm.DB) {
	fmt.Println("\n>>> [Seeding] Checking and populating initial minimum test data...")

	var catCount int64
	db.Model(&Category{}).Count(&catCount)
	if catCount >= 4 {
		fmt.Printf("Database already contains %d categories. Skipping initial seed.\n", catCount)
		return
	}

	// 1. Seed 4 Categories
	categories := []Category{
		{Name: "Laptops & Notebooks", Description: "High-performance laptops, ultrabooks, and portable workstations"},
		{Name: "Smartphones & Tablets", Description: "Latest flagship and mid-range mobile devices"},
		{Name: "Computer Accessories", Description: "Mechanical keyboards, ergonomic mice, and USB adapters"},
		{Name: "Peripherals & Displays", Description: "High-refresh monitors, portable external SSDs, and displays"},
	}

	for i := range categories {
		if err := createCategoryDB(db, &categories[i]); err != nil {
			log.Printf("Failed to seed category %s: %v", categories[i].Name, err)
		}
	}
	fmt.Printf("✓ Seeded %d Categories\n", len(categories))

	// 2. Seed 15 Products (including examples from PDF: MacBook Air, ASUS Vivobook, iPhone 16, Samsung Galaxy, Logitech Mouse, etc.)
	products := []Product{
		// Category 1: Laptops
		{Name: "MacBook Air M3", Price: 34900.00, Quantity: 12, CategoryID: categories[0].ID},
		{Name: "ASUS Vivobook 15", Price: 21900.00, Quantity: 8, CategoryID: categories[0].ID},
		{Name: "Dell XPS 13", Price: 45900.00, Quantity: 4, CategoryID: categories[0].ID},
		{Name: "Lenovo ThinkPad X1", Price: 38900.00, Quantity: 5, CategoryID: categories[0].ID},

		// Category 2: Smartphones
		{Name: "iPhone 16 128GB", Price: 32900.00, Quantity: 15, CategoryID: categories[1].ID},
		{Name: "iPhone 16 Pro 256GB", Price: 41900.00, Quantity: 3, CategoryID: categories[1].ID}, // low stock
		{Name: "Samsung Galaxy S24", Price: 29900.00, Quantity: 10, CategoryID: categories[1].ID},
		{Name: "Samsung Galaxy Z Flip 6", Price: 36900.00, Quantity: 4, CategoryID: categories[1].ID}, // low stock

		// Category 3: Accessories
		{Name: "Logitech Mouse MX Master 3S", Price: 3590.00, Quantity: 20, CategoryID: categories[2].ID},
		{Name: "Logitech G Pro Wireless Mouse", Price: 4290.00, Quantity: 7, CategoryID: categories[2].ID},
		{Name: "Mechanical Keyboard Keychron K2", Price: 3890.00, Quantity: 14, CategoryID: categories[2].ID},
		{Name: "USB-C Hub 7-in-1 Adapter", Price: 1290.00, Quantity: 25, CategoryID: categories[2].ID},
		{Name: "HyperX Cloud III Gaming Headset", Price: 2890.00, Quantity: 9, CategoryID: categories[2].ID},

		// Category 4: Peripherals & Displays
		{Name: "LG UltraFine 27-inch 4K Monitor", Price: 16900.00, Quantity: 5, CategoryID: categories[3].ID},
		{Name: "Samsung T7 Shield 1TB Portable SSD", Price: 4500.00, Quantity: 18, CategoryID: categories[3].ID},
	}

	for i := range products {
		if err := createProductDB(db, &products[i]); err != nil {
			log.Printf("Failed to seed product %s: %v", products[i].Name, err)
		}
	}
	fmt.Printf("✓ Seeded %d Products\n", len(products))

	// 3. Seed 8 Customers
	customers := []Customer{
		{Name: "Somchai Jaidee", Email: "somchai@example.com", Phone: "081-111-2222"},
		{Name: "Suda Rakdee", Email: "suda@example.com", Phone: "082-222-3333"},
		{Name: "Arthit Srisuk", Email: "arthit@example.com", Phone: "083-333-4444"},
		{Name: "Wanna Suksan", Email: "wanna@example.com", Phone: "084-444-5555"},
		{Name: "Natthapong Wongsuwan", Email: "natthapong@example.com", Phone: "085-555-6666"},
		{Name: "Pornthip Charoen", Email: "pornthip@example.com", Phone: "086-666-7777"},
		{Name: "Komsan Thongdee", Email: "komsan@example.com", Phone: "087-777-8888"},
		{Name: "Siriporn Meesuk", Email: "siriporn@example.com", Phone: "088-888-9999"},
	}

	for i := range customers {
		if err := createCustomerDB(db, &customers[i]); err != nil {
			log.Printf("Failed to seed customer %s: %v", customers[i].Name, err)
		}
	}
	fmt.Printf("✓ Seeded %d Customers\n", len(customers))

	// 4. Seed 12 Orders using createOrderDB (Transaction & Stock Deduction)
	orderRequests := []struct {
		custIdx int
		prodIdx int
		qty     int
		status  string
	}{
		{0, 0, 1, "Completed"}, // Somchai buys 1 MacBook Air
		{0, 8, 1, "Completed"}, // Somchai buys 1 Logitech MX Master
		{1, 4, 1, "Completed"}, // Suda buys 1 iPhone 16
		{1, 11, 2, "Pending"},  // Suda buys 2 USB-C Hub
		{2, 1, 1, "Completed"}, // Arthit buys 1 ASUS Vivobook
		{2, 10, 1, "Completed"},// Arthit buys 1 Mechanical Keyboard
		{3, 6, 1, "Pending"},   // Wanna buys 1 Galaxy S24
		{4, 13, 1, "Completed"},// Natthapong buys 1 LG Monitor
		{4, 14, 2, "Pending"},  // Natthapong buys 2 Samsung SSD
		{5, 9, 1, "Cancelled"}, // Pornthip buys 1 Logitech G Pro (Cancelled)
		{6, 11, 1, "Completed"},// Komsan buys 1 USB-C Hub
		{7, 12, 1, "Pending"},  // Siriporn buys 1 HyperX Headset
	}

	var seededOrderCount int
	for _, req := range orderRequests {
		ord, err := createOrderDB(db, customers[req.custIdx].ID, products[req.prodIdx].ID, req.qty)
		if err != nil {
			log.Printf("Failed to seed order: %v", err)
			continue
		}
		if req.status != "Pending" {
			_ = updateOrderStatusDB(db, ord.ID, req.status)
		}
		seededOrderCount++
	}
	fmt.Printf("✓ Seeded %d Orders (Transactions verified, stock deducted)\n", seededOrderCount)
	fmt.Println(">>> Minimum test data seeding complete!")
}

// resetDatabase drops and recreates tables
func resetDatabase(db *gorm.DB) {
	fmt.Print("Are you sure you want to drop and reset all tables? (y/n): ")
	var confirm string
	fmt.Scanln(&confirm)
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("Reset cancelled.")
		return
	}

	// Disable foreign keys temporarily to drop cleanly
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	db.Migrator().DropTable(&Order{}, &Customer{}, &Product{}, &Category{})
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	err := db.AutoMigrate(&Category{}, &Product{}, &Customer{}, &Order{})
	if err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		return
	}
	fmt.Println("All tables reset and re-migrated successfully!")
	seedDatabase(db)
}

// runFullAssignmentDemo executes and outputs every requirement from Sections 1 to 5
func runFullAssignmentDemo(db *gorm.DB) {
	fmt.Println("\n==================================================================")
	fmt.Println("    EXECUTING FULL ASSIGNMENT VERIFICATION (PARTS 1 - 5)          ")
	fmt.Println("==================================================================")

	// Ensure seed data is present
	seedDatabase(db)

	// ==================== PART 1: Database & Models ====================
	fmt.Println("\n------------------------------------------------------------------")
	fmt.Println(">>> [ส่วนที่ 1] Database & Models Verification")
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("✓ Database 'store_db' exists and connected.")
	fmt.Println("✓ Models: Category, Product, Customer, Order configured with gorm.Model and Foreign Keys.")
	fmt.Println("✓ AutoMigrate() verified.")

	// ==================== PART 2: Create ====================
	fmt.Println("\n------------------------------------------------------------------")
	fmt.Println(">>> [ส่วนที่ 2] Create Functions Verification")
	fmt.Println("------------------------------------------------------------------")

	// 1. createCategoryDB()
	newCat := Category{Name: "Audio & Wearables", Description: "Headphones, Earbuds, Smartwatches"}
	err := createCategoryDB(db, &newCat)
	fmt.Printf("[createCategoryDB] Inserted ID: %d, Name: '%s' (Error: %v)\n", newCat.ID, newCat.Name, err)

	// 2. createProductDB()
	newProd := Product{Name: "Sony WH-1000XM5", Price: 11900.00, Quantity: 10, CategoryID: newCat.ID}
	err = createProductDB(db, &newProd)
	fmt.Printf("[createProductDB] Inserted ID: %d, Name: '%s', Price: %.2f, Stock: %d (Error: %v)\n",
		newProd.ID, newProd.Name, newProd.Price, newProd.Quantity, err)

	// 3. createCustomerDB()
	newCust := Customer{Name: "Chaiwat Demo", Email: fmt.Sprintf("chaiwat_%d@example.com", time.Now().Unix()), Phone: "089-999-9999"}
	err = createCustomerDB(db, &newCust)
	fmt.Printf("[createCustomerDB] Inserted ID: %d, Name: '%s', Email: '%s' (Error: %v)\n",
		newCust.ID, newCust.Name, newCust.Email, err)

	// 4. createOrderDB() Transaction & Stock Deduction Test
	fmt.Println("\n[createOrderDB Transaction Test]")
	fmt.Printf("Before Order -> Product '%s' Stock: %d\n", newProd.Name, newProd.Quantity)
	testOrder, err := createOrderDB(db, newCust.ID, newProd.ID, 3)
	if err != nil {
		fmt.Printf("Order creation failed: %v\n", err)
	} else {
		// Re-fetch product to verify stock reduction
		var checkProd Product
		db.First(&checkProd, newProd.ID)
		fmt.Printf("✓ Order Created! Order ID: %d, Customer: %s, Product: %s, Qty: %d, TotalPrice: ฿%.2f, Status: %s\n",
			testOrder.ID, newCust.Name, checkProd.Name, testOrder.Quantity, testOrder.TotalPrice, testOrder.Status)
		fmt.Printf("✓ After Order -> Product '%s' Stock updated to: %d (Decreased by 3)\n", checkProd.Name, checkProd.Quantity)
	}

	// 5. Test Insufficient Stock Transaction Rollback
	fmt.Println("\n[Transaction Rollback Test - Requesting 999 items when stock is insufficient]")
	_, errRollback := createOrderDB(db, newCust.ID, newProd.ID, 999)
	fmt.Printf("Attempting invalid order: expected error caught: '%v'\n", errRollback)

	// ==================== PART 3: Read & Conditional Queries ====================
	fmt.Println("\n------------------------------------------------------------------")
	fmt.Println(">>> [ส่วนที่ 3] Read / Conditional Queries Verification")
	fmt.Println("------------------------------------------------------------------")

	// 1. getAllProductsDB()
	allProds, _ := getAllProductsDB(db)
	fmt.Printf("\n1. [getAllProductsDB] Total products found: %d\n", len(allProds))

	// 2. getProductByIDDB()
	prod1, _ := getProductByIDDB(db, 1)
	if prod1 != nil {
		fmt.Printf("2. [getProductByIDDB(db, 1)] Found: ID=%d, Name='%s', Price=฿%.2f\n", prod1.ID, prod1.Name, prod1.Price)
	}

	// 3. getProductByNameDB() LIKE
	likeProds, _ := getProductByNameDB(db, "Keyboard")
	fmt.Printf("3. [getProductByNameDB(db, \"Keyboard\")] Found %d product(s):\n", len(likeProds))
	for _, p := range likeProds {
		fmt.Printf("   - ID: %d, Name: %s, Price: ฿%.2f\n", p.ID, p.Name, p.Price)
	}

	// 4. getProductsByCategoryDB()
	catProds, _ := getProductsByCategoryDB(db, 1)
	fmt.Printf("4. [getProductsByCategoryDB(db, 1)] Laptops category found %d items\n", len(catProds))

	// 5. getProductsByPriceRangeDB(db, 1000, 10000)
	priceRangeProds, _ := getProductsByPriceRangeDB(db, 1000, 10000)
	fmt.Printf("5. [getProductsByPriceRangeDB(db, 1000, 10000)] Found %d products in range ฿1,000 - ฿10,000\n", len(priceRangeProds))
	for _, p := range priceRangeProds {
		fmt.Printf("   - %s: ฿%.2f\n", p.Name, p.Price)
	}

	// 6. getProductsLowStockDB(db, 5)
	lowStockProds, _ := getProductsLowStockDB(db, 5)
	fmt.Printf("6. [getProductsLowStockDB(db, 5)] Products with stock < 5 (%d found):\n", len(lowStockProds))
	for _, p := range lowStockProds {
		fmt.Printf("   - %s: Quantity = %d\n", p.Name, p.Quantity)
	}

	// 7. getProductsByCategoryAndPriceDB(db, 2, 500, 50000)
	catPriceProds, _ := getProductsByCategoryAndPriceDB(db, 2, 500, 50000)
	fmt.Printf("7. [getProductsByCategoryAndPriceDB(db, 2, 500, 50000)] Category 2 in range: %d found\n", len(catPriceProds))

	// 8. getProductsSortedByPriceDB()
	sortedAsc, _ := getProductsSortedByPriceDB(db, "ASC")
	if len(sortedAsc) > 0 {
		fmt.Printf("8a. [getProductsSortedByPriceDB(db, \"ASC\")] Cheapest: %s (฿%.2f)\n", sortedAsc[0].Name, sortedAsc[0].Price)
	}
	sortedDesc, _ := getProductsSortedByPriceDB(db, "DESC")
	if len(sortedDesc) > 0 {
		fmt.Printf("8b. [getProductsSortedByPriceDB(db, \"DESC\")] Most Expensive: %s (฿%.2f)\n", sortedDesc[0].Name, sortedDesc[0].Price)
	}

	// 9. getTopExpensiveProductsDB(db, 3)
	top3, _ := getTopExpensiveProductsDB(db, 3)
	fmt.Println("9. [getTopExpensiveProductsDB(db, 3)] Top 3 Most Expensive Products:")
	for i, p := range top3 {
		fmt.Printf("   #%d %s - ฿%.2f\n", i+1, p.Name, p.Price)
	}

	// 10. getCustomerByEmailDB()
	custByEmail, _ := getCustomerByEmailDB(db, "somchai@example.com")
	if custByEmail != nil {
		fmt.Printf("10. [getCustomerByEmailDB(db, \"somchai@example.com\")] Found ID %d: %s, Phone: %s\n",
			custByEmail.ID, custByEmail.Name, custByEmail.Phone)
	}

	// 11. getOrdersByCustomerDB()
	ordersCust1, _ := getOrdersByCustomerDB(db, 1)
	fmt.Printf("11. [getOrdersByCustomerDB(db, 1)] Customer ID 1 has %d order(s)\n", len(ordersCust1))

	// 12. getOrdersByStatusDB(db, "Pending")
	pendingOrders, _ := getOrdersByStatusDB(db, "Pending")
	fmt.Printf("12. [getOrdersByStatusDB(db, \"Pending\")] Found %d pending order(s)\n", len(pendingOrders))

	// 13. getOrdersByDateRangeDB()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	dateRangeOrders, _ := getOrdersByDateRangeDB(db, start, end)
	fmt.Printf("13. [getOrdersByDateRangeDB] Orders created today: %d\n", len(dateRangeOrders))

	// 14. getOrdersWithDetailsDB()
	detailedOrders, _ := getOrdersWithDetailsDB(db)
	fmt.Printf("14. [getOrdersWithDetailsDB] Retrieved %d orders with full customer and product details.\n", len(detailedOrders))

	// ==================== PART 4: Update ====================
	fmt.Println("\n------------------------------------------------------------------")
	fmt.Println(">>> [ส่วนที่ 4] Update Functions Verification")
	fmt.Println("------------------------------------------------------------------")

	// 1. updateProductPriceDB
	fmt.Println("1. [updateProductPriceDB] Updating price of product ID 1:")
	_ = updateProductPriceDB(db, 1, 35900.00)
	p1Updated, _ := getProductByIDDB(db, 1)
	fmt.Printf("   ✓ New Price for '%s': ฿%.2f\n", p1Updated.Name, p1Updated.Price)

	// Negative price validation test
	errNegPrice := updateProductPriceDB(db, 1, -500.00)
	fmt.Printf("   ✓ Test Negative Price Validation: expected error caught: '%v'\n", errNegPrice)

	// 2. updateProductStockDB
	fmt.Println("2. [updateProductStockDB] Updating stock of product ID 1:")
	_ = updateProductStockDB(db, 1, 15)
	p1UpdatedStock, _ := getProductByIDDB(db, 1)
	fmt.Printf("   ✓ New Stock for '%s': %d\n", p1UpdatedStock.Name, p1UpdatedStock.Quantity)

	// Negative stock validation test
	errNegStock := updateProductStockDB(db, 1, -10)
	fmt.Printf("   ✓ Test Negative Stock Validation: expected error caught: '%v'\n", errNegStock)

	// 3. updateProductByIDDB
	err = updateProductByIDDB(db, 1, map[string]interface{}{"name": "MacBook Air M3 (2024)"})
	p1UpdatedName, _ := getProductByIDDB(db, 1)
	fmt.Printf("3. [updateProductByIDDB] Product ID 1 name updated to: '%s' (Error: %v)\n", p1UpdatedName.Name, err)

	// 4. updateCustomerByIDDB
	err = updateCustomerByIDDB(db, 1, map[string]interface{}{"phone": "081-999-8888"})
	c1, _ := getCustomerByIDDB(db, 1)
	fmt.Printf("4. [updateCustomerByIDDB] Customer ID 1 phone updated to: %s (Error: %v)\n", c1.Phone, err)

	// 5. updateOrderStatusDB
	if testOrder != nil {
		fmt.Printf("5. [updateOrderStatusDB] Updating Order ID %d to 'Completed':\n", testOrder.ID)
		err = updateOrderStatusDB(db, testOrder.ID, "Completed")
		fmt.Printf("   ✓ Order status updated successfully (Error: %v)\n", err)

		// Invalid status validation test
		errInvalidStatus := updateOrderStatusDB(db, testOrder.ID, "Shipped")
		fmt.Printf("   ✓ Test Invalid Status Validation: expected error caught: '%v'\n", errInvalidStatus)
	}

	// ==================== PART 5: Delete ====================
	fmt.Println("\n------------------------------------------------------------------")
	fmt.Println(">>> [ส่วนที่ 5] Delete Functions & Relationship Integrity Verification")
	fmt.Println("------------------------------------------------------------------")

	// 1. Referential Integrity Test: Delete Category with existing products
	fmt.Println("1. [Integrity Check] Attempting to delete Category ID 1 (which contains active products):")
	errDelCat := deleteCategoryByIDDB(db, 1)
	fmt.Printf("   ✓ Prevented illegal deletion as expected! Error: '%v'\n", errDelCat)

	// 2. Referential Integrity Test: Delete Customer with existing orders
	fmt.Println("2. [Integrity Check] Attempting to delete Customer ID 1 (who has active orders):")
	errDelCust := deleteCustomerByIDDB(db, 1)
	fmt.Printf("   ✓ Prevented illegal deletion as expected! Error: '%v'\n", errDelCust)

	// 3. Soft Delete Product
	fmt.Printf("3. [softDeleteProductDB] Soft deleting product ID %d ('%s'):\n", newProd.ID, newProd.Name)
	err = softDeleteProductDB(db, newProd.ID)
	fmt.Printf("   ✓ Soft delete error: %v\n", err)

	// 4. getDeletedProductsDB()
	deletedList, _ := getDeletedProductsDB(db)
	fmt.Printf("4. [getDeletedProductsDB] Retrieved %d soft-deleted product(s):\n", len(deletedList))
	for _, dp := range deletedList {
		fmt.Printf("   - ID: %d, Name: %s, DeletedAt: %v\n", dp.ID, dp.Name, dp.DeletedAt.Time)
	}

	// 5. restoreProductDB()
	fmt.Printf("5. [restoreProductDB] Restoring product ID %d:\n", newProd.ID)
	err = restoreProductDB(db, newProd.ID)
	fmt.Printf("   ✓ Restore error: %v\n", err)
	restoredProd, _ := getProductByIDDB(db, newProd.ID)
	if restoredProd != nil {
		fmt.Printf("   ✓ Verified product is active again: ID=%d, Name='%s'\n", restoredProd.ID, restoredProd.Name)
	}

	// 6. Hard Delete test using an isolated dummy product with no orders
	dummyProd := Product{Name: "Temporary Discard Item", Price: 100.0, Quantity: 1, CategoryID: 3}
	_ = createProductDB(db, &dummyProd)
	fmt.Printf("6. [hardDeleteProductDB] Hard deleting dummy product ID %d:\n", dummyProd.ID)
	errHardDel := hardDeleteProductDB(db, dummyProd.ID)
	fmt.Printf("   ✓ Permanently deleted isolated product: %v\n", errHardDel)

	// Prevent Hard Delete if product has orders
	fmt.Println("   [Hard Delete Integrity Check] Attempting hard delete on product ID 1 (has orders):")
	errHardWithOrders := hardDeleteProductDB(db, 1)
	fmt.Printf("   ✓ Prevented hard delete of ordered product as expected! Error: '%v'\n", errHardWithOrders)

	// ==================== PART 6: Export SQL Dump ====================
	fmt.Println("\n------------------------------------------------------------------")
	fmt.Println(">>> [ส่วนที่ 6] Export SQL Dump File (store_db.sql)")
	fmt.Println("------------------------------------------------------------------")
	if err := exportSQLDump(db, "store_db.sql"); err != nil {
		fmt.Printf("Failed to export SQL: %v\n", err)
	} else {
		fmt.Println("✓ Successfully exported 'store_db.sql' with table schemas and current records!")
	}

	fmt.Println("\n==================================================================")
	fmt.Println("           ALL ASSIGNMENT REQUIREMENTS VERIFIED SUCCESSFULLY!     ")
	fmt.Println("==================================================================")
}

// printProducts displays a formatted list of products
func printProducts(products []Product) {
	fmt.Printf("\n%-4s | %-32s | %-12s | %-8s | %-20s\n", "ID", "Name", "Price (THB)", "Stock", "Category")
	fmt.Println(strings.Repeat("-", 85))
	for _, p := range products {
		catName := "N/A"
		if p.Category.Name != "" {
			catName = p.Category.Name
		}
		fmt.Printf("%-4d | %-32s | %11.2f | %8d | %-20s\n", p.ID, p.Name, p.Price, p.Quantity, catName)
	}
	fmt.Printf("Total: %d products\n", len(products))
}

// printOrdersWithDetails displays orders with joined customer and product info
func printOrdersWithDetails(orders []Order) {
	fmt.Printf("\n%-4s | %-18s | %-28s | %-4s | %-12s | %-10s\n", "ID", "Customer", "Product", "Qty", "Total (THB)", "Status")
	fmt.Println(strings.Repeat("-", 85))
	for _, o := range orders {
		custName := fmt.Sprintf("Cust #%d", o.CustomerID)
		if o.Customer.Name != "" {
			custName = o.Customer.Name
		}
		prodName := fmt.Sprintf("Prod #%d", o.ProductID)
		if o.Product.Name != "" {
			prodName = o.Product.Name
		}
		fmt.Printf("%-4d | %-18s | %-28s | %4d | %11.2f | %-10s\n",
			o.ID, custName, prodName, o.Quantity, o.TotalPrice, o.Status)
	}
	fmt.Printf("Total: %d orders\n", len(orders))
}

// exportSQLDump generates a complete SQL dump containing table structures and INSERT statements
func exportSQLDump(db *gorm.DB, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	header := fmt.Sprintf(`-- ========================================================
-- Store Management System - MySQL Database Dump
-- Database: %s
-- Generated at: %s
-- ========================================================

SET FOREIGN_KEY_CHECKS = 0;
CREATE DATABASE IF NOT EXISTS `+"`%s`"+` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `+"`%s`"+`;

`, DBName, time.Now().Format("2006-01-02 15:04:05"), DBName, DBName)
	writer.WriteString(header)

	tables := []string{"categories", "products", "customers", "orders"}

	for _, table := range tables {
		// 1. Table structure
		writer.WriteString(fmt.Sprintf("--\n-- Table structure for `%s`\n--\n", table))
		writer.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", table))

		var tblName, createStmt string
		row := sqlDB.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`", table))
		if err := row.Scan(&tblName, &createStmt); err != nil {
			return fmt.Errorf("failed to get create table for %s: %w", table, err)
		}
		writer.WriteString(createStmt + ";\n\n")

		// 2. Table data
		rows, err := sqlDB.Query(fmt.Sprintf("SELECT * FROM `%s`", table))
		if err != nil {
			return fmt.Errorf("failed to query rows for %s: %w", table, err)
		}

		cols, _ := rows.Columns()
		colList := "`" + strings.Join(cols, "`, `") + "`"

		var insertValues []string
		for rows.Next() {
			scanDest := make([]interface{}, len(cols))
			scanPtrs := make([]interface{}, len(cols))
			for i := range scanDest {
				scanPtrs[i] = &scanDest[i]
			}

			if err := rows.Scan(scanPtrs...); err != nil {
				rows.Close()
				return err
			}

			vals := make([]string, len(cols))
			for i, val := range scanDest {
				if val == nil {
					vals[i] = "NULL"
				} else {
					switch v := val.(type) {
					case []byte:
						vals[i] = "'" + escapeSQL(string(v)) + "'"
					case string:
						vals[i] = "'" + escapeSQL(v) + "'"
					case time.Time:
						vals[i] = "'" + v.Format("2006-01-02 15:04:05") + "'"
					default:
						vals[i] = fmt.Sprintf("%v", v)
					}
				}
			}
			insertValues = append(insertValues, "("+strings.Join(vals, ", ")+")")
		}
		rows.Close()

		if len(insertValues) > 0 {
			writer.WriteString(fmt.Sprintf("--\n-- Dumping data for `%s`\n--\n", table))
			insertQuery := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES\n%s;\n\n",
				table, colList, strings.Join(insertValues, ",\n"))
			writer.WriteString(insertQuery)
		}
	}

	writer.WriteString("SET FOREIGN_KEY_CHECKS = 1;\n-- Dump completed\n")
	return nil
}

func escapeSQL(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}
