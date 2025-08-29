package api

import (
	"database/sql"
	"net/http"

	"crypto-analytics/backend/internal/config"
	"crypto-analytics/backend/internal/ethereum"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server represents the API server
type Server struct {
	db        *sql.DB
	ethClient *ethereum.Client
	config    *config.Config
	router    *gin.Engine
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

	server := &Server{
		db:        db,
		ethClient: ethClient,
		config:    cfg,
		router:    router,
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

// healthCheck returns the health status of the API
func (s *Server) healthCheck(c *gin.Context) {
	// Check database connection
	if err := s.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	// Check Ethereum connection
	if _, err := s.ethClient.GetNetworkID(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "ethereum connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"environment": s.config.Environment,
	})
}
