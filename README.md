# ระบบจัดการร้านค้า (Store Management System)
**Go + GORM + MySQL Practical Assignment**

โปรเจกต์นี้พัฒนาด้วยภาษา **Go** ร่วมกับ **GORM (ORM for Go)** และฐานข้อมูล **MySQL** ตามข้อกำหนดของใบงานปฏิบัติการทุกประการ

---

## 📁 โครงสร้างโปรเจกต์ (Project Structure)

```text
GO-GORM/
├── go.mod                  # Go Module definition
├── go.sum                  # Checksums ของ dependencies
├── models.go               # Struct ทั้ง 4 ตาราง (Category, Product, Customer, Order) และ InitDB()
├── category.go             # ฟังก์ชันจัดการ Category (Create, Delete พร้อมตรวจสอบเงื่อนไข)
├── product.go              # ฟังก์ชันจัดการ Product (Create, Read หลากหลายเงื่อนไข, Update, Soft/Hard Delete, Restore)
├── customer.go             # ฟังก์ชันจัดการ Customer (Create, Read, Update, Delete พร้อมตรวจสอบประวัติการสั่งซื้อ)
├── order.go                # ฟังก์ชันจัดการ Order (Create Transaction พร้อมตัดสต็อก, Read, Update สถานะ)
├── main.go                 # หน้าเมนู Console, ตัวรัน Demonstration ทดสอบทุกฟังก์ชัน, ตัวสร้าง SQL Dump
└── store_db.sql            # ไฟล์ SQL Export (โครงสร้างตารางและข้อมูลตั้งต้น)
```

---

## ⚙️ การตั้งค่าฐานข้อมูล (Database Configuration)
โปรแกรมถูกตั้งค่าให้เชื่อมต่อไปยัง:
- **Host**: `127.0.0.1` (localhost)
- **Port**: `3306`
- **Username**: `root`
- **Password**: *(ว่าง)*
- **Database**: `store_db` *(โปรแกรมจะสร้างฐานข้อมูลให้อัตโนมัติหากยังไม่มี)*

---

## 🚀 วิธีการรันโปรแกรม (How to Run)

ทุกไฟล์ใช้ `package main` สามารถรันผ่านคำสั่ง:

```bash
go run .
```

หรือรันโหมดทดสอบอัตโนมัติ (Automated Verification of All Parts 1–5):

```bash
go run . --demo
```

---

## 📋 สรุปฟังก์ชันที่พัฒนาตามใบงาน

### ส่วนที่ 1: Database & Models (`models.go`)
- ทุก Model สืบทอด `gorm.Model` (`ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`)
- กำหนด Primary Key และ Foreign Key:
  - `Category` 1 : N `Product`
  - `Customer` 1 : N `Order`
  - `Product` 1 : N `Order`
- ฟังก์ชัน `InitDB()` ตรวจสอบและสร้างฐานข้อมูล `store_db` อัตโนมัติ พร้อมรัน `db.AutoMigrate()`

### ส่วนที่ 2: Create
- `createCategoryDB()`: เพิ่มประเภทสินค้า
- `createProductDB()`: เพิ่มสินค้าพร้อมประเภท
- `createCustomerDB()`: เพิ่มลูกค้า
- `createOrderDB()`: เพิ่มคำสั่งซื้อโดยใช้ **Transaction (`db.Transaction`)** ตรวจสอบลูกค้า สินค้า เช็คสต็อก คำนวณ `TotalPrice = Price × Quantity` และตัดสต็อกสินค้า หากขั้นตอนใดล้มเหลวจะ Rollback ทันที

