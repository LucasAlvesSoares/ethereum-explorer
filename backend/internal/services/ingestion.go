package services

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"crypto-analytics/backend/internal/ethereum"
	"crypto-analytics/backend/internal/websocket"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// IngestionService handles blockchain data ingestion
type IngestionService struct {
	db        *sql.DB
	ethClient *ethereum.Client
	wsHub     *websocket.Hub
	logger    *logrus.Logger
}

// NewIngestionService creates a new ingestion service
func NewIngestionService(db *sql.DB, ethClient *ethereum.Client, wsHub *websocket.Hub) *IngestionService {
	return &IngestionService{
		db:        db,
		ethClient: ethClient,
		wsHub:     wsHub,
		logger:    logrus.New(),
	}
}

// Start begins the ingestion service with real-time block monitoring
func (s *IngestionService) Start() {
	s.logger.Info("Starting blockchain ingestion service...")

	// Try to ingest some older blocks that might work better
	if err := s.IngestOlderBlocks(5); err != nil {
		s.logger.Errorf("Failed to ingest older blocks: %v", err)
	}

	// Skip real-time ingestion for now since WebSocket subscriptions aren't supported
	s.logger.Info("Skipping real-time ingestion (WebSocket subscriptions not supported by RPC endpoint)")
}

// IngestBlock fetches and stores a block with all its transactions
func (s *IngestionService) IngestBlock(blockNumber *big.Int) error {
	s.logger.Infof("Ingesting block %s", blockNumber.String())

	// Fetch block from Ethereum
	block, err := s.ethClient.GetBlockByNumber(blockNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch block %s: %w", blockNumber.String(), err)
	}

	// Start database transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %w", err)
	}
	defer tx.Rollback()

	// Store block
	if err := s.storeBlock(tx, block); err != nil {
		return fmt.Errorf("failed to store block: %w", err)
	}

	// Store transactions
	for _, ethTx := range block.Transactions() {
		if err := s.storeTransaction(tx, ethTx, block); err != nil {
			return fmt.Errorf("failed to store transaction %s: %w", ethTx.Hash().Hex(), err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Broadcast new block to WebSocket clients
	s.broadcastNewBlock(block)

	s.logger.Infof("Successfully ingested block %s with %d transactions",
		blockNumber.String(), len(block.Transactions()))
	return nil
}

// storeBlock stores a block in the database
func (s *IngestionService) storeBlock(tx *sql.Tx, block *types.Block) error {
	query := `
		INSERT INTO blocks (
			number, hash, parent_hash, timestamp, gas_limit, gas_used,
			difficulty, total_difficulty, size, transaction_count, miner, extra_data, base_fee_per_gas
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (number) DO UPDATE SET
			hash = EXCLUDED.hash,
			parent_hash = EXCLUDED.parent_hash,
			timestamp = EXCLUDED.timestamp,
			gas_limit = EXCLUDED.gas_limit,
			gas_used = EXCLUDED.gas_used,
			difficulty = EXCLUDED.difficulty,
			total_difficulty = EXCLUDED.total_difficulty,
			size = EXCLUDED.size,
			transaction_count = EXCLUDED.transaction_count,
			miner = EXCLUDED.miner,
			extra_data = EXCLUDED.extra_data,
			base_fee_per_gas = EXCLUDED.base_fee_per_gas,
			updated_at = CURRENT_TIMESTAMP
	`

	timestamp := time.Unix(int64(block.Time()), 0)

	var baseFeePerGasInt *int64
	if block.BaseFee() != nil {
		// Convert big.Int to int64 for database storage
		if block.BaseFee().IsInt64() {
			val := block.BaseFee().Int64()
			baseFeePerGasInt = &val
		} else {
			// If the value is too large for int64, store 0 as fallback
			val := int64(0)
			baseFeePerGasInt = &val
		}
	}

	_, err := tx.Exec(query,
		block.Number().Int64(),
		block.Hash().Hex(),
		block.ParentHash().Hex(),
		timestamp,
		block.GasLimit(),
		block.GasUsed(),
		block.Difficulty().String(),
		"0", // total_difficulty - would need to calculate or fetch separately
		int(block.Size()),
		len(block.Transactions()),
		block.Coinbase().Hex(),
		fmt.Sprintf("0x%x", block.Extra()),
		baseFeePerGasInt,
	)

	return err
}

// storeTransaction stores a transaction in the database
func (s *IngestionService) storeTransaction(tx *sql.Tx, ethTx *types.Transaction, block *types.Block) error {
	// Skip transactions that might cause issues with unsupported types
	defer func() {
		if r := recover(); r != nil {
			s.logger.Warnf("Recovered from panic while processing transaction %s: %v", ethTx.Hash().Hex(), r)
		}
	}()

	// Get transaction receipt for additional data
	receipt, err := s.ethClient.GetTransactionReceipt(ethTx.Hash())
	if err != nil {
		s.logger.Warnf("Failed to get receipt for transaction %s: %v", ethTx.Hash().Hex(), err)
		// Continue without receipt data
	}

	query := `
		INSERT INTO transactions (
			hash, block_number, transaction_index, from_address, to_address, value,
			gas_limit, gas_used, gas_price, max_fee_per_gas, max_priority_fee_per_gas,
			nonce, input_data, status, contract_address, logs_bloom
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (hash) DO UPDATE SET
			block_number = EXCLUDED.block_number,
			transaction_index = EXCLUDED.transaction_index,
			from_address = EXCLUDED.from_address,
			to_address = EXCLUDED.to_address,
			value = EXCLUDED.value,
			gas_limit = EXCLUDED.gas_limit,
			gas_used = EXCLUDED.gas_used,
			gas_price = EXCLUDED.gas_price,
			max_fee_per_gas = EXCLUDED.max_fee_per_gas,
			max_priority_fee_per_gas = EXCLUDED.max_priority_fee_per_gas,
			nonce = EXCLUDED.nonce,
			input_data = EXCLUDED.input_data,
			status = EXCLUDED.status,
			contract_address = EXCLUDED.contract_address,
			logs_bloom = EXCLUDED.logs_bloom,
			updated_at = CURRENT_TIMESTAMP
	`

	// Extract sender address
	signer := types.LatestSignerForChainID(ethTx.ChainId())
	from, err := types.Sender(signer, ethTx)
	if err != nil {
		return fmt.Errorf("failed to get sender address: %w", err)
	}

	var toAddress *string
	if ethTx.To() != nil {
		addr := ethTx.To().Hex()
		toAddress = &addr
	}

	var gasUsed *uint64
	var status *int
	var contractAddress *string
	var logsBloom *string

	if receipt != nil {
		gasUsed = &receipt.GasUsed
		statusInt := int(receipt.Status)
		status = &statusInt

		if receipt.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
			addr := receipt.ContractAddress.Hex()
			contractAddress = &addr
		}

		bloom := fmt.Sprintf("%x", receipt.Bloom.Bytes())
		logsBloom = &bloom
	}

	// Find transaction index in block
	var txIndex int
	for i, blockTx := range block.Transactions() {
		if blockTx.Hash() == ethTx.Hash() {
			txIndex = i
			break
		}
	}

	// Handle gas fee values safely
	var maxFeePerGas, maxPriorityFeePerGas string
	if ethTx.GasFeeCap() != nil {
		maxFeePerGas = ethTx.GasFeeCap().String()
	} else {
		maxFeePerGas = "0"
	}
	if ethTx.GasTipCap() != nil {
		maxPriorityFeePerGas = ethTx.GasTipCap().String()
	} else {
		maxPriorityFeePerGas = "0"
	}

	_, err = tx.Exec(query,
		ethTx.Hash().Hex(),
		block.Number().Int64(),
		txIndex,
		from.Hex(),
		toAddress,
		ethTx.Value().String(),
		ethTx.Gas(),
		gasUsed,
		ethTx.GasPrice().String(),
		maxFeePerGas,
		maxPriorityFeePerGas,
		ethTx.Nonce(),
		fmt.Sprintf("%x", ethTx.Data()),
		status,
		contractAddress,
		logsBloom,
	)

	return err
}

// IngestLatestBlocks ingests the latest N blocks
func (s *IngestionService) IngestLatestBlocks(count int) error {
	latestBlockNumber, err := s.ethClient.GetLatestBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %w", err)
	}

	s.logger.Infof("Starting ingestion of latest %d blocks from block %s", count, latestBlockNumber.String())

	// Try to ingest recent blocks, but if they fail, try older blocks
	successCount := 0
	for i := 0; i < count*3 && successCount < count; i++ {
		blockNumber := new(big.Int).Sub(latestBlockNumber, big.NewInt(int64(i)))
		if blockNumber.Sign() < 0 {
			break
		}

		if err := s.IngestBlockSafely(blockNumber); err != nil {
			s.logger.Warnf("Failed to ingest block %s: %v", blockNumber.String(), err)
			continue
		}
		successCount++
	}

	s.logger.Infof("Successfully ingested %d blocks", successCount)
	return nil
}

