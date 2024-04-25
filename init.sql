-- Creation of tax_discount table
CREATE TYPE discount_type AS ENUM('donation', 'personal', 'k-receipt');

CREATE TABLE IF NOT EXISTS tax_discount(
    id SERIAL PRIMARY KEY,
    discount_name VARCHAR(250) UNIQUE NOT NULL,
    discount_type discount_type NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL,
    min_discount_value DECIMAL(10, 2),
    max_discount_value DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tax_discount(discount_name,discount_type, discount_value, min_discount_value, max_discount_value) VALUES
('donation','donation', 100000.00, 0.00, 100000.00),
('personalDeduction','personal', 60000.00, 10000.00, 100000.00),
('kreceipt','k-receipt', 50000.00, 0.00, 100000.00);

CREATE TABLE IF NOT EXISTS  users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL  ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    
);


INSERT INTO users(username,password,role) VALUES
('adminTax', 'admin!', 'admin');