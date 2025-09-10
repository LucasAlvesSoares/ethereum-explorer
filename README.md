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


### 🔄 Transaction Flow Analysis
- **Flow Visualization**: Interactive transaction flow diagrams
- **Address Clustering**: Group related addresses and identify patterns
- **Fund Tracking**: Follow transaction paths and money flows
- **Risk Assessment**: Identify suspicious transaction patterns
- **Network Analysis**: Understand transaction relationships

### 📈 Advanced Analytics
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
   
   # Configure your Ethereum RPC URL in backend/.env
   ETHEREUM_RPC_URL=https://your-ethereum-node-url
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
- **Home** (`/`) - Platform overview and quick stats
- **Blocks** (`/blocks`) - Block explorer with pagination
- **Block Details** (`/blocks/[id]`) - Individual block information
- **Transactions** (`/transactions`) - Transaction list and search
- **Transaction Details** (`/transactions/[hash]`) - Detailed transaction view
- **Address Profile** (`/addresses/[address]`) - Address activity and history

### Analytics Pages
- **Transaction Flow** (`/transaction-flow`) - Flow analysis and visualization
- **Gas Analytics** (`/gas-analytics`) - Gas price monitoring and analysis tools

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

## 🔧 Configuration

### Environment Variables

#### Backend (.env)
```bash
# Database
DATABASE_URL=postgres://user:password@localhost:5432/crypto_analytics

# Ethereum
ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/your-key

# Redis
REDIS_URL=redis://localhost:6379

# Server
PORT=8080
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
