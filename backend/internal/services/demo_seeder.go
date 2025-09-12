package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// DemoSeeder handles seeding demo data into the database
type DemoSeeder struct {
	db       *sql.DB
	dataPath string
}

// NewDemoSeeder creates a new demo seeder
func NewDemoSeeder(db *sql.DB, dataPath string) *DemoSeeder {
	return &DemoSeeder{
		db:       db,
		dataPath: dataPath,
	}
}

// SeedDatabase populates the database with demo data
func (ds *DemoSeeder) SeedDatabase() error {
	logrus.Info("Starting demo data seeding...")

	// Clear existing data first (in demo mode only)
	if err := ds.clearExistingData(); err != nil {
		logrus.Warnf("Failed to clear existing data: %v", err)
	}

	// Seed data in order of dependencies
	if err := ds.seedBlocks(); err != nil {
		logrus.Errorf("Failed to seed blocks: %v", err)
		return err
	}

	if err := ds.seedAddresses(); err != nil {
		logrus.Errorf("Failed to seed addresses: %v", err)
		return err
	}

	if err := ds.seedTransactions(); err != nil {
		logrus.Errorf("Failed to seed transactions: %v", err)
		return err
	}

	if err := ds.seedGasPrices(); err != nil {
		logrus.Errorf("Failed to seed gas prices: %v", err)
		return err
	}

	logrus.Info("Demo data seeding completed successfully")
	return nil
}

// clearExistingData removes existing data from tables
func (ds *DemoSeeder) clearExistingData() error {
	tables := []string{"events", "gas_prices", "transactions", "addresses", "blocks", "tokens"}

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		if _, err := ds.db.Exec(query); err != nil {
			logrus.Warnf("Failed to clear table %s: %v", table, err)
		}
	}

	return nil
}

// seedBlocks loads and inserts demo blocks
func (ds *DemoSeeder) seedBlocks() error {
	filePath := filepath.Join(ds.dataPath, "demo_blocks.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Warnf("Demo blocks file not found: %v", err)
		return ds.generateDefaultBlocks()
	}

	var blocks []Block
	if err := json.Unmarshal(data, &blocks); err != nil {
		logrus.Warnf("Failed to parse demo blocks: %v", err)
		return ds.generateDefaultBlocks()
	}

	query := `
		INSERT INTO blocks (number, hash, parent_hash, timestamp, gas_limit, gas_used,
			difficulty, total_difficulty, size, transaction_count, miner, extra_data, base_fee_per_gas)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (number) DO NOTHING`

	for _, block := range blocks {
		_, err := ds.db.Exec(query,
			block.Number, block.Hash, block.ParentHash, block.Timestamp,
			block.GasLimit, block.GasUsed, block.Difficulty, block.TotalDifficulty,
			block.Size, block.TransactionCount, block.Miner, block.ExtraData, block.BaseFeePerGas)
		if err != nil {
			return fmt.Errorf("failed to insert block %d: %w", block.Number, err)
		}
	}

	logrus.Infof("Seeded %d blocks", len(blocks))
	return nil
}

// seedAddresses loads and inserts demo addresses
func (ds *DemoSeeder) seedAddresses() error {
	filePath := filepath.Join(ds.dataPath, "demo_addresses.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Warnf("Demo addresses file not found: %v", err)
		return ds.generateDefaultAddresses()
	}

	var addresses []Address
	if err := json.Unmarshal(data, &addresses); err != nil {
		logrus.Warnf("Failed to parse demo addresses: %v", err)
		return ds.generateDefaultAddresses()
	}

	query := `
		INSERT INTO addresses (address, balance, nonce, is_contract, contract_creator,
			creation_transaction, first_seen_block, last_seen_block, transaction_count, label)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (address) DO NOTHING`

	for _, addr := range addresses {
		_, err := ds.db.Exec(query,
			addr.Address, addr.Balance, addr.Nonce, addr.IsContract,
			addr.ContractCreator, addr.CreationTransaction, addr.FirstSeenBlock,
			addr.LastSeenBlock, addr.TransactionCount, addr.Label)
		if err != nil {
			return fmt.Errorf("failed to insert address %s: %w", addr.Address, err)
		}
	}

	logrus.Infof("Seeded %d addresses", len(addresses))
	return nil
}

