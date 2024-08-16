-- Wallets table --
CREATE TABLE IF NOT EXISTS wallets (
    wallet_id TEXT PRIMARY KEY,
    balance REAL,
    currency VARCHAR
);
CREATE Unique INDEX IF NOT EXISTS wallet_id_currency_index ON wallets(wallet_id,currency);
-- Wallets table --