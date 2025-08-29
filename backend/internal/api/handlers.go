package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Block represents a blockchain block
type Block struct {
	Number           int64  `json:"number"`
	Hash             string `json:"hash"`
	ParentHash       string `json:"parent_hash"`
	Timestamp        string `json:"timestamp"`
	GasLimit         uint64 `json:"gas_limit"`
	GasUsed          uint64 `json:"gas_used"`
	Difficulty       string `json:"difficulty"`
	TotalDifficulty  string `json:"total_difficulty"`
	Size             uint64 `json:"size"`
	TransactionCount int    `json:"transaction_count"`
	Miner            string `json:"miner"`
	ExtraData        string `json:"extra_data"`
	BaseFeePerGas    string `json:"base_fee_per_gas,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash                 string  `json:"hash"`
	BlockNumber          int64   `json:"block_number"`
	TransactionIndex     int     `json:"transaction_index"`
	FromAddress          string  `json:"from_address"`
	ToAddress            *string `json:"to_address"`
	Value                string  `json:"value"`
	GasLimit             uint64  `json:"gas_limit"`
	GasUsed              *uint64 `json:"gas_used"`
	GasPrice             string  `json:"gas_price"`
	MaxFeePerGas         string  `json:"max_fee_per_gas,omitempty"`
	MaxPriorityFeePerGas string  `json:"max_priority_fee_per_gas,omitempty"`
	Nonce                uint64  `json:"nonce"`
	InputData            string  `json:"input_data"`
	Status               *int    `json:"status"`
	ContractAddress      *string `json:"contract_address"`
	LogsBloom            *string `json:"logs_bloom"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

// Address represents an Ethereum address
type Address struct {
	Address             string   `json:"address"`
	Balance             string   `json:"balance"`
	Nonce               uint64   `json:"nonce"`
	IsContract          bool     `json:"is_contract"`
	ContractCreator     *string  `json:"contract_creator"`
	CreationTransaction *string  `json:"creation_transaction"`
	FirstSeenBlock      *int64   `json:"first_seen_block"`
	LastSeenBlock       *int64   `json:"last_seen_block"`
	TransactionCount    int64    `json:"transaction_count"`
	Label               *string  `json:"label"`
	Tags                []string `json:"tags"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
}

// getBlocks returns a paginated list of blocks
func (s *Server) getBlocks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	query := `
		SELECT number, hash, parent_hash, timestamp, gas_limit, gas_used,
			   difficulty, total_difficulty, size, transaction_count, miner,
			   extra_data, base_fee_per_gas, created_at, updated_at
		FROM blocks
		ORDER BY number DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		logrus.Errorf("Failed to query blocks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocks"})
		return
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var block Block
		var baseFeePerGas sql.NullString

		err := rows.Scan(
			&block.Number, &block.Hash, &block.ParentHash, &block.Timestamp,
			&block.GasLimit, &block.GasUsed, &block.Difficulty, &block.TotalDifficulty,
			&block.Size, &block.TransactionCount, &block.Miner, &block.ExtraData,
			&baseFeePerGas, &block.CreatedAt, &block.UpdatedAt,
		)
		if err != nil {
			logrus.Errorf("Failed to scan block: %v", err)
			continue
		}

		if baseFeePerGas.Valid {
			block.BaseFeePerGas = baseFeePerGas.String
		}

		blocks = append(blocks, block)
	}

	// Get total count for pagination
	var totalCount int64
	err = s.db.QueryRow("SELECT COUNT(*) FROM blocks").Scan(&totalCount)
	if err != nil {
		logrus.Errorf("Failed to get block count: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"blocks": blocks,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       totalCount,
			"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
		},
	})
}