### ส่วนที่ 3: Read / Conditional Queries
- `getAllProductsDB()`: ดึงสินค้าทั้งหมด
- `getProductByIDDB(db, id)`: ค้นหาสินค้าจาก ID
- `getProductByNameDB(db, name)`: ค้นหาสินค้าด้วย LIKE (`%name%`)
- `getProductsByCategoryDB(db, categoryID)`: ค้นหาสินค้าตามประเภท
- `getProductsByPriceRangeDB(db, min, max)`: ค้นหาสินค้าตามช่วงราคา
- `getProductsLowStockDB(db, threshold)`: ค้นหาสินค้าที่สต็อกต่ำกว่าค่าที่กำหนด
- `getProductsByCategoryAndPriceDB(db, categoryID, min, max)`: ค้นหาสินค้าตามประเภทและช่วงราคา
- `getProductsSortedByPriceDB(db, order)`: เรียงลำดับราคาสินค้า ASC / DESC
- `getTopExpensiveProductsDB(db, limit)`: หาสินค้าที่แพงที่สุด N รายการ
- `getCustomerByEmailDB(db, email)`: ค้นหาลูกค้าด้วย Email
- `getOrdersByCustomerDB(db, customerID)`: ค้นหาประวัติคำสั่งซื้อของลูกค้า
- `getOrdersByStatusDB(db, status)`: ค้นหาคำสั่งซื้อตามสถานะ
- `getOrdersByDateRangeDB(db, start, end)`: ค้นหาคำสั่งซื้อตามช่วงวันที่
- `getOrdersWithDetailsDB(db)`: ดึงคำสั่งซื้อพร้อม Preload ข้อมูลลูกค้าและสินค้า

### ส่วนที่ 4: Update
- `updateProductByIDDB()`: แก้ไขข้อมูลสินค้า
- `updateProductPriceDB()`: แก้ไขราคา (ตรวจสอบไม่ให้ติดลบ)
- `updateProductStockDB()`: เพิ่ม/ลดสต็อกสินค้า (ตรวจสอบไม่ให้ติดลบ)
- `updateCustomerByIDDB()`: แก้ไขข้อมูลลูกค้า
- `updateOrderStatusDB()`: เปลี่ยนสถานะคำสั่งซื้อ (ตรวจสอบให้เป็น `Pending`, `Completed` หรือ `Cancelled` เท่านั้น หากยกเลิกจะคืนสต็อกสินค้าให้โดยอัตโนมัติ)

### ส่วนที่ 5: Delete
- `softDeleteProductDB()`: Soft Delete สินค้า (ยังคงข้อมูลในฐานข้อมูลแต่ใส่ `deleted_at`)
- `hardDeleteProductDB()`: ลบสินค้าถาวร (ไม่อนุญาตให้ลบสินค้าที่มีในประวัติคำสั่งซื้อ)
- `deleteCustomerByIDDB()`: ลบลูกค้า (ไม่อนุญาตให้ลบหากลูกค้ามีประวัติคำสั่งซื้อ)
- `deleteCategoryByIDDB()`: ลบประเภทสินค้า (ไม่อนุญาตให้ลบหากยังมีสินค้าอยู่ในประเภทนั้น)
- `getDeletedProductsDB()`: ดึงรายการสินค้าที่ถูก Soft Delete
- `restoreProductDB()`: กู้คืนสินค้าที่ถูก Soft Delete กลับมาใช้งานใหม่

### ส่วนที่ 6: ข้อมูลทดสอบและไฟล์ส่งมอบ
- **Seed Data ขั้นต่ำครบถ้วน**:
  - Category: 4 รายการ
  - Product: 15 รายการ (MacBook Air, ASUS Vivobook, iPhone 16, Samsung Galaxy, Logitech Mouse, Mechanical Keyboard, USB-C Hub ฯลฯ)
  - Customer: 8 รายการ
  - Order: 12 รายการ
- **ไฟล์ SQL Export**: `store_db.sql` ถูกสร้างขึ้นโดยอัตโนมัติ พร้อมโครงสร้างตาราง (DDL) และข้อมูลตั้งต้น (DML) สำหรับนำไปส่งได้ทันที