// IngestBlockSafely tries to ingest a block with better error handling
func (s *IngestionService) IngestBlockSafely(blockNumber *big.Int) error {
	s.logger.Infof("Safely ingesting block %s", blockNumber.String())

	// Fetch block from Ethereum
	block, err := s.ethClient.GetBlockByNumber(blockNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch block %s: %w", blockNumber.String(), err)
	}

	// Start database transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %w", err)
	}
	defer tx.Rollback()

	// Store block (this is the most important part)
	if err := s.storeBlock(tx, block); err != nil {
		return fmt.Errorf("failed to store block: %w", err)
	}

	// Try to store transactions, but don't fail the whole block if transactions fail
	successfulTxs := 0
	for _, ethTx := range block.Transactions() {
		if err := s.storeTransaction(tx, ethTx, block); err != nil {
			s.logger.Warnf("Failed to store transaction %s: %v", ethTx.Hash().Hex(), err)
			continue
		}
		successfulTxs++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Broadcast new block to WebSocket clients
	s.broadcastNewBlock(block)

	s.logger.Infof("Successfully ingested block %s with %d/%d transactions",
		blockNumber.String(), successfulTxs, len(block.Transactions()))
	return nil
}

// StartRealTimeIngestion starts real-time block ingestion
func (s *IngestionService) StartRealTimeIngestion() error {
	s.logger.Info("Starting real-time block ingestion...")

	headerChan := make(chan *types.Header)
	unsubscribe, err := s.ethClient.SubscribeNewHead(headerChan)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %w", err)
	}
	defer unsubscribe()

	for header := range headerChan {
		s.logger.Infof("New block received: %s", header.Number.String())

		if err := s.IngestBlock(header.Number); err != nil {
			s.logger.Errorf("Failed to ingest new block %s: %v", header.Number.String(), err)
		} else {
			// Broadcast network stats update
			s.broadcastNetworkStats()
		}
	}

	return nil
}