// getBlock returns a specific block by number or hash
func (s *Server) getBlock(c *gin.Context) {
	identifier := c.Param("identifier")

	var query string
	var args []interface{}

	// Check if identifier is a number or hash
	if strings.HasPrefix(identifier, "0x") {
		// It's a hash
		query = `
			SELECT number, hash, parent_hash, timestamp, gas_limit, gas_used,
				   difficulty, total_difficulty, size, transaction_count, miner,
				   extra_data, base_fee_per_gas, created_at, updated_at
			FROM blocks WHERE hash = $1
		`
		args = []interface{}{identifier}
	} else {
		// It's a block number
		blockNumber, err := strconv.ParseInt(identifier, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block identifier"})
			return
		}

		query = `
			SELECT number, hash, parent_hash, timestamp, gas_limit, gas_used,
				   difficulty, total_difficulty, size, transaction_count, miner,
				   extra_data, base_fee_per_gas, created_at, updated_at
			FROM blocks WHERE number = $1
		`
		args = []interface{}{blockNumber}
	}

	var block Block
	var baseFeePerGas sql.NullString

	err := s.db.QueryRow(query, args...).Scan(
		&block.Number, &block.Hash, &block.ParentHash, &block.Timestamp,
		&block.GasLimit, &block.GasUsed, &block.Difficulty, &block.TotalDifficulty,
		&block.Size, &block.TransactionCount, &block.Miner, &block.ExtraData,
		&baseFeePerGas, &block.CreatedAt, &block.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}
	if err != nil {
		logrus.Errorf("Failed to query block: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch block"})
		return
	}

	if baseFeePerGas.Valid {
		block.BaseFeePerGas = baseFeePerGas.String
	}

	c.JSON(http.StatusOK, block)
}

// getTransactions returns a paginated list of transactions
func (s *Server) getTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	blockNumber := c.Query("block")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var query string
	var args []interface{}
	var countQuery string
	var countArgs []interface{}

	if blockNumber != "" {
		// Filter by block number
		blockNum, err := strconv.ParseInt(blockNumber, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block number"})
			return
		}

		query = `
			SELECT hash, block_number, transaction_index, from_address, to_address,
				   value, gas_limit, gas_used, gas_price, max_fee_per_gas,
				   max_priority_fee_per_gas, nonce, input_data, status,
				   contract_address, logs_bloom, created_at, updated_at
			FROM transactions
			WHERE block_number = $1
			ORDER BY transaction_index ASC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{blockNum, limit, offset}
		countQuery = "SELECT COUNT(*) FROM transactions WHERE block_number = $1"
		countArgs = []interface{}{blockNum}
	} else {
		// Get all transactions
		query = `
			SELECT hash, block_number, transaction_index, from_address, to_address,
				   value, gas_limit, gas_used, gas_price, max_fee_per_gas,
				   max_priority_fee_per_gas, nonce, input_data, status,
				   contract_address, logs_bloom, created_at, updated_at
			FROM transactions
			ORDER BY block_number DESC, transaction_index ASC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
		countQuery = "SELECT COUNT(*) FROM transactions"
		countArgs = []interface{}{}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		logrus.Errorf("Failed to query transactions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var tx Transaction
		var toAddress, contractAddress, logsBloom sql.NullString
		var gasUsed sql.NullInt64
		var status sql.NullInt32

		err := rows.Scan(
			&tx.Hash, &tx.BlockNumber, &tx.TransactionIndex, &tx.FromAddress,
			&toAddress, &tx.Value, &tx.GasLimit, &gasUsed, &tx.GasPrice,
			&tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas, &tx.Nonce, &tx.InputData,
			&status, &contractAddress, &logsBloom, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			logrus.Errorf("Failed to scan transaction: %v", err)
			continue
		}

		if toAddress.Valid {
			tx.ToAddress = &toAddress.String
		}
		if gasUsed.Valid {
			gasUsedVal := uint64(gasUsed.Int64)
			tx.GasUsed = &gasUsedVal
		}
		if status.Valid {
			statusVal := int(status.Int32)
			tx.Status = &statusVal
		}
		if contractAddress.Valid {
			tx.ContractAddress = &contractAddress.String
		}
		if logsBloom.Valid {
			tx.LogsBloom = &logsBloom.String
		}

		transactions = append(transactions, tx)
	}

	// Get total count for pagination
	var totalCount int64
	err = s.db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		logrus.Errorf("Failed to get transaction count: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       totalCount,
			"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
		},
	})
}

