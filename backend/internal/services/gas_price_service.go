package services

import (
	"database/sql"
	"math/big"
	"sort"
	"time"

	"crypto-analytics/backend/internal/ethereum"

	"github.com/sirupsen/logrus"
)

// GasPriceService handles gas price tracking and analysis
type GasPriceService struct {
	db        *sql.DB
	ethClient *ethereum.Client
	logger    *logrus.Logger
	stopChan  chan struct{}
}

// NewGasPriceService creates a new gas price service
func NewGasPriceService(db *sql.DB, ethClient *ethereum.Client) *GasPriceService {
	return &GasPriceService{
		db:        db,
		ethClient: ethClient,
		logger:    logrus.New(),
		stopChan:  make(chan struct{}),
	}
}

// Start begins the gas price polling service
func (s *GasPriceService) Start() {
	s.logger.Info("Starting gas price service with 2-minute polling...")

	// Do initial fetch
	if err := s.fetchAndStoreGasPrices(); err != nil {
		s.logger.Errorf("Initial gas price fetch failed: %v", err)
	}

	// Start polling every 2 minutes
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.fetchAndStoreGasPrices(); err != nil {
				s.logger.Errorf("Failed to fetch gas prices: %v", err)
			}
		case <-s.stopChan:
			s.logger.Info("Gas price service stopped")
			return
		}
	}
}

// Stop stops the gas price service
func (s *GasPriceService) Stop() {
	close(s.stopChan)
}

// fetchAndStoreGasPrices fetches current gas prices using eth_feeHistory and stores them
func (s *GasPriceService) fetchAndStoreGasPrices() error {
	s.logger.Debug("Fetching gas prices using eth_feeHistory...")

	// Get latest block number
	latestBlock, err := s.ethClient.GetLatestBlockNumber()
	if err != nil {
		return err
	}

	// Get fee history for last 20 blocks with 25th, 50th, and 75th percentiles
	rewardPercentiles := []float64{25.0, 50.0, 75.0}
	feeHistory, err := s.ethClient.FeeHistory(20, latestBlock, rewardPercentiles)
	if err != nil {
		return err
	}

	// Calculate gas prices from fee history
	gasPrices := s.calculateGasPricesFromFeeHistory(feeHistory)

	// Calculate network utilization from latest block
	utilization, err := s.calculateNetworkUtilization()
	if err != nil {
		s.logger.Warnf("Failed to calculate network utilization: %v", err)
		utilization = 50.0 // default fallback
	}

	// Store in database
	return s.storeGasPrices(gasPrices, utilization, latestBlock)
}

// calculateGasPricesFromFeeHistory calculates slow/standard/fast gas prices from fee history
func (s *GasPriceService) calculateGasPricesFromFeeHistory(feeHistory *ethereum.FeeHistory) map[string]int64 {
	// Collect all priority fees from all blocks
	var allPriorityFees []*big.Int

	for _, blockRewards := range feeHistory.Reward {
		for _, reward := range blockRewards {
			if reward != nil && reward.Sign() > 0 {
				allPriorityFees = append(allPriorityFees, reward)
			}
		}
	}

	// Get latest base fee (most recent non-nil value)
	var latestBaseFee *big.Int
	for i := len(feeHistory.BaseFeePerGas) - 1; i >= 0; i-- {
		if feeHistory.BaseFeePerGas[i] != nil {
			latestBaseFee = feeHistory.BaseFeePerGas[i]
			break
		}
	}

	// If no base fee found, use a reasonable default (15 gwei)
	if latestBaseFee == nil {
		latestBaseFee = big.NewInt(15 * 1e9) // 15 gwei in wei
	}

	// Calculate priority fee percentiles
	priorityFeePercentiles := s.calculatePercentiles(allPriorityFees, []float64{25.0, 50.0, 75.0})

	// Convert to gwei and add base fee to get total gas prices
	baseFeeGwei := latestBaseFee.Int64() / 1e9

	return map[string]int64{
		"slow":     baseFeeGwei + (priorityFeePercentiles[0] / 1e9),
		"standard": baseFeeGwei + (priorityFeePercentiles[1] / 1e9),
		"fast":     baseFeeGwei + (priorityFeePercentiles[2] / 1e9),
		"base_fee": baseFeeGwei,
	}
}