// broadcastNewBlock sends new block information to WebSocket clients
func (s *IngestionService) broadcastNewBlock(block *types.Block) {
	if s.wsHub == nil {
		return
	}

	timestamp := time.Unix(int64(block.Time()), 0)

	blockUpdate := websocket.BlockUpdate{
		Number:           block.Number().Int64(),
		Hash:             block.Hash().Hex(),
		TransactionCount: len(block.Transactions()),
		GasUsed:          block.GasUsed(),
		GasLimit:         block.GasLimit(),
		Timestamp:        timestamp.Format(time.RFC3339),
		Miner:            block.Coinbase().Hex(),
	}

	s.wsHub.BroadcastBlockUpdate(blockUpdate)
}

// broadcastNetworkStats sends updated network statistics to WebSocket clients
func (s *IngestionService) broadcastNetworkStats() {
	if s.wsHub == nil {
		return
	}

	// Get latest block number for network stats
	latestBlockNumber, err := s.ethClient.GetLatestBlockNumber()
	if err != nil {
		s.logger.Errorf("Failed to get latest block number for stats: %v", err)
		return
	}

	// Get network ID
	networkID, err := s.ethClient.GetNetworkID()
	if err != nil {
		s.logger.Errorf("Failed to get network ID for stats: %v", err)
		return
	}

	statsData := map[string]interface{}{
		"latest_block": latestBlockNumber.String(),
		"network_id":   networkID.String(),
		"timestamp":    time.Now().Unix(),
	}

	s.wsHub.BroadcastNetworkStats(statsData)
}

// IngestOlderBlocks ingests older blocks that are more likely to work
func (s *IngestionService) IngestOlderBlocks(count int) error {
	latestBlockNumber, err := s.ethClient.GetLatestBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %w", err)
	}

	// Try blocks from 1000 blocks ago to avoid newer transaction types
	startBlock := new(big.Int).Sub(latestBlockNumber, big.NewInt(1000))
	s.logger.Infof("Starting ingestion of %d older blocks from block %s", count, startBlock.String())

	successCount := 0
	for i := 0; i < count*3 && successCount < count; i++ {
		blockNumber := new(big.Int).Sub(startBlock, big.NewInt(int64(i)))
		if blockNumber.Sign() < 0 {
			break
		}

		if err := s.IngestBlockSafely(blockNumber); err != nil {
			s.logger.Warnf("Failed to ingest block %s: %v", blockNumber.String(), err)
			continue
		}
		successCount++
	}

	s.logger.Infof("Successfully ingested %d older blocks", successCount)
	return nil
}