// seedTransactions loads and inserts demo transactions
func (ds *DemoSeeder) seedTransactions() error {
	filePath := filepath.Join(ds.dataPath, "demo_transactions.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Warnf("Demo transactions file not found: %v", err)
		return ds.generateDefaultTransactions()
	}

	var transactions []Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		logrus.Warnf("Failed to parse demo transactions: %v", err)
		return ds.generateDefaultTransactions()
	}

	query := `
		INSERT INTO transactions (hash, block_number, transaction_index, from_address, to_address,
			value, gas_limit, gas_used, gas_price, max_fee_per_gas, max_priority_fee_per_gas,
			nonce, input_data, status, contract_address, logs_bloom)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (hash) DO NOTHING`

	for _, tx := range transactions {
		_, err := ds.db.Exec(query,
			tx.Hash, tx.BlockNumber, tx.TransactionIndex, tx.FromAddress, tx.ToAddress,
			tx.Value, tx.GasLimit, tx.GasUsed, tx.GasPrice, tx.MaxFeePerGas,
			tx.MaxPriorityFeePerGas, tx.Nonce, tx.InputData, tx.Status,
			tx.ContractAddress, tx.LogsBloom)
		if err != nil {
			return fmt.Errorf("failed to insert transaction %s: %w", tx.Hash, err)
		}
	}

	logrus.Infof("Seeded %d transactions", len(transactions))
	return nil
}

// seedGasPrices loads and inserts demo gas prices
func (ds *DemoSeeder) seedGasPrices() error {
	filePath := filepath.Join(ds.dataPath, "demo_gas_prices.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Warnf("Demo gas prices file not found: %v", err)
		return ds.generateDefaultGasPrices()
	}

	var gasPrices []GasPrice
	if err := json.Unmarshal(data, &gasPrices); err != nil {
		logrus.Warnf("Failed to parse demo gas prices: %v", err)
		return ds.generateDefaultGasPrices()
	}

	query := `
		INSERT INTO gas_prices (block_number, timestamp, base_fee_per_gas, slow_gas_price,
			standard_gas_price, fast_gas_price, slow_wait_time, standard_wait_time,
			fast_wait_time, network_utilization)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (block_number) DO NOTHING`

	for _, gp := range gasPrices {
		_, err := ds.db.Exec(query,
			gp.BlockNumber, gp.Timestamp, gp.BaseFeePerGas, gp.SlowGasPrice,
			gp.StandardGasPrice, gp.FastGasPrice, gp.SlowWaitTime,
			gp.StandardWaitTime, gp.FastWaitTime, gp.NetworkUtilization)
		if err != nil {
			return fmt.Errorf("failed to insert gas price for block %d: %w", gp.BlockNumber, err)
		}
	}

	logrus.Infof("Seeded %d gas price records", len(gasPrices))
	return nil
}

// generateDefaultBlocks creates basic demo blocks if no file is found
func (ds *DemoSeeder) generateDefaultBlocks() error {
	logrus.Info("Generating default demo blocks...")

	baseTime := time.Now().Add(-24 * time.Hour)
	blocks := []Block{
		{
			Number:           20850000,
			Hash:             "0x123abc456def789012345678901234567890123456789012345678901234567890",
			ParentHash:       "0x098765432109876543210987654321098765432109876543210987654321098765",
			Timestamp:        baseTime,
			GasLimit:         30000000,
			GasUsed:          15000000,
			Difficulty:       "0",
			TotalDifficulty:  "58750003716598352816469",
			Size:             50000,
			TransactionCount: 150,
			Miner:            "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326",
			ExtraData:        "0x",
		},
		{
			Number:           20850001,
			Hash:             "0x234bcd567ef890123456789012345678901234567890123456789012345678901",
			ParentHash:       "0x123abc456def789012345678901234567890123456789012345678901234567890",
			Timestamp:        baseTime.Add(12 * time.Second),
			GasLimit:         30000000,
			GasUsed:          18000000,
			Difficulty:       "0",
			TotalDifficulty:  "58750003716598352816470",
			Size:             60000,
			TransactionCount: 200,
			Miner:            "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326",
			ExtraData:        "0x",
		},
	}

	query := `
		INSERT INTO blocks (number, hash, parent_hash, timestamp, gas_limit, gas_used,
			difficulty, total_difficulty, size, transaction_count, miner, extra_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (number) DO NOTHING`

	for _, block := range blocks {
		_, err := ds.db.Exec(query,
			block.Number, block.Hash, block.ParentHash, block.Timestamp,
			block.GasLimit, block.GasUsed, block.Difficulty, block.TotalDifficulty,
			block.Size, block.TransactionCount, block.Miner, block.ExtraData)
		if err != nil {
			return fmt.Errorf("failed to insert default block %d: %w", block.Number, err)
		}
	}

	return nil
}

