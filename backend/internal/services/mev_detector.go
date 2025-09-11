package services

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
)

// MEVDetector handles MEV (Maximal Extractable Value) detection and analysis
type MEVDetector struct {
	db *sql.DB
}

// NewMEVDetector creates a new MEV detector instance
func NewMEVDetector(db *sql.DB) *MEVDetector {
	return &MEVDetector{
		db: db,
	}
}

// MEVTransaction represents a transaction with potential MEV activity
type MEVTransaction struct {
	Hash             string    `json:"hash"`
	BlockNumber      int64     `json:"block_number"`
	TransactionIndex int       `json:"transaction_index"`
	FromAddress      string    `json:"from_address"`
	ToAddress        *string   `json:"to_address"`
	Value            *big.Int  `json:"value"`
	GasPrice         *big.Int  `json:"gas_price"`
	GasUsed          *int64    `json:"gas_used"`
	GasLimit         int64     `json:"gas_limit"`
	Timestamp        time.Time `json:"timestamp"`
	MEVType          string    `json:"mev_type"`
	MEVScore         float64   `json:"mev_score"`
	PotentialProfit  *big.Int  `json:"potential_profit"`
	RelatedTxHashes  []string  `json:"related_tx_hashes"`
}

// MEVSuspiciousAddress represents an address with potential MEV bot behavior
type MEVSuspiciousAddress struct {
	Address             string   `json:"address"`
	TransactionCount    int64    `json:"transaction_count"`
	HighGasTransactions int64    `json:"high_gas_transactions"`
	AverageGasPrice     *big.Int `json:"average_gas_price"`
	MEVScore            float64  `json:"mev_score"`
	FirstSeenBlock      int64    `json:"first_seen_block"`
	LastSeenBlock       int64    `json:"last_seen_block"`
	SuspiciousPatterns  []string `json:"suspicious_patterns"`
}

// MEVAnalysis represents comprehensive MEV analysis for a block or time period
type MEVAnalysis struct {
	BlockNumber       int64                  `json:"block_number,omitempty"`
	TimeRange         *TimeRange             `json:"time_range,omitempty"`
	TotalTransactions int64                  `json:"total_transactions"`
	MEVTransactions   int64                  `json:"mev_transactions"`
	MEVPercentage     float64                `json:"mev_percentage"`
	TotalMEVValue     *big.Int               `json:"total_mev_value"`
	AverageGasPrice   *big.Int               `json:"average_gas_price"`
	HighGasTxCount    int64                  `json:"high_gas_tx_count"`
	SandwichAttacks   int64                  `json:"sandwich_attacks"`
	ArbitrageOps      int64                  `json:"arbitrage_ops"`
	FrontRunningOps   int64                  `json:"front_running_ops"`
	TopMEVBots        []MEVSuspiciousAddress `json:"top_mev_bots"`
}

type TimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// DetectHighGasTransactions identifies transactions with unusually high gas prices
func (m *MEVDetector) DetectHighGasTransactions(blockNumber int64, gasThresholdMultiplier float64) ([]MEVTransaction, error) {
	// First, get the average gas price for the block
	avgGasQuery := `
		SELECT AVG(COALESCE(gas_price, 0)) as avg_gas_price
		FROM transactions 
		WHERE block_number = $1 AND gas_price > 0
	`
	var avgGasPrice sql.NullInt64
	err := m.db.QueryRow(avgGasQuery, blockNumber).Scan(&avgGasPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to get average gas price: %w", err)
	}

	if !avgGasPrice.Valid {
		return []MEVTransaction{}, nil
	}

	// Calculate threshold (e.g., 2x average gas price)
	threshold := float64(avgGasPrice.Int64) * gasThresholdMultiplier

	// Find transactions with high gas prices
	query := `
		SELECT t.hash, t.block_number, t.transaction_index, t.from_address, 
		       t.to_address, t.value, t.gas_price, t.gas_used, t.gas_limit,
		       b.timestamp
		FROM transactions t
		JOIN blocks b ON t.block_number = b.number
		WHERE t.block_number = $1 
		  AND t.gas_price > $2
		ORDER BY t.gas_price DESC, t.transaction_index ASC
	`

	rows, err := m.db.Query(query, blockNumber, int64(threshold))
	if err != nil {
		return nil, fmt.Errorf("failed to query high gas transactions: %w", err)
	}
	defer rows.Close()

	var mevTxs []MEVTransaction
	for rows.Next() {
		var tx MEVTransaction
		var valueStr, gasPriceStr string
		var gasUsed sql.NullInt64

		err := rows.Scan(
			&tx.Hash, &tx.BlockNumber, &tx.TransactionIndex,
			&tx.FromAddress, &tx.ToAddress, &valueStr, &gasPriceStr,
			&gasUsed, &tx.GasLimit, &tx.Timestamp,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan high gas transaction")
			continue
		}

		// Convert big integers
		tx.Value = new(big.Int)
		tx.Value.SetString(valueStr, 10)
		tx.GasPrice = new(big.Int)
		tx.GasPrice.SetString(gasPriceStr, 10)

		if gasUsed.Valid {
			gasUsedVal := gasUsed.Int64
			tx.GasUsed = &gasUsedVal
		}

		// Calculate MEV score based on gas price premium
		gasPremium := float64(tx.GasPrice.Int64()) / float64(avgGasPrice.Int64)
		tx.MEVScore = gasPremium
		tx.MEVType = "high_gas_price"

		mevTxs = append(mevTxs, tx)
	}

	return mevTxs, nil
}

// DetectSandwichPatterns identifies potential sandwich attacks in a block
func (m *MEVDetector) DetectSandwichPatterns(blockNumber int64) ([]MEVTransaction, error) {
	// Look for A->B->A patterns where A transactions have high gas prices
	query := `
		WITH ranked_txs AS (
			SELECT t.hash, t.block_number, t.transaction_index, t.from_address, 
			       t.to_address, t.value, t.gas_price, t.gas_used, t.gas_limit,
			       b.timestamp,
			       LAG(t.from_address, 1) OVER (ORDER BY t.transaction_index) as prev_from,
			       LAG(t.from_address, 2) OVER (ORDER BY t.transaction_index) as prev_prev_from,
			       LAG(t.gas_price, 1) OVER (ORDER BY t.transaction_index) as prev_gas_price,
			       LAG(t.gas_price, 2) OVER (ORDER BY t.transaction_index) as prev_prev_gas_price
			FROM transactions t
			JOIN blocks b ON t.block_number = b.number
			WHERE t.block_number = $1
			ORDER BY t.transaction_index
		)
		SELECT hash, block_number, transaction_index, from_address, 
		       to_address, value, gas_price, gas_used, gas_limit, timestamp
		FROM ranked_txs
		WHERE from_address = prev_prev_from  -- Same address as 2 transactions ago
		  AND from_address != prev_from      -- Different from immediate previous
		  AND gas_price > (SELECT AVG(gas_price) * 1.5 FROM transactions WHERE block_number = $1)
		  AND prev_prev_gas_price > (SELECT AVG(gas_price) * 1.5 FROM transactions WHERE block_number = $1)
		ORDER BY transaction_index
	`

	rows, err := m.db.Query(query, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query sandwich patterns: %w", err)
	}
	defer rows.Close()

	var mevTxs []MEVTransaction
	for rows.Next() {
		var tx MEVTransaction
		var valueStr, gasPriceStr string
		var gasUsed sql.NullInt64

		err := rows.Scan(
			&tx.Hash, &tx.BlockNumber, &tx.TransactionIndex,
			&tx.FromAddress, &tx.ToAddress, &valueStr, &gasPriceStr,
			&gasUsed, &tx.GasLimit, &tx.Timestamp,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan sandwich transaction")
			continue
		}

		// Convert big integers
		tx.Value = new(big.Int)
		tx.Value.SetString(valueStr, 10)
		tx.GasPrice = new(big.Int)
		tx.GasPrice.SetString(gasPriceStr, 10)

		if gasUsed.Valid {
			gasUsedVal := gasUsed.Int64
			tx.GasUsed = &gasUsedVal
		}

		tx.MEVType = "sandwich_attack"
		tx.MEVScore = 8.0 // High score for sandwich patterns
		mevTxs = append(mevTxs, tx)
	}

	return mevTxs, nil
}

