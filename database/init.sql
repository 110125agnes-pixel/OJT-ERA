-- SQLite initialization script

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

-- Sample data (optional)

INSERT INTO items (lastname, firstname, middlename, suffix, birthdate, sex, civil_status) VALUES 
    ('Doe', 'John', 'A', '', '1990-01-01', 'Male', 'Single'),
    ('Smith', 'Jane', 'B', 'Jr.', '1985-05-15', 'Female', 'Married'),
    ('Lee', 'Chris', '', '', '2000-12-31', 'Other', 'Single');
