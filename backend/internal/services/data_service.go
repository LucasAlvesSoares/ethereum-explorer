package services

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// Block represents a blockchain block
type Block struct {
	Number           int64     `json:"number" db:"number"`
	Hash             string    `json:"hash" db:"hash"`
	ParentHash       string    `json:"parent_hash" db:"parent_hash"`
	Timestamp        time.Time `json:"timestamp" db:"timestamp"`
	GasLimit         int64     `json:"gas_limit" db:"gas_limit"`
	GasUsed          int64     `json:"gas_used" db:"gas_used"`
	Difficulty       string    `json:"difficulty" db:"difficulty"`
	TotalDifficulty  string    `json:"total_difficulty" db:"total_difficulty"`
	Size             int       `json:"size" db:"size"`
	TransactionCount int       `json:"transaction_count" db:"transaction_count"`
	Miner            string    `json:"miner" db:"miner"`
	ExtraData        string    `json:"extra_data" db:"extra_data"`
	BaseFeePerGas    *int64    `json:"base_fee_per_gas" db:"base_fee_per_gas"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash                 string    `json:"hash" db:"hash"`
	BlockNumber          int64     `json:"block_number" db:"block_number"`
	TransactionIndex     int       `json:"transaction_index" db:"transaction_index"`
	FromAddress          string    `json:"from_address" db:"from_address"`
	ToAddress            *string   `json:"to_address" db:"to_address"`
	Value                string    `json:"value" db:"value"`
	GasLimit             int64     `json:"gas_limit" db:"gas_limit"`
	GasUsed              *int64    `json:"gas_used" db:"gas_used"`
	GasPrice             *int64    `json:"gas_price" db:"gas_price"`
	MaxFeePerGas         *int64    `json:"max_fee_per_gas" db:"max_fee_per_gas"`
	MaxPriorityFeePerGas *int64    `json:"max_priority_fee_per_gas" db:"max_priority_fee_per_gas"`
	Nonce                int64     `json:"nonce" db:"nonce"`
	InputData            string    `json:"input_data" db:"input_data"`
	Status               *int      `json:"status" db:"status"`
	ContractAddress      *string   `json:"contract_address" db:"contract_address"`
	LogsBloom            string    `json:"logs_bloom" db:"logs_bloom"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// Address represents an Ethereum address
type Address struct {
	Address             string    `json:"address" db:"address"`
	Balance             string    `json:"balance" db:"balance"`
	Nonce               int64     `json:"nonce" db:"nonce"`
	IsContract          bool      `json:"is_contract" db:"is_contract"`
	ContractCreator     *string   `json:"contract_creator" db:"contract_creator"`
	CreationTransaction *string   `json:"creation_transaction" db:"creation_transaction"`
	FirstSeenBlock      *int64    `json:"first_seen_block" db:"first_seen_block"`
	LastSeenBlock       *int64    `json:"last_seen_block" db:"last_seen_block"`
	TransactionCount    int64     `json:"transaction_count" db:"transaction_count"`
	Label               *string   `json:"label" db:"label"`
	Tags                []string  `json:"tags" db:"tags"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// GasPrice represents gas price data
type GasPrice struct {
	ID                 int64     `json:"id" db:"id"`
	BlockNumber        int64     `json:"block_number" db:"block_number"`
	Timestamp          time.Time `json:"timestamp" db:"timestamp"`
	BaseFeePerGas      *int64    `json:"base_fee_per_gas" db:"base_fee_per_gas"`
	SlowGasPrice       int64     `json:"slow_gas_price" db:"slow_gas_price"`
	StandardGasPrice   int64     `json:"standard_gas_price" db:"standard_gas_price"`
	FastGasPrice       int64     `json:"fast_gas_price" db:"fast_gas_price"`
	SlowWaitTime       int       `json:"slow_wait_time" db:"slow_wait_time"`
	StandardWaitTime   int       `json:"standard_wait_time" db:"standard_wait_time"`
	FastWaitTime       int       `json:"fast_wait_time" db:"fast_wait_time"`
	NetworkUtilization *float64  `json:"network_utilization" db:"network_utilization"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// BlocksResponse represents the response for blocks list
type BlocksResponse struct {
	Blocks     []Block `json:"blocks"`
	TotalCount int64   `json:"total_count"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	Mode       string  `json:"mode"`
}

// TransactionsResponse represents the response for transactions list
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	TotalCount   int64         `json:"total_count"`
	Page         int           `json:"page"`
	Limit        int           `json:"limit"`
	Mode         string        `json:"mode"`
}

// AddressResponse represents the response for address details
type AddressResponse struct {
	Address      Address       `json:"address"`
	Transactions []Transaction `json:"transactions,omitempty"`
	Mode         string        `json:"mode"`
}

// DataService defines the interface for data access
type DataService interface {
	GetBlocks(page, limit int) (*BlocksResponse, error)
	GetBlock(identifier string) (*Block, error)
	GetTransactions(page, limit int) (*TransactionsResponse, error)
	GetTransaction(hash string) (*Transaction, error)
	GetAddress(address string) (*AddressResponse, error)
	GetAddressTransactions(address string, page, limit int) (*TransactionsResponse, error)
	SearchByQuery(query string) (*SearchResult, error)
	GetGasPrices(hours int) ([]GasPrice, error)
	GetMode() string
}

// SearchResult represents search results
type SearchResult struct {
	Block       *Block       `json:"block,omitempty"`
	Transaction *Transaction `json:"transaction,omitempty"`
	Address     *Address     `json:"address,omitempty"`
	Mode        string       `json:"mode"`
}

// LiveDataService implements DataService using database and Ethereum client
type LiveDataService struct {
	db   *sql.DB
	mode string
}

// NewLiveDataService creates a new live data service
func NewLiveDataService(db *sql.DB) DataService {
	return &LiveDataService{
		db:   db,
		mode: "live",
	}
}

// GetMode returns the current mode
func (l *LiveDataService) GetMode() string {
	return l.mode
}

// Implement LiveDataService methods
func (l *LiveDataService) GetBlocks(page, limit int) (*BlocksResponse, error) {
	offset := (page - 1) * limit

	query := `
		SELECT number, hash, parent_hash, timestamp, gas_limit, gas_used, 
			   difficulty, total_difficulty, size, transaction_count, miner, 
			   extra_data, base_fee_per_gas, created_at, updated_at
		FROM blocks 
		ORDER BY number DESC 
		LIMIT $1 OFFSET $2`

	rows, err := l.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query blocks: %w", err)
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var block Block
		err := rows.Scan(
			&block.Number, &block.Hash, &block.ParentHash, &block.Timestamp,
			&block.GasLimit, &block.GasUsed, &block.Difficulty, &block.TotalDifficulty,
			&block.Size, &block.TransactionCount, &block.Miner, &block.ExtraData,
			&block.BaseFeePerGas, &block.CreatedAt, &block.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan block: %w", err)
		}
		blocks = append(blocks, block)
	}

	// Get total count
	var totalCount int64
	err = l.db.QueryRow("SELECT COUNT(*) FROM blocks").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &BlocksResponse{
		Blocks:     blocks,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		Mode:       l.mode,
	}, nil
}

func (l *LiveDataService) GetBlock(identifier string) (*Block, error) {
	var query string
	var args []interface{}

	// Check if identifier is numeric (block number) or hash
	if blockNum, err := strconv.ParseInt(identifier, 10, 64); err == nil {
		query = `
			SELECT number, hash, parent_hash, timestamp, gas_limit, gas_used,
				   difficulty, total_difficulty, size, transaction_count, miner,
				   extra_data, base_fee_per_gas, created_at, updated_at
			FROM blocks WHERE number = $1`
		args = []interface{}{blockNum}
	} else {
		query = `
			SELECT number, hash, parent_hash, timestamp, gas_limit, gas_used,
				   difficulty, total_difficulty, size, transaction_count, miner,
				   extra_data, base_fee_per_gas, created_at, updated_at
			FROM blocks WHERE hash = $1`
		args = []interface{}{identifier}
	}

	var block Block
	err := l.db.QueryRow(query, args...).Scan(
		&block.Number, &block.Hash, &block.ParentHash, &block.Timestamp,
		&block.GasLimit, &block.GasUsed, &block.Difficulty, &block.TotalDifficulty,
		&block.Size, &block.TransactionCount, &block.Miner, &block.ExtraData,
		&block.BaseFeePerGas, &block.CreatedAt, &block.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("block not found")
		}
		return nil, fmt.Errorf("failed to query block: %w", err)
	}

	return &block, nil
}

func (l *LiveDataService) GetTransactions(page, limit int) (*TransactionsResponse, error) {
	offset := (page - 1) * limit

	query := `
		SELECT hash, block_number, transaction_index, from_address, to_address,
			   value, gas_limit, gas_used, gas_price, max_fee_per_gas,
			   max_priority_fee_per_gas, nonce, input_data, status,
			   contract_address, logs_bloom, created_at, updated_at
		FROM transactions 
		ORDER BY block_number DESC, transaction_index DESC
		LIMIT $1 OFFSET $2`

	rows, err := l.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(
			&tx.Hash, &tx.BlockNumber, &tx.TransactionIndex, &tx.FromAddress,
			&tx.ToAddress, &tx.Value, &tx.GasLimit, &tx.GasUsed, &tx.GasPrice,
			&tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas, &tx.Nonce, &tx.InputData,
			&tx.Status, &tx.ContractAddress, &tx.LogsBloom, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	// Get total count
	var totalCount int64
	err = l.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &TransactionsResponse{
		Transactions: transactions,
		TotalCount:   totalCount,
		Page:         page,
		Limit:        limit,
		Mode:         l.mode,
	}, nil
}

func (l *LiveDataService) GetTransaction(hash string) (*Transaction, error) {
	query := `
		SELECT hash, block_number, transaction_index, from_address, to_address,
			   value, gas_limit, gas_used, gas_price, max_fee_per_gas,
			   max_priority_fee_per_gas, nonce, input_data, status,
			   contract_address, logs_bloom, created_at, updated_at
		FROM transactions WHERE hash = $1`

	var tx Transaction
	err := l.db.QueryRow(query, hash).Scan(
		&tx.Hash, &tx.BlockNumber, &tx.TransactionIndex, &tx.FromAddress,
		&tx.ToAddress, &tx.Value, &tx.GasLimit, &tx.GasUsed, &tx.GasPrice,
		&tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas, &tx.Nonce, &tx.InputData,
		&tx.Status, &tx.ContractAddress, &tx.LogsBloom, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to query transaction: %w", err)
	}

	return &tx, nil
}

func (l *LiveDataService) GetAddress(address string) (*AddressResponse, error) {
	query := `
		SELECT address, balance, nonce, is_contract, contract_creator,
			   creation_transaction, first_seen_block, last_seen_block,
			   transaction_count, label, tags, created_at, updated_at
		FROM addresses WHERE address = $1`

	var addr Address
	var tags []byte
	err := l.db.QueryRow(query, address).Scan(
		&addr.Address, &addr.Balance, &addr.Nonce, &addr.IsContract,
		&addr.ContractCreator, &addr.CreationTransaction, &addr.FirstSeenBlock,
		&addr.LastSeenBlock, &addr.TransactionCount, &addr.Label, &tags,
		&addr.CreatedAt, &addr.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query address: %w", err)
	}

	if err == sql.ErrNoRows {
		// Address not found in database, create minimal response
		addr = Address{
			Address:          address,
			Balance:          "0",
			Nonce:            0,
			IsContract:       false,
			TransactionCount: 0,
			Tags:             []string{},
		}
	} else {
		// Parse tags array
		if len(tags) > 0 {
			// Handle PostgreSQL array format
			addr.Tags = []string{}
		}
	}

	return &AddressResponse{
		Address: addr,
		Mode:    l.mode,
	}, nil
}

func (l *LiveDataService) GetAddressTransactions(address string, page, limit int) (*TransactionsResponse, error) {
	offset := (page - 1) * limit

	query := `
		SELECT hash, block_number, transaction_index, from_address, to_address,
			   value, gas_limit, gas_used, gas_price, max_fee_per_gas,
			   max_priority_fee_per_gas, nonce, input_data, status,
			   contract_address, logs_bloom, created_at, updated_at
		FROM transactions 
		WHERE from_address = $1 OR to_address = $1
		ORDER BY block_number DESC, transaction_index DESC
		LIMIT $2 OFFSET $3`

	rows, err := l.db.Query(query, address, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query address transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(
			&tx.Hash, &tx.BlockNumber, &tx.TransactionIndex, &tx.FromAddress,
			&tx.ToAddress, &tx.Value, &tx.GasLimit, &tx.GasUsed, &tx.GasPrice,
			&tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas, &tx.Nonce, &tx.InputData,
			&tx.Status, &tx.ContractAddress, &tx.LogsBloom, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	// Get total count
	var totalCount int64
	err = l.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE from_address = $1 OR to_address = $1", address).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &TransactionsResponse{
		Transactions: transactions,
		TotalCount:   totalCount,
		Page:         page,
		Limit:        limit,
		Mode:         l.mode,
	}, nil
}

func (l *LiveDataService) SearchByQuery(query string) (*SearchResult, error) {
	result := &SearchResult{Mode: l.mode}

	// Try to find as block number
	if _, err := strconv.ParseInt(query, 10, 64); err == nil {
		if block, err := l.GetBlock(query); err == nil {
			result.Block = block
			return result, nil
		}
	}

	// Try to find as block hash (66 chars with 0x prefix)
	if len(query) == 66 && query[:2] == "0x" {
		if block, err := l.GetBlock(query); err == nil {
			result.Block = block
			return result, nil
		}
		if tx, err := l.GetTransaction(query); err == nil {
			result.Transaction = tx
			return result, nil
		}
	}

	// Try to find as address (42 chars with 0x prefix)
	if len(query) == 42 && query[:2] == "0x" {
		if addr, err := l.GetAddress(query); err == nil {
			result.Address = &addr.Address
			return result, nil
		}
	}

	return nil, fmt.Errorf("no results found")
}

func (l *LiveDataService) GetGasPrices(hours int) ([]GasPrice, error) {
	query := `
		SELECT id, block_number, timestamp, base_fee_per_gas, slow_gas_price,
			   standard_gas_price, fast_gas_price, slow_wait_time,
			   standard_wait_time, fast_wait_time, network_utilization, created_at
		FROM gas_prices 
		WHERE timestamp >= NOW() - INTERVAL '%d hours'
		ORDER BY timestamp DESC`

	rows, err := l.db.Query(fmt.Sprintf(query, hours))
	if err != nil {
		return nil, fmt.Errorf("failed to query gas prices: %w", err)
	}
	defer rows.Close()

	var gasPrices []GasPrice
	for rows.Next() {
		var gp GasPrice
		err := rows.Scan(
			&gp.ID, &gp.BlockNumber, &gp.Timestamp, &gp.BaseFeePerGas,
			&gp.SlowGasPrice, &gp.StandardGasPrice, &gp.FastGasPrice,
			&gp.SlowWaitTime, &gp.StandardWaitTime, &gp.FastWaitTime,
			&gp.NetworkUtilization, &gp.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gas price: %w", err)
		}
		gasPrices = append(gasPrices, gp)
	}

	return gasPrices, nil
}
