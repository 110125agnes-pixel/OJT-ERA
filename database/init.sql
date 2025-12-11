-- SQLite initialization script

-- Employee table
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lastname TEXT NOT NULL,
    firstname TEXT NOT NULL,
    middlename TEXT,
    suffix TEXT,
    birthdate TEXT,
    sex TEXT,
    civil_status TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Inventory table
CREATE TABLE IF NOT EXISTS inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_name TEXT NOT NULL,
    category TEXT NOT NULL,
    brand TEXT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    unit TEXT NOT NULL,
    price REAL NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Sample data (optional)

INSERT INTO items (lastname, firstname, middlename, suffix, birthdate, sex, civil_status) VALUES 
    ('Doe', 'John', 'A', '', '1990-01-01', 'Male', 'Single'),
    ('Smith', 'Jane', 'B', 'Jr.', '1985-05-15', 'Female', 'Married'),
    ('Lee', 'Chris', '', '', '2000-12-31', 'Other', 'Single');
