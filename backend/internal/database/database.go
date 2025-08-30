package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Connect establishes a connection to PostgreSQL database
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	logrus.Info("Successfully connected to PostgreSQL database")
	return db, nil
}

// Migrate runs database migrations
func Migrate(db *sql.DB) error {
	logrus.Info("Running database migrations...")

	migrations := []string{
		createBlocksTable,
		createTransactionsTable,
		createAddressesTable,
		createTokensTable,
		createEventsTable,
		createIndexes,
	}

	for i, migration := range migrations {
		logrus.Debugf("Running migration %d/%d", i+1, len(migrations))
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", i+1, err)
		}
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}

const createBlocksTable = `
CREATE TABLE IF NOT EXISTS blocks (
    number BIGINT PRIMARY KEY,
    hash VARCHAR(66) UNIQUE NOT NULL,
    parent_hash VARCHAR(66) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    gas_limit BIGINT NOT NULL,
    gas_used BIGINT NOT NULL,
    difficulty NUMERIC(78, 0),
    total_difficulty NUMERIC(78, 0),
    size INTEGER NOT NULL,
    transaction_count INTEGER NOT NULL,
    miner VARCHAR(42) NOT NULL,
    extra_data TEXT,
    base_fee_per_gas BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createTransactionsTable = `
CREATE TABLE IF NOT EXISTS transactions (
    hash VARCHAR(66) PRIMARY KEY,
    block_number BIGINT NOT NULL REFERENCES blocks(number),
    transaction_index INTEGER NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    value NUMERIC(78, 0) NOT NULL,
    gas_limit BIGINT NOT NULL,
    gas_used BIGINT,
    gas_price BIGINT,
    max_fee_per_gas BIGINT,
    max_priority_fee_per_gas BIGINT,
    nonce BIGINT NOT NULL,
    input_data TEXT,
    status INTEGER,
    contract_address VARCHAR(42),
    logs_bloom TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createAddressesTable = `
CREATE TABLE IF NOT EXISTS addresses (
    address VARCHAR(42) PRIMARY KEY,
    balance NUMERIC(78, 0) DEFAULT 0,
    nonce BIGINT DEFAULT 0,
    is_contract BOOLEAN DEFAULT FALSE,
    contract_creator VARCHAR(42),
    creation_transaction VARCHAR(66),
    first_seen_block BIGINT,
    last_seen_block BIGINT,
    transaction_count BIGINT DEFAULT 0,
    label VARCHAR(255),
    tags TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createTokensTable = `
CREATE TABLE IF NOT EXISTS tokens (
    address VARCHAR(42) PRIMARY KEY,
    name VARCHAR(255),
    symbol VARCHAR(50),
    decimals INTEGER,
    total_supply NUMERIC(78, 0),
    token_type VARCHAR(20) DEFAULT 'ERC20',
    creator VARCHAR(42),
    creation_block BIGINT,
    creation_transaction VARCHAR(66),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createEventsTable = `
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    transaction_hash VARCHAR(66) NOT NULL REFERENCES transactions(hash),
    block_number BIGINT NOT NULL,
    log_index INTEGER NOT NULL,
    address VARCHAR(42) NOT NULL,
    topic0 VARCHAR(66),
    topic1 VARCHAR(66),
    topic2 VARCHAR(66),
    topic3 VARCHAR(66),
    data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(transaction_hash, log_index)
);`

const createIndexes = `
-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_blocks_timestamp ON blocks(timestamp);
CREATE INDEX IF NOT EXISTS idx_blocks_miner ON blocks(miner);
CREATE INDEX IF NOT EXISTS idx_transactions_block_number ON transactions(block_number);
CREATE INDEX IF NOT EXISTS idx_transactions_from_address ON transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_transactions_to_address ON transactions(to_address);
CREATE INDEX IF NOT EXISTS idx_addresses_balance ON addresses(balance);
CREATE INDEX IF NOT EXISTS idx_addresses_is_contract ON addresses(is_contract);
CREATE INDEX IF NOT EXISTS idx_events_address ON events(address);
CREATE INDEX IF NOT EXISTS idx_events_topic0 ON events(topic0);
CREATE INDEX IF NOT EXISTS idx_events_block_number ON events(block_number);

-- Full text search indexes
CREATE INDEX IF NOT EXISTS idx_addresses_label_fts ON addresses USING gin(to_tsvector('english', COALESCE(label, '')));
CREATE INDEX IF NOT EXISTS idx_tokens_name_symbol_fts ON tokens USING gin(to_tsvector('english', COALESCE(name, '') || ' ' || COALESCE(symbol, '')));
`
