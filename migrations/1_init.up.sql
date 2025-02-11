CREATE TABLE employees (
                           id SERIAL PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           balance INT NOT NULL DEFAULT 1000,
                           password VARCHAR(255) NOT NULL
);

CREATE TABLE merch (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) UNIQUE NOT NULL,
                       price INT NOT NULL
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              sender_id INT NOT NULL,
                              receiver_id INT NOT NULL,
                              amount INT NOT NULL,
                              timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              CONSTRAINT chk_positive_amount CHECK (amount > 0),
                              CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES employees (id) ON DELETE CASCADE,
                              CONSTRAINT fk_receiver FOREIGN KEY (receiver_id) REFERENCES employees (id) ON DELETE CASCADE
);

CREATE TABLE purchases (
                           id SERIAL PRIMARY KEY,
                           employee_id INT NOT NULL,
                           merch_id INT NOT NULL,
                           timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           CONSTRAINT fk_employee FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE,
                           CONSTRAINT fk_merch FOREIGN KEY (merch_id) REFERENCES merch (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_employees_name ON employees(name);
CREATE INDEX idx_transactions_sender ON transactions(sender_id);
CREATE INDEX idx_transactions_receiver ON transactions(receiver_id);
CREATE INDEX idx_purchases_employee ON purchases(employee_id);
