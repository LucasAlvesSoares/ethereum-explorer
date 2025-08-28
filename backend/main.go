package main

import (
	"os"

	"crypto-analytics/backend/internal/api"
	"crypto-analytics/backend/internal/config"
	"crypto-analytics/backend/internal/database"
	"crypto-analytics/backend/internal/ethereum"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()

	// Setup logging
	setupLogging(cfg.LogLevel)

	logrus.Info("Starting Ethereum Blockchain Explorer...")

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logrus.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.Migrate(db); err != nil {
		logrus.Fatal("Failed to run database migrations: ", err)
	}

	// Initialize Ethereum client
	ethClient, err := ethereum.NewClient(cfg.EthereumRPC)
	if err != nil {
		logrus.Fatal("Failed to connect to Ethereum node: ", err)
	}
	defer ethClient.Close()

	// Start API server
	server := api.NewServer(db, ethClient, cfg)
	if err := server.Start(cfg.Port); err != nil {
		logrus.Fatal("Failed to start server: ", err)
	}
}

func setupLogging(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
