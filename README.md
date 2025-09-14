# Crypto Analytics Platform

A blockchain analytics platform with real-time data processing and deep insights into Ethereum network activity.

## 🚀 Features

### 📊 Core Blockchain Explorer
- **Block Explorer**: Navigate through blocks with detailed information including transactions, gas usage, and timestamps
- **Transaction Details**: View comprehensive transaction data with status, gas fees, and execution traces
- **Address Profiles**: Analyze address activity, balance history, and transaction patterns
- **Universal Search**: Smart search across blocks, transactions, and addresses with automatic detection and routing
- **Real-time Updates**: Live data streaming via WebSocket connections

### 🔍 Transaction Flow Analysis
- **Flow Visualization**: Interactive transaction flow diagrams
- **Address Analytics**: Comprehensive address activity analysis
- **Transaction Path Analysis**: Follow transaction paths and relationships
- **Network Analysis**: Understand transaction relationships and patterns

### 📈 Advanced Analytics
- **Network Statistics**: Real-time blockchain metrics and health indicators
- **Performance Monitoring**: Track network throughput and efficiency
- **Data Aggregation**: Historical data analysis and trend identification
- **Custom Dashboards**: Configurable analytics views

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
   
   # Configure your Ethereum RPC endpoint in backend/.env
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

## 📱 Application Pages

### Core Pages
- **Home** (`/`) - Platform overview and real-time activity feed
- **Blocks** (`/blocks`) - Block explorer with pagination
- **Block Details** (`/blocks/[id]`) - Individual block information
- **Transactions** (`/transactions`) - Transaction list and search
- **Transaction Details** (`/transactions/[hash]`) - Detailed transaction view
- **Address Profile** (`/addresses/[address]`) - Address activity and history

### Analytics Pages
- **Transaction Flow** (`/transaction-flow`) - Flow analysis and visualization

## 🔌 API Endpoints

### Core API
- `GET /api/v1/health` - Service health check
- `GET /api/v1/blocks` - List blocks
- `GET /api/v1/blocks/:identifier` - Get block details
- `GET /api/v1/transactions` - List transactions
- `GET /api/v1/transactions/:hash` - Get transaction details
- `GET /api/v1/addresses/:address` - Get address information
- `GET /api/v1/addresses/:address/transactions` - Get address transactions
- `GET /api/v1/search/:query` - Universal search endpoint
- `GET /api/v1/stats` - Network statistics

### Analytics API
- `GET /api/v1/transaction-flow/:address` - Transaction flow analysis
- `GET /api/v1/address-analytics/:address` - Address analytics
- `GET /api/v1/transaction-path` - Transaction path analysis

### WebSocket
- `ws://localhost:8080/api/v1/ws` - Real-time updates for blocks and transactions

## 🗄 Database Schema

### Core Tables
- `blocks` - Blockchain blocks with transaction data
- `transactions` - Transaction records with detailed information
- `addresses` - Address information and activity tracking

## 🔧 Configuration

### Environment Variables

#### Backend (.env)
```bash
# Database
DATABASE_URL=postgres://user:password@localhost:5432/crypto_analytics

# Ethereum
ETHEREUM_RPC=https://mainnet.infura.io/v3/your-key

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
- **Health Checks**: Service health endpoints at `/api/v1/health`
- **Database Monitoring**: Query performance tracking

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request