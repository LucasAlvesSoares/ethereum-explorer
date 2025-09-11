package api

import (
	"crypto-analytics/backend/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MEVAnalyticsHandler handles MEV analytics endpoints
type MEVAnalyticsHandler struct {
	mevDetector *services.MEVDetector
}

// NewMEVAnalyticsHandler creates a new MEV analytics handler
func NewMEVAnalyticsHandler(mevDetector *services.MEVDetector) *MEVAnalyticsHandler {
	return &MEVAnalyticsHandler{
		mevDetector: mevDetector,
	}
}

// RegisterMEVRoutes registers all MEV analytics routes
func (h *MEVAnalyticsHandler) RegisterMEVRoutes(r *gin.RouterGroup) {
	mev := r.Group("/mev-analytics")
	{
		mev.GET("/block/:blockNumber", h.GetBlockMEVAnalysis)
		mev.GET("/suspicious-transactions", h.GetSuspiciousTransactions)
		mev.GET("/high-gas-transactions/:blockNumber", h.GetHighGasTransactions)
		mev.GET("/sandwich-attacks/:blockNumber", h.GetSandwichAttacks)
		mev.GET("/mev-bots", h.GetMEVBots)
		mev.GET("/trends", h.GetMEVTrends)
		mev.GET("/stats", h.GetMEVStats)
	}
}

// GetBlockMEVAnalysis returns MEV analysis for a specific block
func (h *MEVAnalyticsHandler) GetBlockMEVAnalysis(c *gin.Context) {
	blockNumberStr := c.Param("blockNumber")
	blockNumber, err := strconv.ParseInt(blockNumberStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid block number",
		})
		return
	}

	analysis, err := h.mevDetector.AnalyzeBlockForMEV(blockNumber)
	if err != nil {
		logrus.WithError(err).Error("Failed to analyze block for MEV")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze block for MEV",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analysis,
	})
}

// GetSuspiciousTransactions returns suspicious transactions with potential MEV activity
func (h *MEVAnalyticsHandler) GetSuspiciousTransactions(c *gin.Context) {
	// Get query parameters
	blockNumberStr := c.Query("block_number")
	if blockNumberStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "block_number parameter is required",
		})
		return
	}

	blockNumber, err := strconv.ParseInt(blockNumberStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid block number",
		})
		return
	}

	// Get threshold multiplier (default 2.0)
	thresholdStr := c.Query("threshold")
	threshold := 2.0
	if thresholdStr != "" {
		if parsed, err := strconv.ParseFloat(thresholdStr, 64); err == nil {
			threshold = parsed
		}
	}

	// Get high gas transactions
	highGasTxs, err := h.mevDetector.DetectHighGasTransactions(blockNumber, threshold)
	if err != nil {
		logrus.WithError(err).Error("Failed to detect high gas transactions")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to detect suspicious transactions",
		})
		return
	}

	// Get sandwich attacks
	sandwichTxs, err := h.mevDetector.DetectSandwichPatterns(blockNumber)
	if err != nil {
		logrus.WithError(err).Error("Failed to detect sandwich patterns")
		// Don't return error, just log and continue with high gas transactions
	}

	// Combine results
	allSuspiciousTxs := append(highGasTxs, sandwichTxs...)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"block_number":            blockNumber,
			"threshold_multiplier":    threshold,
			"suspicious_transactions": allSuspiciousTxs,
			"total_count":             len(allSuspiciousTxs),
			"high_gas_count":          len(highGasTxs),
			"sandwich_count":          len(sandwichTxs),
		},
	})
}

// GetHighGasTransactions returns transactions with unusually high gas prices
func (h *MEVAnalyticsHandler) GetHighGasTransactions(c *gin.Context) {
	blockNumberStr := c.Param("blockNumber")
	blockNumber, err := strconv.ParseInt(blockNumberStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid block number",
		})
		return
	}

	// Get threshold multiplier (default 2.0)
	thresholdStr := c.Query("threshold")
	threshold := 2.0
	if thresholdStr != "" {
		if parsed, err := strconv.ParseFloat(thresholdStr, 64); err == nil {
			threshold = parsed
		}
	}

	transactions, err := h.mevDetector.DetectHighGasTransactions(blockNumber, threshold)
	if err != nil {
		logrus.WithError(err).Error("Failed to detect high gas transactions")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to detect high gas transactions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"block_number":         blockNumber,
			"threshold_multiplier": threshold,
			"transactions":         transactions,
			"total_count":          len(transactions),
		},
	})
}

