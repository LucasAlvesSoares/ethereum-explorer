# Crypto Analytics Platform

A comprehensive blockchain analytics platform with real-time data processing, advanced visualization, and deep insights into Ethereum network activity.

## 🚀 Features

### 📊 Core Blockchain Explorer
- **Block Explorer**: Navigate through blocks with detailed information including transactions, gas usage, and timestamps
- **Transaction Details**: View comprehensive transaction data with status, gas fees, and execution traces
- **Address Profiles**: Analyze address activity, balance history, and transaction patterns
- **Universal Search**: Smart search across blocks, transactions, and addresses with automatic detection and routing
- **Real-time Updates**: Live data streaming via WebSocket connections

### ⛽ Gas Analytics
- **Gas Price Tracking**: Real-time and historical gas price monitoring
- **Gas Usage Analysis**: Block-level gas consumption patterns
- **Fee Calculator**: Dynamic gas fee estimation tools
- **Price Charts**: Interactive gas price visualization
- **Network Congestion**: Gas usage trends and network health indicators

### 🔍 MEV Analytics
- **MEV Detection Engine**: Advanced pattern recognition for sandwich attacks, front-running, back-running, arbitrage, and liquidations
- **Confidence Scoring**: AI-powered confidence assessment (0-100%) for MEV activity classification
- **Real-time Dashboard**: 24h MEV statistics, active bot tracking, and activity breakdowns
- **MEV Bot Profiling**: Identify and track most profitable MEV addresses with success rates
- **Interactive Filtering**: Filter activities by MEV type, confidence level, and time range
- **Historical Analysis**: Track MEV trends and network impact over time
- **Value Extraction Metrics**: Monitor total value extracted and MEV impact on network

### � Transaction Flow Analysis
- **Flow Visualization**: Interactive transaction flow diagrams
- **Address Clustering**: Group related addresses and identify patterns
- **Fund Tracking**: Follow transaction paths and money flows
- **Risk Assessment**: Identify suspicious transaction patterns
- **Network Analysis**: Understand transaction relationships

### �📈 Advanced Analytics
- **Network Statistics**: Real-time blockchain metrics and health indicators
- **Performance Monitoring**: Track network throughput and efficiency
- **Data Aggregation**: Historical data analysis and trend identification
- **Custom Dashboards**: Configurable analytics views
- **Export Capabilities**: Data export for further analysis

### 🔧 Technical Features
- **High Performance**: Optimized data ingestion and processing
- **Scalable Architecture**: Microservices-based design
- **Real-time Processing**: Live data streaming and updates
- **API Integration**: RESTful APIs for external integrations
- **Docker Support**: Containerized deployment
- **Database Optimization**: Efficient data storage and retrieval

## 🛠 Tech Stack

### Backend
- **Go** - High-performance backend services
- **Gin Framework** - Fast HTTP web framework
- **go-ethereum** - Ethereum client library
- **PostgreSQL** - Primary database
- **Redis** - Caching and session storage
- **WebSocket** - Real-time communication
- **Docker** - Containerization

### Frontend
- **Next.js 14** - React-based frontend framework
- **TypeScript** - Type-safe development
- **Tailwind CSS** - Utility-first styling
- **React Query** - Data fetching and caching
- **Recharts** - Data visualization
- **ethers.js** - Ethereum JavaScript library

### Infrastructure
- **Docker Compose** - Local development environment
- **PostgreSQL** - Relational database
- **Redis** - In-memory data store
- **Nginx** - Reverse proxy (production)

## 🏗 Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   Blockchain    │
│   (Next.js)     │◄──►│     (Go)        │◄──►│   (Ethereum)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐              │
         │              │   PostgreSQL    │              │
         └──────────────►│   Database      │◄─────────────┘
                        └─────────────────┘
                                 │
                        ┌─────────────────┐
                        │     Redis       │
                        │    Cache        │
                        └─────────────────┘
```

## 🚦 Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd crypto-analytics
   ```

2. **Environment Configuration**
   ```bash
   # Copy environment files
   cp backend/.env.example backend/.env
   
   # Configure your setup in backend/.env
   # For quick start with demo data (recommended):
   DEMO_MODE=true
   DEMO_DATA_PATH=./data
   
   # For live blockchain data:
   DEMO_MODE=false
   ETHEREUM_RPC=https://your-ethereum-node-url
   ```

3. **Start the development environment**
   ```bash
   docker-compose up -d
   ```

4. **Access the applications**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

## 🎭 Demo Mode vs Live Mode

The platform supports two operational modes:

### 🎬 Demo Mode (Default)
- **Quick Start**: No external dependencies required
- **Sample Data**: Pre-loaded with realistic blockchain data
- **Development Friendly**: Perfect for testing and development
- **Instant Setup**: Works out of the box

```bash
# Enable demo mode in backend/.env
DEMO_MODE=true
DEMO_DATA_PATH=./data
```

### 🌐 Live Mode
- **Real Data**: Connects to actual Ethereum network
- **Real-time Updates**: Live blockchain data ingestion
- **Production Ready**: Full feature set with live data
- **External Dependencies**: Requires Ethereum RPC endpoint

```bash
# Enable live mode in backend/.env
DEMO_MODE=false
ETHEREUM_RPC=https://mainnet.infura.io/v3/your-key
```

### Mode Switching
- Use the mode selector in the application header to view current mode
- Backend restart required to change modes
- Demo data is automatically seeded on first run

## 📱 Application Pages

### Core Pages
- **Home** (`/`) - Platform overview and quick stats
- **Blocks** (`/blocks`) - Block explorer with pagination
- **Block Details** (`/blocks/[id]`) - Individual block information
- **Transactions** (`/transactions`) - Transaction list and search
- **Transaction Details** (`/transactions/[hash]`) - Detailed transaction view
- **Address Profile** (`/addresses/[address]`) - Address activity and history

### Analytics Pages
- **Transaction Flow** (`/transaction-flow`) - Flow analysis and visualization
- **Gas Analytics** (`/gas-analytics`) - Gas price monitoring and analysis tools
- **MEV Analytics** (`/mev-analytics`) - MEV detection, bot tracking, and activity analysis

## 🔌 API Endpoints

### Core API
- `GET /api/v1/blocks` - List blocks
- `GET /api/v1/blocks/:id` - Get block details
- `GET /api/v1/transactions` - List transactions
- `GET /api/v1/transactions/:hash` - Get transaction details
- `GET /api/v1/addresses/:address` - Get address information
- `GET /api/v1/search/:query` - Universal search endpoint

### Analytics API
- `GET /api/v1/gas-analytics/current` - Current gas prices
- `GET /api/v1/gas-analytics/history` - Historical gas data
- `GET /api/v1/gas-analytics/calculator` - Gas fee calculations
- `GET /api/v1/mev-analytics/activities` - MEV activities with filtering
- `GET /api/v1/mev-analytics/stats` - MEV statistics and aggregated data
- `GET /api/v1/mev-analytics/activities/:id` - Individual MEV activity details
- `GET /api/v1/transaction-flow/:address` - Transaction flow analysis
- `GET /api/v1/address-analytics/:address` - Address analytics
- `GET /api/v1/transaction-path` - Transaction path analysis
- `GET /api/v1/network/stats` - Network statistics

### WebSocket
- `ws://localhost:8080/ws` - Real-time updates

## 🗄 Database Schema

### Core Tables
- `blocks` - Blockchain blocks
- `transactions` - Transaction records
- `addresses` - Address information
- `gas_prices` - Historical gas price data
- `mev_activities` - MEV detection results and analysis

### Demo Data
- `backend/data/demo_blocks.json` - Sample blockchain block data
- `backend/data/demo_transactions.json` - Sample transaction data with MEV patterns
- `backend/data/demo_addresses.json` - Sample address data including MEV bots
- `backend/data/demo_gas_prices.json` - Historical gas price data for analytics

## 🔧 Configuration

### Environment Variables

#### Backend (.env)
```bash
# Application Mode
DEMO_MODE=true
DEMO_DATA_PATH=./data

# Database
DATABASE_URL=postgres://user:password@localhost:5432/crypto_analytics

# Ethereum
ETHEREUM_RPC=https://mainnet.infura.io/v3/your-key

# Redis
REDIS_URL=redis://localhost:6379

# Server
PORT=8080
LOG_LEVEL=info
ENVIRONMENT=local
```

## 🚀 Deployment

### Production Deployment
```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Deploy with production configuration
docker-compose -f docker-compose.prod.yml up -d
```

## 🧪 Development

### Running Tests
```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && npm test
```

### Code Quality
```bash
# Go formatting and linting
cd backend && go fmt ./... && golangci-lint run

# Frontend linting
cd frontend && npm run lint
```

## 📊 Monitoring

- **Logs**: Structured JSON logging
- **Metrics**: Application performance metrics
- **Health Checks**: Service health endpoints
- **Database Monitoring**: Query performance tracking

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🔗 Links

- **Documentation**: [Coming Soon]
- **API Reference**: [Coming Soon]
- **Support**: [Issues](https://github.com/your-repo/issues)