// getTransaction returns a specific transaction by hash
func (s *Server) getTransaction(c *gin.Context) {
	hash := c.Param("hash")

	if !strings.HasPrefix(hash, "0x") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction hash"})
		return
	}

	query := `
		SELECT hash, block_number, transaction_index, from_address, to_address,
			   value, gas_limit, gas_used, gas_price, max_fee_per_gas,
			   max_priority_fee_per_gas, nonce, input_data, status,
			   contract_address, logs_bloom, created_at, updated_at
		FROM transactions WHERE hash = $1
	`

	var tx Transaction
	var toAddress, contractAddress, logsBloom sql.NullString
	var gasUsed sql.NullInt64
	var status sql.NullInt32

	err := s.db.QueryRow(query, hash).Scan(
		&tx.Hash, &tx.BlockNumber, &tx.TransactionIndex, &tx.FromAddress,
		&toAddress, &tx.Value, &tx.GasLimit, &gasUsed, &tx.GasPrice,
		&tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas, &tx.Nonce, &tx.InputData,
		&status, &contractAddress, &logsBloom, &tx.CreatedAt, &tx.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	if err != nil {
		logrus.Errorf("Failed to query transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	if toAddress.Valid {
		tx.ToAddress = &toAddress.String
	}
	if gasUsed.Valid {
		gasUsedVal := uint64(gasUsed.Int64)
		tx.GasUsed = &gasUsedVal
	}
	if status.Valid {
		statusVal := int(status.Int32)
		tx.Status = &statusVal
	}
	if contractAddress.Valid {
		tx.ContractAddress = &contractAddress.String
	}
	if logsBloom.Valid {
		tx.LogsBloom = &logsBloom.String
	}

	c.JSON(http.StatusOK, tx)
}

// getAddress returns address information
func (s *Server) getAddress(c *gin.Context) {
	address := c.Param("address")

	if !strings.HasPrefix(address, "0x") || len(address) != 42 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
		return
	}

	query := `
		SELECT address, balance, nonce, is_contract, contract_creator,
			   creation_transaction, first_seen_block, last_seen_block,
			   transaction_count, label, tags, created_at, updated_at
		FROM addresses WHERE address = $1
	`

	var addr Address
	var contractCreator, creationTransaction sql.NullString
	var firstSeenBlock, lastSeenBlock sql.NullInt64
	var label sql.NullString
	var tags sql.NullString

	err := s.db.QueryRow(query, address).Scan(
		&addr.Address, &addr.Balance, &addr.Nonce, &addr.IsContract,
		&contractCreator, &creationTransaction, &firstSeenBlock, &lastSeenBlock,
		&addr.TransactionCount, &label, &tags, &addr.CreatedAt, &addr.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Address not in database, return empty address info
		addr = Address{
			Address:          address,
			Balance:          "0",
			Nonce:            0,
			IsContract:       false,
			TransactionCount: 0,
			Tags:             []string{},
		}
	} else if err != nil {
		logrus.Errorf("Failed to query address: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch address"})
		return
	} else {
		if contractCreator.Valid {
			addr.ContractCreator = &contractCreator.String
		}
		if creationTransaction.Valid {
			addr.CreationTransaction = &creationTransaction.String
		}
		if firstSeenBlock.Valid {
			addr.FirstSeenBlock = &firstSeenBlock.Int64
		}
		if lastSeenBlock.Valid {
			addr.LastSeenBlock = &lastSeenBlock.Int64
		}
		if label.Valid {
			addr.Label = &label.String
		}

		// Parse tags from database (stored as PostgreSQL array)
		if tags.Valid && tags.String != "" {
			// Simple parsing - in production, use proper array parsing
			addr.Tags = strings.Split(strings.Trim(tags.String, "{}"), ",")
		} else {
			addr.Tags = []string{}
		}
	}

	c.JSON(http.StatusOK, addr)
}

