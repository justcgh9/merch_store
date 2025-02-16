DROP INDEX IF EXISTS idx_history_from_user;
DROP INDEX IF EXISTS idx_history_to_user;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_balance_username;
DROP INDEX IF EXISTS idx_inventory_username;

DROP TABLE IF EXISTS History;
DROP TABLE IF EXISTS Inventory;
DROP TABLE IF EXISTS Balance;
DROP TABLE IF EXISTS Users;
