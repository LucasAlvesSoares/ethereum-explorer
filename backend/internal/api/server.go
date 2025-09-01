package api

import (
	"database/sql"
	"net/http"

	"crypto-analytics/backend/internal/config"
	"crypto-analytics/backend/internal/ethereum"
	"crypto-analytics/backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server represents the API server
type Server struct {
	db        *sql.DB
	ethClient *ethereum.Client
	config    *config.Config
	router    *gin.Engine
	wsHub     *websocket.Hub
}

// NewServer creates a new API server instance
func NewServer(db *sql.DB, ethClient *ethereum.Client, cfg *config.Config) *Server {
	// Set Gin mode based on environment
	if cfg.Environment == "production" || cfg.Environment == "staging" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Create and start WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	server := &Server{
		db:        db,
		ethClient: ethClient,
		config:    cfg,
		router:    router,
		wsHub:     wsHub,
	}

	server.setupRoutes()
	return server
}

// Start starts the API server
func (s *Server) Start(port string) error {
	logrus.Infof("Starting API server on port %s", port)
	return s.router.Run(":" + port)
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")

	// Health check
	api.GET("/health", s.healthCheck)

	// WebSocket endpoint
	api.GET("/ws", s.handleWebSocket)

	// Blocks
	api.GET("/blocks", s.getBlocks)
	api.GET("/blocks/:identifier", s.getBlock)

	// Transactions
	api.GET("/transactions", s.getTransactions)
	api.GET("/transactions/:hash", s.getTransaction)

	// Addresses
	api.GET("/addresses/:address", s.getAddress)
	api.GET("/addresses/:address/transactions", s.getAddressTransactions)

	// Search
	api.GET("/search/:query", s.search)

	// Network stats
	api.GET("/stats", s.getNetworkStats)

	// Gas Analytics
	api.GET("/gas/prices", s.getGasPrices)
	api.GET("/gas/stats", s.getGasPriceStats)
	api.GET("/gas/history", s.getGasPriceHistory)
	api.GET("/gas/calculate", s.calculateGasFee)
	api.GET("/gas/recommendations", s.getGasPriceRecommendations)

	// Transaction Flow Visualization
	api.GET("/transaction-flow/:address", s.GetTransactionFlow)
	api.GET("/address-analytics/:address", s.GetAddressAnalytics)
	api.GET("/transaction-path", s.GetTransactionPath)
}

// corsMiddleware handles CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(c *gin.Context) {
	s.wsHub.HandleWebSocket(c.Writer, c.Request)
}

// GetWebSocketHub returns the WebSocket hub for broadcasting messages
func (s *Server) GetWebSocketHub() *websocket.Hub {
	return s.wsHub
}

// healthCheck returns the health status of the API
func (s *Server) healthCheck(c *gin.Context) {
	response := gin.H{
		"status":      "healthy",
		"environment": s.config.Environment,
		"database":    "connected",
		"ethereum":    "disconnected",
	}

	// Check database connection
	if err := s.db.Ping(); err != nil {
		response["status"] = "degraded"
		response["database"] = "disconnected"
		logrus.Warnf("Database health check failed: %v", err)
	}

	// Check Ethereum connection (optional in demo mode)
	if s.ethClient != nil {
		if s.ethClient.IsConnected() {
			if _, err := s.ethClient.GetNetworkID(); err != nil {
				response["ethereum"] = "disconnected"
				logrus.Warnf("Ethereum health check failed: %v", err)
				// Don't mark as unhealthy in demo mode - just log the warning
			} else {
				response["ethereum"] = "connected"
			}
		} else {
			response["ethereum"] = "disconnected"
		}
	}

	// Only return unhealthy if database is down (critical)
	if response["database"] == "disconnected" {
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