// getAddressTransactions returns transactions for a specific address
func (s *Server) getAddressTransactions(c *gin.Context) {
	address := c.Param("address")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if !strings.HasPrefix(address, "0x") || len(address) != 42 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	query := `
		SELECT hash, block_number, transaction_index, from_address, to_address,
			   value, gas_limit, gas_used, gas_price, max_fee_per_gas,
			   max_priority_fee_per_gas, nonce, input_data, status,
			   contract_address, logs_bloom, created_at, updated_at
		FROM transactions
		WHERE from_address = $1 OR to_address = $1
		ORDER BY block_number DESC, transaction_index ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, address, limit, offset)
	if err != nil {
		logrus.Errorf("Failed to query address transactions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var tx Transaction
		var toAddress, contractAddress, logsBloom sql.NullString
		var gasUsed sql.NullInt64
		var status sql.NullInt32

		err := rows.Scan(
			&tx.Hash, &tx.BlockNumber, &tx.TransactionIndex, &tx.FromAddress,
			&toAddress, &tx.Value, &tx.GasLimit, &gasUsed, &tx.GasPrice,
			&tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas, &tx.Nonce, &tx.InputData,
			&status, &contractAddress, &logsBloom, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			logrus.Errorf("Failed to scan transaction: %v", err)
			continue
		}

		if toAddress.Valid {
			tx.ToAddress = &toAddress.String
		}
		if gasUsed.Valid {
			gasUsedVal := uint64(gasUsed.Int64)
			tx.GasUsed = &gasUsedVal
		}
		if status.Valid {
			statusVal := int(status.Int32)
			tx.Status = &statusVal
		}
		if contractAddress.Valid {
			tx.ContractAddress = &contractAddress.String
		}
		if logsBloom.Valid {
			tx.LogsBloom = &logsBloom.String
		}

		transactions = append(transactions, tx)
	}

	// Get total count for pagination
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM transactions WHERE from_address = $1 OR to_address = $1"
	err = s.db.QueryRow(countQuery, address).Scan(&totalCount)
	if err != nil {
		logrus.Errorf("Failed to get address transaction count: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       totalCount,
			"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
		},
	})
}

// search handles multi-type search (blocks, transactions, addresses)
func (s *Server) search(c *gin.Context) {
	query := c.Param("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	results := gin.H{}

	// Check if it's a block number
	if blockNumber, err := strconv.ParseInt(query, 10, 64); err == nil {
		var exists bool
		err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM blocks WHERE number = $1)", blockNumber).Scan(&exists)
		if err == nil && exists {
			results["block"] = gin.H{
				"type":   "block",
				"number": blockNumber,
				"url":    fmt.Sprintf("/api/v1/blocks/%d", blockNumber),
			}
		}
	}

	// Check if it's a hash (block or transaction)
	if strings.HasPrefix(query, "0x") && len(query) == 66 {
		// Check if it's a block hash
		var blockNumber sql.NullInt64
		err := s.db.QueryRow("SELECT number FROM blocks WHERE hash = $1", query).Scan(&blockNumber)
		if err == nil && blockNumber.Valid {
			results["block"] = gin.H{
				"type":   "block",
				"hash":   query,
				"number": blockNumber.Int64,
				"url":    fmt.Sprintf("/api/v1/blocks/%s", query),
			}
		}

		// Check if it's a transaction hash
		var txBlockNumber sql.NullInt64
		err = s.db.QueryRow("SELECT block_number FROM transactions WHERE hash = $1", query).Scan(&txBlockNumber)
		if err == nil && txBlockNumber.Valid {
			results["transaction"] = gin.H{
				"type":         "transaction",
				"hash":         query,
				"block_number": txBlockNumber.Int64,
				"url":          fmt.Sprintf("/api/v1/transactions/%s", query),
			}
		}
	}

	// Check if it's an address
	if strings.HasPrefix(query, "0x") && len(query) == 42 {
		results["address"] = gin.H{
			"type":    "address",
			"address": query,
			"url":     fmt.Sprintf("/api/v1/addresses/%s", query),
		}
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No results found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// getNetworkStats returns basic network statistics
func (s *Server) getNetworkStats(c *gin.Context) {
	stats := gin.H{}

	// Get latest block number
	var latestBlock sql.NullInt64
	err := s.db.QueryRow("SELECT MAX(number) FROM blocks").Scan(&latestBlock)
	if err == nil && latestBlock.Valid {
		stats["latest_block"] = latestBlock.Int64
	}

	// Get total blocks
	var totalBlocks int64
	err = s.db.QueryRow("SELECT COUNT(*) FROM blocks").Scan(&totalBlocks)
	if err == nil {
		stats["total_blocks"] = totalBlocks
	}

	// Get total transactions
	var totalTransactions int64
	err = s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&totalTransactions)
	if err == nil {
		stats["total_transactions"] = totalTransactions
	}

	// Get average gas price from recent blocks
	var avgGasPrice sql.NullString
	err = s.db.QueryRow(`
		SELECT AVG(CAST(gas_price AS NUMERIC))::TEXT 
		FROM transactions 
		WHERE block_number > (SELECT MAX(number) - 100 FROM blocks)
	`).Scan(&avgGasPrice)
	if err == nil && avgGasPrice.Valid {
		stats["avg_gas_price"] = avgGasPrice.String
	}

	c.JSON(http.StatusOK, stats)
}
