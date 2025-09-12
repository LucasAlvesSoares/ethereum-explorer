package main

import (
	"os"

	"crypto-analytics/backend/internal/api"
	"crypto-analytics/backend/internal/config"
	"crypto-analytics/backend/internal/database"
	"crypto-analytics/backend/internal/ethereum"
	"crypto-analytics/backend/internal/services"

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

	// Initialize data service based on mode
	var dataService services.DataService
	var ethClient *ethereum.Client

	if cfg.IsDemoMode() {
		logrus.Info("Running in DEMO mode")

		// Initialize demo data service
		dataService = services.NewDemoDataService(db, cfg.DemoDataPath)

		// Seed demo data if needed
		seeder := services.NewDemoSeeder(db, cfg.DemoDataPath)
		if err := seeder.SeedDatabase(); err != nil {
			logrus.Warnf("Failed to seed demo data: %v", err)
		}

		ethClient = nil // No Ethereum client needed in demo mode
	} else {
		logrus.Info("Running in LIVE mode")

		// Initialize Ethereum client
		var err error
		ethClient, err = ethereum.NewClient(cfg.EthereumRPC)
		if err != nil {
			logrus.Fatal("Failed to connect to Ethereum node in live mode: ", err)
		}
		defer ethClient.Close()

		// Initialize live data service
		dataService = services.NewLiveDataService(db)
	}

	// Start API server
	server := api.NewServer(db, ethClient, cfg, dataService)

	// Start ingestion service with WebSocket hub for real-time updates (only in live mode)
	if cfg.IsLiveMode() && ethClient != nil {
		ingestionService := services.NewIngestionService(db, ethClient, server.GetWebSocketHub())
		go ingestionService.Start()
		logrus.Info("Started blockchain ingestion service")
	} else {
		logrus.Info("Skipping blockchain ingestion service (demo mode)")
	}

	// Start API server (this blocks)
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
