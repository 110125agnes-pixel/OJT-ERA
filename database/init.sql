-- SQLite initialization script
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Sample data (optional)
INSERT INTO items (name) VALUES 
    ('Sample Item 1'),
    ('Sample Item 2'),
    ('Sample Item 3');
