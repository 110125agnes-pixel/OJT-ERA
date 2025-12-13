-- MySQL initialization script for OJT-ERA Database

-- Create database
CREATE DATABASE IF NOT EXISTS ojt_era_db;
USE ojt_era_db;

-- Employee table
CREATE TABLE IF NOT EXISTS items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    lastname VARCHAR(100) NOT NULL,
    firstname VARCHAR(100) NOT NULL,
    middlename VARCHAR(100),
    suffix VARCHAR(20),
    birthdate DATE,
    sex VARCHAR(20),
    civil_status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Inventory table
CREATE TABLE IF NOT EXISTS inventory (
    id INT AUTO_INCREMENT PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    unit VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Sample data (optional)
INSERT INTO items (lastname, firstname, middlename, suffix, birthdate, sex, civil_status) VALUES 
    ('Doe', 'John', 'A', '', '1990-01-01', 'Male', 'Single'),
    ('Smith', 'Jane', 'B', 'Jr.', '1985-05-15', 'Female', 'Married'),
    ('Lee', 'Chris', '', '', '2000-12-31', 'Other', 'Single');
