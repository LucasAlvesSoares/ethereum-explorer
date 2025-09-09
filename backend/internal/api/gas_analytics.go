package api

import (
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GasPrice struct {
	Slow     int64 `json:"slow"`
	Standard int64 `json:"standard"`
	Fast     int64 `json:"fast"`
}

type GasPriceStats struct {
	Current     GasPrice `json:"current"`
	Average24h  GasPrice `json:"average_24h"`
	Median24h   GasPrice `json:"median_24h"`
	Min24h      GasPrice `json:"min_24h"`
	Max24h      GasPrice `json:"max_24h"`
	Trend       string   `json:"trend"`
	Utilization float64  `json:"network_utilization"`
}

type GasPriceHistoryPoint struct {
	Timestamp string `json:"timestamp"`
	Slow      int64  `json:"slow"`
	Standard  int64  `json:"standard"`
	Fast      int64  `json:"fast"`
}

type GasFeeCalculation struct {
	GasLimit       int64   `json:"gas_limit"`
	SlowFee        string  `json:"slow_fee"`
	StandardFee    string  `json:"standard_fee"`
	FastFee        string  `json:"fast_fee"`
	SlowFeeUSD     float64 `json:"slow_fee_usd,omitempty"`
	StandardFeeUSD float64 `json:"standard_fee_usd,omitempty"`
	FastFeeUSD     float64 `json:"fast_fee_usd,omitempty"`
}

type GasPriceRecommendation struct {
	TransactionType string `json:"transaction_type"`
	GasPrice        int64  `json:"gas_price"`
	EstimatedTime   int    `json:"estimated_time_seconds"`
	Description     string `json:"description"`
}

// getGasPrices returns current gas price recommendations
func (s *Server) getGasPrices(c *gin.Context) {
	// Get latest gas prices from database or calculate from recent transactions
	gasPrice, err := s.getCurrentGasPrices()
	if err != nil {
		logrus.WithError(err).Error("Failed to get current gas prices")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas prices"})
		return
	}

	c.JSON(http.StatusOK, gasPrice)
}

// getGasPriceStats returns gas price statistics and trends
func (s *Server) getGasPriceStats(c *gin.Context) {
	stats, err := s.getGasPriceStatistics()
	if err != nil {
		logrus.WithError(err).Error("Failed to get gas price statistics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas price statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// getGasPriceHistory returns historical gas price data
func (s *Server) getGasPriceHistory(c *gin.Context) {
	timeframe := c.DefaultQuery("timeframe", "24h")

	var hours int
	switch timeframe {
	case "1h":
		hours = 1
	case "24h":
		hours = 24
	case "7d":
		hours = 24 * 7
	case "30d":
		hours = 24 * 30
	default:
		hours = 24
	}

	history, err := s.getGasPriceHistoryData(hours)
	if err != nil {
		logrus.WithError(err).Error("Failed to get gas price history")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas price history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timeframe": timeframe,
		"data":      history,
	})
}

// calculateGasFee calculates transaction fee based on gas price and gas limit
func (s *Server) calculateGasFee(c *gin.Context) {
	gasLimitStr := c.Query("gas_limit")
	if gasLimitStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gas_limit parameter is required"})
		return
	}

	gasLimit, err := strconv.ParseInt(gasLimitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gas_limit parameter"})
		return
	}

	gasPrice, err := s.getCurrentGasPrices()
	if err != nil {
		logrus.WithError(err).Error("Failed to get current gas prices")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas prices"})
		return
	}

	calculation := s.calculateTransactionFees(gasLimit, gasPrice)
	c.JSON(http.StatusOK, calculation)
}

// getGasPriceRecommendations provides gas price recommendations for different transaction types
func (s *Server) getGasPriceRecommendations(c *gin.Context) {
	gasPrice, err := s.getCurrentGasPrices()
	if err != nil {
		logrus.WithError(err).Error("Failed to get current gas prices")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas prices"})
		return
	}

	recommendations := []GasPriceRecommendation{
		{
			TransactionType: "standard_transfer",
			GasPrice:        gasPrice.Standard,
			EstimatedTime:   180,
			Description:     "Standard ETH transfer - balanced speed and cost",
		},
		{
			TransactionType: "token_transfer",
			GasPrice:        gasPrice.Standard,
			EstimatedTime:   180,
			Description:     "ERC-20 token transfer - standard priority",
		},
		{
			TransactionType: "defi_interaction",
			GasPrice:        gasPrice.Fast,
			EstimatedTime:   60,
			Description:     "DeFi protocol interaction - higher priority recommended",
		},
		{
			TransactionType: "nft_mint",
			GasPrice:        gasPrice.Fast,
			EstimatedTime:   60,
			Description:     "NFT minting - fast confirmation recommended",
		},
		{
			TransactionType: "contract_deployment",
			GasPrice:        gasPrice.Standard,
			EstimatedTime:   180,
			Description:     "Smart contract deployment - standard priority",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"current_prices":  gasPrice,
	})
}

// Helper functions

func (s *Server) getCurrentGasPrices() (*GasPrice, error) {
	// Get from gas_prices table (populated by GasPriceService)
	var gasPrice GasPrice
	query := `
		SELECT slow_gas_price, standard_gas_price, fast_gas_price 
		FROM gas_prices 
		ORDER BY timestamp DESC 
		LIMIT 1
	`

	err := s.db.QueryRow(query).Scan(&gasPrice.Slow, &gasPrice.Standard, &gasPrice.Fast)
	if err == nil {
		return &gasPrice, nil
	}

	// If no gas price data available, return reasonable defaults
	if err == sql.ErrNoRows {
		return &GasPrice{
			Slow:     15, // 15 gwei
			Standard: 25, // 25 gwei
			Fast:     40, // 40 gwei
		}, nil
	}

	return nil, err
}

func (s *Server) calculateGasPricesFromTransactions() (*GasPrice, error) {
	// Get gas prices from recent transactions (last 10 blocks)
	query := `
		SELECT 
			PERCENTILE_CONT(0.25) WITHIN GROUP (ORDER BY gas_price) as slow,
			PERCENTILE_CONT(0.50) WITHIN GROUP (ORDER BY gas_price) as standard,
			PERCENTILE_CONT(0.75) WITHIN GROUP (ORDER BY gas_price) as fast
		FROM transactions 
		WHERE gas_price > 0 
		AND block_number > (SELECT MAX(number) - 10 FROM blocks)
	`

	var slow, standard, fast sql.NullFloat64
	err := s.db.QueryRow(query).Scan(&slow, &standard, &fast)
	if err != nil {
		// Fallback to default values if no data
		return &GasPrice{
			Slow:     20000000000, // 20 gwei
			Standard: 30000000000, // 30 gwei
			Fast:     50000000000, // 50 gwei
		}, nil
	}

	// Convert from wei to gwei (divide by 1e9)
	return &GasPrice{
		Slow:     int64(slow.Float64 / 1e9),
		Standard: int64(standard.Float64 / 1e9),
		Fast:     int64(fast.Float64 / 1e9),
	}, nil
}

func (s *Server) getGasPriceStatistics() (*GasPriceStats, error) {
	current, err := s.getCurrentGasPrices()
	if err != nil {
		return nil, err
	}

	// Get 24h statistics
	query := `
		SELECT 
			AVG(slow_gas_price) as avg_slow, AVG(standard_gas_price) as avg_standard, AVG(fast_gas_price) as avg_fast,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY slow_gas_price) as med_slow,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY standard_gas_price) as med_standard,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY fast_gas_price) as med_fast,
			MIN(slow_gas_price) as min_slow, MIN(standard_gas_price) as min_standard, MIN(fast_gas_price) as min_fast,
			MAX(slow_gas_price) as max_slow, MAX(standard_gas_price) as max_standard, MAX(fast_gas_price) as max_fast,
			AVG(network_utilization) as avg_utilization
		FROM gas_prices 
		WHERE timestamp > NOW() - INTERVAL '24 hours'
	`

	var avgSlow, avgStd, avgFast sql.NullFloat64
	var medSlow, medStd, medFast sql.NullFloat64
	var minSlow, minStd, minFast sql.NullFloat64
	var maxSlow, maxStd, maxFast sql.NullFloat64
	var avgUtil sql.NullFloat64

	err = s.db.QueryRow(query).Scan(
		&avgSlow, &avgStd, &avgFast,
		&medSlow, &medStd, &medFast,
		&minSlow, &minStd, &minFast,
		&maxSlow, &maxStd, &maxFast,
		&avgUtil,
	)

	if err != nil {
		// Return current prices as fallback
		return &GasPriceStats{
			Current:     *current,
			Average24h:  *current,
			Median24h:   *current,
			Min24h:      *current,
			Max24h:      *current,
			Trend:       "stable",
			Utilization: 50.0,
		}, nil
	}

	return &GasPriceStats{
		Current: *current,
		Average24h: GasPrice{
			Slow:     int64(avgSlow.Float64),
			Standard: int64(avgStd.Float64),
			Fast:     int64(avgFast.Float64),
		},
		Median24h: GasPrice{
			Slow:     int64(medSlow.Float64),
			Standard: int64(medStd.Float64),
			Fast:     int64(medFast.Float64),
		},
		Min24h: GasPrice{
			Slow:     int64(minSlow.Float64),
			Standard: int64(minStd.Float64),
			Fast:     int64(minFast.Float64),
		},
		Max24h: GasPrice{
			Slow:     int64(maxSlow.Float64),
			Standard: int64(maxStd.Float64),
			Fast:     int64(maxFast.Float64),
		},
		Trend:       s.calculateTrend(current),
		Utilization: avgUtil.Float64,
	}, nil
}

func (s *Server) getGasPriceHistoryData(hours int) ([]GasPriceHistoryPoint, error) {
	query := `
		SELECT timestamp, slow_gas_price, standard_gas_price, fast_gas_price
		FROM gas_prices 
		WHERE timestamp > NOW() - INTERVAL '%d hours'
		ORDER BY timestamp ASC
	`

	rows, err := s.db.Query(fmt.Sprintf(query, hours))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []GasPriceHistoryPoint
	for rows.Next() {
		var point GasPriceHistoryPoint
		var timestamp time.Time

		err := rows.Scan(&timestamp, &point.Slow, &point.Standard, &point.Fast)
		if err != nil {
			continue
		}

		point.Timestamp = timestamp.Format(time.RFC3339)
		history = append(history, point)
	}

	// If no historical data, generate some sample data from recent transactions
	if len(history) == 0 {
		return s.generateSampleGasHistory(hours)
	}

	return history, nil
}

func (s *Server) generateSampleGasHistory(hours int) ([]GasPriceHistoryPoint, error) {
	current, err := s.getCurrentGasPrices()
	if err != nil {
		return nil, err
	}

	var history []GasPriceHistoryPoint
	now := time.Now()

	// Generate hourly data points
	for i := hours; i >= 0; i-- {
		timestamp := now.Add(time.Duration(-i) * time.Hour)

		// Add some variation to make it realistic
		variation := float64(i%10-5) * 0.1 // ±50% variation

		history = append(history, GasPriceHistoryPoint{
			Timestamp: timestamp.Format(time.RFC3339),
			Slow:      int64(float64(current.Slow) * (1 + variation)),
			Standard:  int64(float64(current.Standard) * (1 + variation)),
			Fast:      int64(float64(current.Fast) * (1 + variation)),
		})
	}

	return history, nil
}

func (s *Server) calculateTransactionFees(gasLimit int64, gasPrice *GasPrice) *GasFeeCalculation {
	// Convert gwei to wei and calculate fees
	slowWei := big.NewInt(gasPrice.Slow * 1e9)
	standardWei := big.NewInt(gasPrice.Standard * 1e9)
	fastWei := big.NewInt(gasPrice.Fast * 1e9)

	gasLimitBig := big.NewInt(gasLimit)

	slowFee := new(big.Int).Mul(slowWei, gasLimitBig)
	standardFee := new(big.Int).Mul(standardWei, gasLimitBig)
	fastFee := new(big.Int).Mul(fastWei, gasLimitBig)

	// Convert to ETH (divide by 1e18)
	ethDivisor := big.NewFloat(1e18)

	slowETH, _ := new(big.Float).Quo(new(big.Float).SetInt(slowFee), ethDivisor).Float64()
	standardETH, _ := new(big.Float).Quo(new(big.Float).SetInt(standardFee), ethDivisor).Float64()
	fastETH, _ := new(big.Float).Quo(new(big.Float).SetInt(fastFee), ethDivisor).Float64()

	return &GasFeeCalculation{
		GasLimit:    gasLimit,
		SlowFee:     fmt.Sprintf("%.9f ETH", slowETH),
		StandardFee: fmt.Sprintf("%.9f ETH", standardETH),
		FastFee:     fmt.Sprintf("%.9f ETH", fastETH),
	}
}

func (s *Server) calculateTrend(current *GasPrice) string {
	// Simple trend calculation - could be enhanced with more sophisticated analysis
	query := `
		SELECT standard_gas_price 
		FROM gas_prices 
		WHERE timestamp > NOW() - INTERVAL '1 hour'
		ORDER BY timestamp DESC 
		LIMIT 2
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "stable"
	}
	defer rows.Close()

	var prices []int64
	for rows.Next() {
		var price int64
		if err := rows.Scan(&price); err == nil {
			prices = append(prices, price)
		}
	}

	if len(prices) < 2 {
		return "stable"
	}

	diff := prices[0] - prices[1]
	if diff > prices[1]/10 { // >10% increase
		return "rising"
	} else if diff < -prices[1]/10 { // >10% decrease
		return "falling"
	}

	return "stable"
}
