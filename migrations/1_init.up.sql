CREATE TABLE employees (
                           id SERIAL PRIMARY KEY,
                           name VARCHAR(255) UNIQUE NOT NULL,
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
                              sender VARCHAR(255) NOT NULL,
                              receiver VARCHAR(255) NOT NULL,
                              amount INT NOT NULL,
                              CONSTRAINT chk_positive_amount CHECK (amount > 0),
                              CONSTRAINT fk_sender FOREIGN KEY (sender) REFERENCES employees (name) ON DELETE CASCADE,
                              CONSTRAINT fk_receiver FOREIGN KEY (receiver) REFERENCES employees (name) ON DELETE CASCADE
);

CREATE TABLE purchases (
                           id SERIAL PRIMARY KEY,
                           employee_id INT NOT NULL,
                           merch_id INT NOT NULL,
                           quantity INT NOT NULL,
                           CONSTRAINT fk_employee FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE,
                           CONSTRAINT fk_merch FOREIGN KEY (merch_id) REFERENCES merch (id) ON DELETE CASCADE,
                           CONSTRAINT uniq_employee_merch UNIQUE (employee_id, merch_id)
);

CREATE UNIQUE INDEX idx_employees_name ON employees(name);
CREATE INDEX idx_transactions_sender ON transactions(sender);
CREATE INDEX idx_transactions_receiver ON transactions(receiver);
CREATE INDEX idx_purchases_employee ON purchases(employee_id);

INSERT INTO merch (name, price) VALUES
                                    ('t-shirt', 80),
                                    ('cup', 20),
                                    ('book', 50),
                                    ('pen', 10),
                                    ('powerbank', 200),
                                    ('hoody', 300),
                                    ('umbrella', 200),
                                    ('socks', 10),
                                    ('wallet', 50),
                                    ('pink-hoody', 500);
