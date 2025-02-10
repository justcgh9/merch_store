CREATE TABLE IF NOT EXISTS Users (
    username VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_username ON Users(username);

CREATE TABLE IF NOT EXISTS Balance (
    username VARCHAR(255) PRIMARY KEY REFERENCES Users(username) ON DELETE CASCADE,
    balance INTEGER NOT NULL CHECK (balance >= 0)
);

CREATE INDEX IF NOT EXISTS idx_balance_username ON Balance(username);

CREATE TABLE IF NOT EXISTS Inventory (
    username VARCHAR(255) PRIMARY KEY REFERENCES Users(username) ON DELETE CASCADE,
    t_shirt INTEGER DEFAULT 0 CHECK (t_shirt >= 0),
    cup INTEGER DEFAULT 0 CHECK (cup >= 0),
    book INTEGER DEFAULT 0 CHECK (book >= 0),
    pen INTEGER DEFAULT 0 CHECK (pen >= 0),
    powerbank INTEGER DEFAULT 0 CHECK (powerbank >= 0),
    hoody INTEGER DEFAULT 0 CHECK (hoody >= 0),
    umbrella INTEGER DEFAULT 0 CHECK (umbrella >= 0),
    socks INTEGER DEFAULT 0 CHECK (socks >= 0),
    wallet INTEGER DEFAULT 0 CHECK (wallet >= 0),
    pink_hoody INTEGER DEFAULT 0 CHECK (pink_hoody >= 0)
);

CREATE INDEX IF NOT EXISTS idx_inventory_username ON Inventory(username);

CREATE TABLE IF NOT EXISTS History (
    id SERIAL PRIMARY KEY,
    from_user VARCHAR(255) REFERENCES Users(username) ON DELETE CASCADE,
    to_user VARCHAR(255) REFERENCES Users(username) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_history_from_user ON History(from_user);
CREATE INDEX IF NOT EXISTS idx_history_to_user ON History(to_user);
CREATE INDEX IF NOT EXISTS idx_history_to_user ON History(from_user, to_user);