// calculatePercentiles calculates percentiles from a slice of big.Int values
func (s *GasPriceService) calculatePercentiles(values []*big.Int, percentiles []float64) []int64 {
	if len(values) == 0 {
		// Return reasonable defaults if no data (in wei)
		return []int64{2 * 1e9, 5 * 1e9, 10 * 1e9} // 2, 5, 10 gwei
	}

	// Convert to int64 and sort
	intValues := make([]int64, 0, len(values))
	for _, val := range values {
		if val.IsInt64() {
			intValues = append(intValues, val.Int64())
		}
	}

	if len(intValues) == 0 {
		return []int64{2 * 1e9, 5 * 1e9, 10 * 1e9}
	}

	sort.Slice(intValues, func(i, j int) bool {
		return intValues[i] < intValues[j]
	})

	results := make([]int64, len(percentiles))
	for i, p := range percentiles {
		index := int(float64(len(intValues)) * p / 100.0)
		if index >= len(intValues) {
			index = len(intValues) - 1
		}
		results[i] = intValues[index]
	}

	return results
}

// calculateNetworkUtilization calculates current network utilization
func (s *GasPriceService) calculateNetworkUtilization() (float64, error) {
	// Get latest block
	latestBlock, err := s.ethClient.GetLatestBlockNumber()
	if err != nil {
		return 0, err
	}

	block, err := s.ethClient.GetBlockByNumber(latestBlock)
	if err != nil {
		return 0, err
	}

	// Calculate utilization percentage
	utilization := float64(block.GasUsed()) / float64(block.GasLimit()) * 100.0
	return utilization, nil
}

// storeGasPrices stores the calculated gas prices in the database
func (s *GasPriceService) storeGasPrices(gasPrices map[string]int64, utilization float64, blockNumber *big.Int) error {
	query := `
		INSERT INTO gas_prices (
			block_number, timestamp, base_fee_per_gas, slow_gas_price, 
			standard_gas_price, fast_gas_price, network_utilization
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (block_number) DO UPDATE SET
			timestamp = EXCLUDED.timestamp,
			base_fee_per_gas = EXCLUDED.base_fee_per_gas,
			slow_gas_price = EXCLUDED.slow_gas_price,
			standard_gas_price = EXCLUDED.standard_gas_price,
			fast_gas_price = EXCLUDED.fast_gas_price,
			network_utilization = EXCLUDED.network_utilization,
			created_at = CURRENT_TIMESTAMP
	`

	baseFee := gasPrices["base_fee"] * 1e9 // Convert back to wei for storage

	_, err := s.db.Exec(query,
		blockNumber.Int64(),
		time.Now(),
		baseFee,
		gasPrices["slow"],
		gasPrices["standard"],
		gasPrices["fast"],
		utilization,
	)

	if err != nil {
		return err
	}

	s.logger.Infof("Stored gas prices - Slow: %d gwei, Standard: %d gwei, Fast: %d gwei (Utilization: %.1f%%)",
		gasPrices["slow"], gasPrices["standard"], gasPrices["fast"], utilization)

	return nil
}

// GetCurrentGasPrices returns the most recent gas prices from database
func (s *GasPriceService) GetCurrentGasPrices() (map[string]int64, error) {
	query := `
		SELECT slow_gas_price, standard_gas_price, fast_gas_price 
		FROM gas_prices 
		ORDER BY timestamp DESC 
		LIMIT 1
	`

	var slow, standard, fast int64
	err := s.db.QueryRow(query).Scan(&slow, &standard, &fast)
	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"slow":     slow,
		"standard": standard,
		"fast":     fast,
	}, nil
}
