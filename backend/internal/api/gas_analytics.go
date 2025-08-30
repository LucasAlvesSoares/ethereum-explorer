package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// GasPriceData represents gas price information
type GasPriceData struct {
	Timestamp   time.Time `json:"timestamp"`
	Slow        int64     `json:"slow"`     // Gas price for slow transactions (gwei)
	Standard    int64     `json:"standard"` // Gas price for standard transactions (gwei)
	Fast        int64     `json:"fast"`     // Gas price for fast transactions (gwei)
	Instant     int64     `json:"instant"`  // Gas price for instant transactions (gwei)
	BlockNumber int64     `json:"block_number"`
}

// GasPriceStats represents aggregated gas price statistics
type GasPriceStats struct {
	Current    GasPriceData `json:"current"`
	Average24h struct {
		Slow     float64 `json:"slow"`
		Standard float64 `json:"standard"`
		Fast     float64 `json:"fast"`
		Instant  float64 `json:"instant"`
	} `json:"average_24h"`
	Trend struct {
		Direction  string  `json:"direction"` // "up", "down", "stable"
		Percentage float64 `json:"percentage"`
	} `json:"trend"`
}

// GasPriceHistory represents historical gas price data
type GasPriceHistory struct {
	Data       []GasPriceData `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

// getGasPrices returns current gas price recommendations
func (s *Server) getGasPrices(c *gin.Context) {
	// In demo mode, return mock data
	mockData := GasPriceData{
		Timestamp:   time.Now(),
		Slow:        15,
		Standard:    25,
		Fast:        35,
		Instant:     50,
		BlockNumber: 18500000,
	}

	c.JSON(http.StatusOK, gin.H{
		"gas_prices": mockData,
		"demo_mode":  true,
	})
}

// getGasPriceStats returns gas price statistics and trends
func (s *Server) getGasPriceStats(c *gin.Context) {
	// In demo mode, return mock statistics
	stats := GasPriceStats{
		Current: GasPriceData{
			Timestamp:   time.Now(),
			Slow:        15,
			Standard:    25,
			Fast:        35,
			Instant:     50,
			BlockNumber: 18500000,
		},
	}

	stats.Average24h.Slow = 18.5
	stats.Average24h.Standard = 28.2
	stats.Average24h.Fast = 38.7
	stats.Average24h.Instant = 55.3

	stats.Trend.Direction = "down"
	stats.Trend.Percentage = -12.5

	c.JSON(http.StatusOK, gin.H{
		"stats":     stats,
		"demo_mode": true,
	})
}

// getGasPriceHistory returns historical gas price data
func (s *Server) getGasPriceHistory(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
	period := c.DefaultQuery("period", "24h") // 1h, 24h, 7d, 30d

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 24
	}

	// In demo mode, generate mock historical data
	var mockData []GasPriceData
	now := time.Now()

	// Generate data points based on period
	var interval time.Duration
	var dataPoints int

	switch period {
	case "1h":
		interval = 5 * time.Minute
		dataPoints = 12 // 12 points for 1 hour
	case "24h":
		interval = time.Hour
		dataPoints = 24 // 24 points for 24 hours
	case "7d":
		interval = 6 * time.Hour
		dataPoints = 28 // 28 points for 7 days
	case "30d":
		interval = 24 * time.Hour
		dataPoints = 30 // 30 points for 30 days
	default:
		interval = time.Hour
		dataPoints = 24
	}

	for i := 0; i < dataPoints; i++ {
		timestamp := now.Add(-time.Duration(dataPoints-i-1) * interval)

		// Generate realistic-looking gas price variations
		baseGas := 25.0
		variation := float64(i%10-5) * 2.0 // Creates some variation

		mockData = append(mockData, GasPriceData{
			Timestamp:   timestamp,
			Slow:        int64(baseGas - 10 + variation),
			Standard:    int64(baseGas + variation),
			Fast:        int64(baseGas + 10 + variation),
			Instant:     int64(baseGas + 25 + variation),
			BlockNumber: 18500000 - int64((dataPoints-i)*10),
		})
	}

	// Apply pagination
	start := (page - 1) * limit
	end := start + limit
	if start >= len(mockData) {
		mockData = []GasPriceData{}
	} else {
		if end > len(mockData) {
			end = len(mockData)
		}
		mockData = mockData[start:end]
	}

	totalPages := (dataPoints + limit - 1) / limit

	history := GasPriceHistory{
		Data: mockData,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      dataPoints,
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"history":   history,
		"period":    period,
		"demo_mode": true,
	})
}

// calculateGasFee calculates transaction fee based on gas price and gas limit
func (s *Server) calculateGasFee(c *gin.Context) {
	gasPrice, err := strconv.ParseInt(c.Query("gas_price"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gas_price parameter",
		})
		return
	}

	gasLimit, err := strconv.ParseInt(c.Query("gas_limit"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gas_limit parameter",
		})
		return
	}

	// Calculate fee in wei
	feeWei := gasPrice * gasLimit * 1e9 // Convert gwei to wei

	// Convert to ETH
	feeEth := float64(feeWei) / 1e18

	// Mock ETH price for USD calculation
	ethPriceUSD := 2500.0
	feeUSD := feeEth * ethPriceUSD

	c.JSON(http.StatusOK, gin.H{
		"gas_price": gasPrice,
		"gas_limit": gasLimit,
		"fee_wei":   feeWei,
		"fee_eth":   feeEth,
		"fee_usd":   feeUSD,
		"eth_price": ethPriceUSD,
		"demo_mode": true,
	})
}

// getGasPriceRecommendations provides gas price recommendations for different transaction types
func (s *Server) getGasPriceRecommendations(c *gin.Context) {
	transactionType := c.DefaultQuery("type", "transfer")

	// Base gas limits for different transaction types
	gasLimits := map[string]int64{
		"transfer": 21000,
		"erc20":    65000,
		"uniswap":  150000,
		"nft":      85000,
		"contract": 200000,
	}

	gasLimit, exists := gasLimits[transactionType]
	if !exists {
		gasLimit = 21000 // Default to simple transfer
	}

	// Current gas prices (mock data)
	gasPrices := GasPriceData{
		Timestamp:   time.Now(),
		Slow:        15,
		Standard:    25,
		Fast:        35,
		Instant:     50,
		BlockNumber: 18500000,
	}

	// Calculate fees for each speed
	ethPrice := 2500.0

	recommendations := gin.H{
		"transaction_type": transactionType,
		"gas_limit":        gasLimit,
		"recommendations": gin.H{
			"slow": gin.H{
				"gas_price":      gasPrices.Slow,
				"estimated_time": "5-10 minutes",
				"fee_eth":        float64(gasPrices.Slow*gasLimit) * 1e9 / 1e18,
				"fee_usd":        float64(gasPrices.Slow*gasLimit) * 1e9 / 1e18 * ethPrice,
			},
			"standard": gin.H{
				"gas_price":      gasPrices.Standard,
				"estimated_time": "2-5 minutes",
				"fee_eth":        float64(gasPrices.Standard*gasLimit) * 1e9 / 1e18,
				"fee_usd":        float64(gasPrices.Standard*gasLimit) * 1e9 / 1e18 * ethPrice,
			},
			"fast": gin.H{
				"gas_price":      gasPrices.Fast,
				"estimated_time": "1-2 minutes",
				"fee_eth":        float64(gasPrices.Fast*gasLimit) * 1e9 / 1e18,
				"fee_usd":        float64(gasPrices.Fast*gasLimit) * 1e9 / 1e18 * ethPrice,
			},
			"instant": gin.H{
				"gas_price":      gasPrices.Instant,
				"estimated_time": "< 1 minute",
				"fee_eth":        float64(gasPrices.Instant*gasLimit) * 1e9 / 1e18,
				"fee_usd":        float64(gasPrices.Instant*gasLimit) * 1e9 / 1e18 * ethPrice,
			},
		},
		"demo_mode": true,
	}

	c.JSON(http.StatusOK, recommendations)
}