// generateDefaultAddresses creates basic demo addresses if no file is found
func (ds *DemoSeeder) generateDefaultAddresses() error {
	logrus.Info("Generating default demo addresses...")

	addresses := []Address{
		{
			Address:          "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326",
			Balance:          "1000000000000000000000",
			Nonce:            42,
			IsContract:       false,
			TransactionCount: 150,
			Label:            stringPtr("Demo Whale Address"),
		},
		{
			Address:          "0xA0b86a33E6a0dE2a7c4D1A7F6Bc2B9e0B3e5C8a3",
			Balance:          "500000000000000000000",
			Nonce:            25,
			IsContract:       true,
			TransactionCount: 500,
			Label:            stringPtr("Demo DeFi Contract"),
		},
	}

	query := `
		INSERT INTO addresses (address, balance, nonce, is_contract, transaction_count, label)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (address) DO NOTHING`

	for _, addr := range addresses {
		_, err := ds.db.Exec(query,
			addr.Address, addr.Balance, addr.Nonce, addr.IsContract,
			addr.TransactionCount, addr.Label)
		if err != nil {
			return fmt.Errorf("failed to insert default address %s: %w", addr.Address, err)
		}
	}

	return nil
}

// generateDefaultTransactions creates basic demo transactions if no file is found
func (ds *DemoSeeder) generateDefaultTransactions() error {
	logrus.Info("Generating default demo transactions...")

	status1 := 1
	toAddr := "0xA0b86a33E6a0dE2a7c4D1A7F6Bc2B9e0B3e5C8a3"
	gasUsed := int64(21000)
	gasPrice := int64(20000000000)

	transactions := []Transaction{
		{
			Hash:             "0xabc123def456789012345678901234567890123456789012345678901234567890",
			BlockNumber:      20850000,
			TransactionIndex: 0,
			FromAddress:      "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326",
			ToAddress:        &toAddr,
			Value:            "1000000000000000000",
			GasLimit:         21000,
			GasUsed:          &gasUsed,
			GasPrice:         &gasPrice,
			Nonce:            42,
			InputData:        "0x",
			Status:           &status1,
			LogsBloom:        "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	query := `
		INSERT INTO transactions (hash, block_number, transaction_index, from_address, to_address,
			value, gas_limit, gas_used, gas_price, nonce, input_data, status, logs_bloom)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (hash) DO NOTHING`

	for _, tx := range transactions {
		_, err := ds.db.Exec(query,
			tx.Hash, tx.BlockNumber, tx.TransactionIndex, tx.FromAddress, tx.ToAddress,
			tx.Value, tx.GasLimit, tx.GasUsed, tx.GasPrice, tx.Nonce,
			tx.InputData, tx.Status, tx.LogsBloom)
		if err != nil {
			return fmt.Errorf("failed to insert default transaction %s: %w", tx.Hash, err)
		}
	}

	return nil
}

// generateDefaultGasPrices creates basic demo gas prices if no file is found
func (ds *DemoSeeder) generateDefaultGasPrices() error {
	logrus.Info("Generating default demo gas prices...")

	baseTime := time.Now().Add(-24 * time.Hour)
	baseFee := int64(15000000000)
	utilization := 75.5

	gasPrices := []GasPrice{
		{
			BlockNumber:        20850000,
			Timestamp:          baseTime,
			BaseFeePerGas:      &baseFee,
			SlowGasPrice:       18000000000,
			StandardGasPrice:   22000000000,
			FastGasPrice:       28000000000,
			SlowWaitTime:       300,
			StandardWaitTime:   180,
			FastWaitTime:       60,
			NetworkUtilization: &utilization,
		},
	}

	query := `
		INSERT INTO gas_prices (block_number, timestamp, base_fee_per_gas, slow_gas_price,
			standard_gas_price, fast_gas_price, slow_wait_time, standard_wait_time,
			fast_wait_time, network_utilization)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (block_number) DO NOTHING`

	for _, gp := range gasPrices {
		_, err := ds.db.Exec(query,
			gp.BlockNumber, gp.Timestamp, gp.BaseFeePerGas, gp.SlowGasPrice,
			gp.StandardGasPrice, gp.FastGasPrice, gp.SlowWaitTime,
			gp.StandardWaitTime, gp.FastWaitTime, gp.NetworkUtilization)
		if err != nil {
			return fmt.Errorf("failed to insert default gas price for block %d: %w", gp.BlockNumber, err)
		}
	}

	return nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
