DROP INDEX IF EXISTS idx_transactions_sender;
DROP INDEX IF EXISTS idx_transactions_receiver;
DROP INDEX IF EXISTS idx_purchases_employee;

DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS merch;
DROP TABLE IF EXISTS employees;