// IdentifyMEVBots identifies addresses with MEV bot-like behavior
func (m *MEVDetector) IdentifyMEVBots(timeRange TimeRange, minTransactions int64) ([]MEVSuspiciousAddress, error) {
	query := `
		WITH address_stats AS (
			SELECT 
				t.from_address as address,
				COUNT(*) as tx_count,
				COUNT(CASE WHEN t.gas_price > (
					SELECT AVG(gas_price) * 2 FROM transactions t2 
					WHERE t2.block_number = t.block_number
				) THEN 1 END) as high_gas_count,
				AVG(t.gas_price) as avg_gas_price,
				MIN(t.block_number) as first_block,
				MAX(t.block_number) as last_block,
				-- Check for rapid consecutive transactions
				COUNT(CASE WHEN EXISTS (
					SELECT 1 FROM transactions t2 
					WHERE t2.from_address = t.from_address 
					  AND t2.block_number = t.block_number 
					  AND t2.transaction_index = t.transaction_index + 1
				) THEN 1 END) as consecutive_tx_count
			FROM transactions t
			JOIN blocks b ON t.block_number = b.number
			WHERE b.timestamp BETWEEN $1 AND $2
			  AND t.gas_price > 0
			GROUP BY t.from_address
			HAVING COUNT(*) >= $3
		)
		SELECT address, tx_count, high_gas_count, avg_gas_price, 
		       first_block, last_block, consecutive_tx_count
		FROM address_stats
		WHERE (high_gas_count::float / tx_count::float) > 0.3  -- >30% high gas transactions
		   OR consecutive_tx_count > 5  -- Multiple consecutive transactions
		ORDER BY (high_gas_count::float / tx_count::float) DESC, tx_count DESC
		LIMIT 50
	`

	rows, err := m.db.Query(query, timeRange.StartTime, timeRange.EndTime, minTransactions)
	if err != nil {
		return nil, fmt.Errorf("failed to query MEV bots: %w", err)
	}
	defer rows.Close()

	var mevBots []MEVSuspiciousAddress
	for rows.Next() {
		var bot MEVSuspiciousAddress
		var avgGasPriceStr string
		var consecutiveTxCount int64

		err := rows.Scan(
			&bot.Address, &bot.TransactionCount, &bot.HighGasTransactions,
			&avgGasPriceStr, &bot.FirstSeenBlock, &bot.LastSeenBlock,
			&consecutiveTxCount,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan MEV bot data")
			continue
		}

		// Convert average gas price
		bot.AverageGasPrice = new(big.Int)
		bot.AverageGasPrice.SetString(avgGasPriceStr, 10)

		// Calculate MEV score
		highGasRatio := float64(bot.HighGasTransactions) / float64(bot.TransactionCount)
		consecutiveRatio := float64(consecutiveTxCount) / float64(bot.TransactionCount)
		bot.MEVScore = (highGasRatio * 5) + (consecutiveRatio * 3)

		// Identify suspicious patterns
		var patterns []string
		if highGasRatio > 0.5 {
			patterns = append(patterns, "frequent_high_gas")
		}
		if consecutiveTxCount > 10 {
			patterns = append(patterns, "consecutive_transactions")
		}
		if bot.TransactionCount > 1000 {
			patterns = append(patterns, "high_volume_trader")
		}

		bot.SuspiciousPatterns = patterns
		mevBots = append(mevBots, bot)
	}

	return mevBots, nil
}

