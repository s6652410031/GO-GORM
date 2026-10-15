-- ========================================================
-- Store Management System - MySQL Database Dump
-- Database: store_db
-- Generated at: 2026-10-15 20:58:53
-- ========================================================

SET FOREIGN_KEY_CHECKS = 0;
CREATE DATABASE IF NOT EXISTS `store_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `store_db`;

--
-- Table structure for `categories`
--
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `description` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_categories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for `categories`
--
INSERT INTO `categories` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `description`) VALUES
(1, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Laptops & Notebooks', 'High-performance laptops, ultrabooks, and portable workstations'),
(2, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Smartphones & Tablets', 'Latest flagship and mid-range mobile devices'),
(3, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Computer Accessories', 'Mechanical keyboards, ergonomic mice, and USB adapters'),
(4, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Peripherals & Displays', 'High-refresh monitors, portable external SSDs, and displays'),
(5, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 'Audio & Wearables', 'Headphones, Earbuds, Smartwatches'),
(6, '2026-10-15 20:58:53', '2026-10-15 20:58:53', NULL, 'Audio & Wearables', 'Headphones, Earbuds, Smartwatches');

--
-- Table structure for `products`
--
DROP TABLE IF EXISTS `products`;
CREATE TABLE `products` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(150) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `quantity` bigint(20) NOT NULL,
  `category_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`),
  KEY `idx_products_category_id` (`category_id`),
  CONSTRAINT `fk_categories_products` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for `products`
--
INSERT INTO `products` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `price`, `quantity`, `category_id`) VALUES
(1, '2026-10-15 20:58:29', '2026-10-15 20:58:53', NULL, 'MacBook Air M3 (2024)', '35900.00', 15, 1),
(2, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'ASUS Vivobook 15', '21900.00', 7, 1),
(3, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Dell XPS 13', '45900.00', 4, 1),
(4, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Lenovo ThinkPad X1', '38900.00', 5, 1),
(5, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'iPhone 16 128GB', '32900.00', 14, 2),
(6, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'iPhone 16 Pro 256GB', '41900.00', 3, 2),
(7, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'Samsung Galaxy S24', '29900.00', 9, 2),
(8, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Samsung Galaxy Z Flip 6', '36900.00', 4, 2),
(9, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'Logitech Mouse MX Master 3S', '3590.00', 19, 3),
(10, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'Logitech G Pro Wireless Mouse', '4290.00', 7, 3),
(11, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'Mechanical Keyboard Keychron K2', '3890.00', 13, 3),
(12, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'USB-C Hub 7-in-1 Adapter', '1290.00', 22, 3),
(13, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'HyperX Cloud III Gaming Headset', '2890.00', 8, 3),
(14, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'LG UltraFine 27-inch 4K Monitor', '16900.00', 4, 4),
(15, '2026-10-15 20:58:29', '2026-10-15 20:58:30', NULL, 'Samsung T7 Shield 1TB Portable SSD', '4500.00', 16, 4),
(16, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 'Sony WH-1000XM5', '11900.00', 7, 5),
(18, '2026-10-15 20:58:53', '2026-10-15 20:58:53', NULL, 'Sony WH-1000XM5', '11900.00', 7, 6);

--
-- Table structure for `customers`
--
DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_customers_email` (`email`),
  KEY `idx_customers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for `customers`
--
INSERT INTO `customers` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `email`, `phone`) VALUES
(1, '2026-10-15 20:58:29', '2026-10-15 20:58:53', NULL, 'Somchai Jaidee', 'somchai@example.com', '081-999-8888'),
(2, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Suda Rakdee', 'suda@example.com', '082-222-3333'),
(3, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Arthit Srisuk', 'arthit@example.com', '083-333-4444'),
(4, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Wanna Suksan', 'wanna@example.com', '084-444-5555'),
(5, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Natthapong Wongsuwan', 'natthapong@example.com', '085-555-6666'),
(6, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Pornthip Charoen', 'pornthip@example.com', '086-666-7777'),
(7, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Komsan Thongdee', 'komsan@example.com', '087-777-8888'),
(8, '2026-10-15 20:58:29', '2026-10-15 20:58:29', NULL, 'Siriporn Meesuk', 'siriporn@example.com', '088-888-9999'),
(9, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 'Chaiwat Demo', 'chaiwat_1792072710@example.com', '089-999-9999'),
(10, '2026-10-15 20:58:53', '2026-10-15 20:58:53', NULL, 'Chaiwat Demo', 'chaiwat_1792072733@example.com', '089-999-9999');

--
-- Table structure for `orders`
--
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `customer_id` bigint(20) unsigned NOT NULL,
  `product_id` bigint(20) unsigned NOT NULL,
  `quantity` bigint(20) NOT NULL,
  `total_price` decimal(12,2) NOT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'Pending',
  PRIMARY KEY (`id`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `idx_orders_customer_id` (`customer_id`),
  KEY `idx_orders_product_id` (`product_id`),
  CONSTRAINT `fk_customers_orders` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_products_orders` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for `orders`
--
INSERT INTO `orders` (`id`, `created_at`, `updated_at`, `deleted_at`, `customer_id`, `product_id`, `quantity`, `total_price`, `status`) VALUES
(1, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 1, 1, 1, '34900.00', 'Completed'),
(2, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 1, 9, 1, '3590.00', 'Completed'),
(3, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 2, 5, 1, '32900.00', 'Completed'),
(4, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 2, 12, 2, '2580.00', 'Pending'),
(5, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 3, 2, 1, '21900.00', 'Completed'),
(6, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 3, 11, 1, '3890.00', 'Completed'),
(7, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 4, 7, 1, '29900.00', 'Pending'),
(8, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 5, 14, 1, '16900.00', 'Completed'),
(9, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 5, 15, 2, '9000.00', 'Pending'),
(10, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 6, 10, 1, '4290.00', 'Cancelled'),
(11, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 7, 12, 1, '1290.00', 'Completed'),
(12, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 8, 13, 1, '2890.00', 'Pending'),
(13, '2026-10-15 20:58:30', '2026-10-15 20:58:30', NULL, 9, 16, 3, '35700.00', 'Completed'),
(14, '2026-10-15 20:58:53', '2026-10-15 20:58:53', NULL, 10, 18, 3, '35700.00', 'Completed');

SET FOREIGN_KEY_CHECKS = 1;
-- Dump completed
