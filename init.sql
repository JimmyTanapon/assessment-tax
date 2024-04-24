-- Creation of tax_discount table
CREATE TYPE discount_type AS ENUM('donation', 'personalDeduction', 'k-receipt');

CREATE TABLE IF NOT EXISTS tax_discount(
    id SERIAL PRIMARY KEY,
    discount_type discount_type NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL,
    min_discount_value DECIMAL(10, 2),
    max_discount_value DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tax_discount(discount_type, discount_value, min_discount_value, max_discount_value) VALUES
('donation', 60000.00, 10000.00, 100000.00),
('personalDeduction', 100000.00, 0.00, 100000.00),
('k-receipt', 50000.00, 0.00, 100000.00);