// AnalyzeBlockForMEV performs comprehensive MEV analysis for a specific block
func (m *MEVDetector) AnalyzeBlockForMEV(blockNumber int64) (*MEVAnalysis, error) {
	analysis := &MEVAnalysis{
		BlockNumber: blockNumber,
	}

	// Get basic block statistics
	statsQuery := `
		SELECT COUNT(*) as total_txs, AVG(gas_price) as avg_gas_price
		FROM transactions 
		WHERE block_number = $1 AND gas_price > 0
	`
	var avgGasPriceStr sql.NullString
	err := m.db.QueryRow(statsQuery, blockNumber).Scan(&analysis.TotalTransactions, &avgGasPriceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get block statistics: %w", err)
	}

	if avgGasPriceStr.Valid {
		analysis.AverageGasPrice = new(big.Int)
		analysis.AverageGasPrice.SetString(avgGasPriceStr.String, 10)
	}

	// Detect high gas transactions
	highGasTxs, err := m.DetectHighGasTransactions(blockNumber, 2.0)
	if err != nil {
		logrus.WithError(err).Error("Failed to detect high gas transactions")
	} else {
		analysis.HighGasTxCount = int64(len(highGasTxs))
	}

	// Detect sandwich attacks
	sandwichTxs, err := m.DetectSandwichPatterns(blockNumber)
	if err != nil {
		logrus.WithError(err).Error("Failed to detect sandwich patterns")
	} else {
		analysis.SandwichAttacks = int64(len(sandwichTxs))
	}

	// Calculate MEV metrics
	analysis.MEVTransactions = analysis.HighGasTxCount + analysis.SandwichAttacks
	if analysis.TotalTransactions > 0 {
		analysis.MEVPercentage = float64(analysis.MEVTransactions) / float64(analysis.TotalTransactions) * 100
	}

	return analysis, nil
}

// GetMEVTrends returns MEV trends over a time period
func (m *MEVDetector) GetMEVTrends(timeRange TimeRange) (*MEVAnalysis, error) {
	analysis := &MEVAnalysis{
		TimeRange: &timeRange,
	}

	// Get comprehensive statistics for the time range
	query := `
		SELECT 
			COUNT(*) as total_txs,
			AVG(t.gas_price) as avg_gas_price,
			COUNT(CASE WHEN t.gas_price > (
				SELECT AVG(gas_price) * 2 FROM transactions t2 
				WHERE t2.block_number = t.block_number AND t2.gas_price > 0
			) THEN 1 END) as high_gas_count
		FROM transactions t
		JOIN blocks b ON t.block_number = b.number
		WHERE b.timestamp BETWEEN $1 AND $2 AND t.gas_price > 0
	`

	var avgGasPriceStr sql.NullString
	err := m.db.QueryRow(query, timeRange.StartTime, timeRange.EndTime).Scan(
		&analysis.TotalTransactions, &avgGasPriceStr, &analysis.HighGasTxCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get MEV trends: %w", err)
	}

	if avgGasPriceStr.Valid {
		analysis.AverageGasPrice = new(big.Int)
		analysis.AverageGasPrice.SetString(avgGasPriceStr.String, 10)
	}

	// Get top MEV bots for the period
	mevBots, err := m.IdentifyMEVBots(timeRange, 10)
	if err != nil {
		logrus.WithError(err).Error("Failed to identify MEV bots")
	} else {
		if len(mevBots) > 10 {
			analysis.TopMEVBots = mevBots[:10]
		} else {
			analysis.TopMEVBots = mevBots
		}
	}

	// Calculate percentages
	if analysis.TotalTransactions > 0 {
		analysis.MEVPercentage = float64(analysis.HighGasTxCount) / float64(analysis.TotalTransactions) * 100
	}

	return analysis, nil
}