// GetSandwichAttacks returns detected sandwich attack patterns
func (h *MEVAnalyticsHandler) GetSandwichAttacks(c *gin.Context) {
	blockNumberStr := c.Param("blockNumber")
	blockNumber, err := strconv.ParseInt(blockNumberStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid block number",
		})
		return
	}

	transactions, err := h.mevDetector.DetectSandwichPatterns(blockNumber)
	if err != nil {
		logrus.WithError(err).Error("Failed to detect sandwich patterns")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to detect sandwich attacks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"block_number": blockNumber,
			"transactions": transactions,
			"attack_count": len(transactions),
		},
	})
}

// GetMEVBots returns addresses with MEV bot-like behavior
func (h *MEVAnalyticsHandler) GetMEVBots(c *gin.Context) {
	// Parse time range parameters
	hoursStr := c.Query("hours")
	hours := 24 // Default to last 24 hours
	if hoursStr != "" {
		if parsed, err := strconv.Atoi(hoursStr); err == nil && parsed > 0 {
			hours = parsed
		}
	}

	// Parse minimum transactions parameter
	minTxStr := c.Query("min_transactions")
	minTransactions := int64(10) // Default minimum
	if minTxStr != "" {
		if parsed, err := strconv.ParseInt(minTxStr, 10, 64); err == nil && parsed > 0 {
			minTransactions = parsed
		}
	}

	// Create time range
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)
	timeRange := services.TimeRange{
		StartTime: startTime,
		EndTime:   endTime,
	}

	mevBots, err := h.mevDetector.IdentifyMEVBots(timeRange, minTransactions)
	if err != nil {
		logrus.WithError(err).Error("Failed to identify MEV bots")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to identify MEV bots",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"time_range": gin.H{
				"start_time": startTime,
				"end_time":   endTime,
				"hours":      hours,
			},
			"min_transactions": minTransactions,
			"mev_bots":         mevBots,
			"total_count":      len(mevBots),
		},
	})
}

// GetMEVTrends returns MEV trends over a specified time period
func (h *MEVAnalyticsHandler) GetMEVTrends(c *gin.Context) {
	// Parse time range parameters
	hoursStr := c.Query("hours")
	hours := 24 // Default to last 24 hours
	if hoursStr != "" {
		if parsed, err := strconv.Atoi(hoursStr); err == nil && parsed > 0 {
			hours = parsed
		}
	}

	// Create time range
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)
	timeRange := services.TimeRange{
		StartTime: startTime,
		EndTime:   endTime,
	}

	trends, err := h.mevDetector.GetMEVTrends(timeRange)
	if err != nil {
		logrus.WithError(err).Error("Failed to get MEV trends")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get MEV trends",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trends,
	})
}

// GetMEVStats returns general MEV statistics
func (h *MEVAnalyticsHandler) GetMEVStats(c *gin.Context) {
	// Parse time range parameters
	hoursStr := c.Query("hours")
	hours := 24 // Default to last 24 hours
	if hoursStr != "" {
		if parsed, err := strconv.Atoi(hoursStr); err == nil && parsed > 0 {
			hours = parsed
		}
	}

	// Create time range
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)
	timeRange := services.TimeRange{
		StartTime: startTime,
		EndTime:   endTime,
	}

	// Get trends (which includes comprehensive stats)
	stats, err := h.mevDetector.GetMEVTrends(timeRange)
	if err != nil {
		logrus.WithError(err).Error("Failed to get MEV stats")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get MEV statistics",
		})
		return
	}

	// Return simplified stats format
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"time_period": gin.H{
				"hours":      hours,
				"start_time": startTime,
				"end_time":   endTime,
			},
			"total_transactions": stats.TotalTransactions,
			"mev_transactions":   stats.HighGasTxCount,
			"mev_percentage":     stats.MEVPercentage,
			"average_gas_price":  stats.AverageGasPrice,
			"top_mev_bots_count": len(stats.TopMEVBots),
			"network_health": gin.H{
				"mev_activity_level": func() string {
					if stats.MEVPercentage > 15 {
						return "high"
					} else if stats.MEVPercentage > 5 {
						return "moderate"
					}
					return "low"
				}(),
			},
		},
	})
}
