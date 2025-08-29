package services

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"crypto-analytics/backend/internal/ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// IngestionService handles blockchain data ingestion
type IngestionService struct {
	db        *sql.DB
	ethClient *ethereum.Client
	logger    *logrus.Logger
}

// NewIngestionService creates a new ingestion service
func NewIngestionService(db *sql.DB, ethClient *ethereum.Client) *IngestionService {
	return &IngestionService{
		db:        db,
		ethClient: ethClient,
		logger:    logrus.New(),
	}
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
	var baseFeePerGas *big.Int
	if block.BaseFee() != nil {
		baseFeePerGas = block.BaseFee()
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
		block.Size(),
		len(block.Transactions()),
		block.Coinbase().Hex(),
		fmt.Sprintf("%x", block.Extra()),
		baseFeePerGas,
	)

	return err
}

// storeTransaction stores a transaction in the database
func (s *IngestionService) storeTransaction(tx *sql.Tx, ethTx *types.Transaction, block *types.Block) error {
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
		ethTx.GasFeeCap(),
		ethTx.GasTipCap(),
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

	for i := 0; i < count; i++ {
		blockNumber := new(big.Int).Sub(latestBlockNumber, big.NewInt(int64(i)))
		if blockNumber.Sign() < 0 {
			break
		}

		if err := s.IngestBlock(blockNumber); err != nil {
			s.logger.Errorf("Failed to ingest block %s: %v", blockNumber.String(), err)
			continue
		}
	}

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
		}
	}

	return nil
}
