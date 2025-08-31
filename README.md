# Ethereum Blockchain Explorer

A production-ready Ethereum blockchain explorer with advanced analytics capabilities.

## Architecture

- **Backend**: Go with Gin framework
- **Frontend**: Next.js with TypeScript
- **Database**: PostgreSQL
- **Blockchain**: Ethereum via go-ethereum client

## Quick Start

```bash
# Start local development environment
docker-compose up -d

# Backend will be available at http://localhost:8080
# Frontend will be available at http://localhost:3000
```

## Project Structure

```
crypto-analytics/
├── backend/           # Go API server
├── frontend/          # Next.js application
├── database/          # PostgreSQL schemas and migrations
├── docker/           # Docker configurations
├── docs/             # Documentation
└── scripts/          # Utility scripts
```

## Environment Configuration

- **Local**: Full Docker Compose stack
- **Staging**: AWS deployment (EC2 + RDS)

## Features

### Core Explorer
- [x] Block details and navigation
- [x] Transaction details and search
- [x] Address profiles and history
- [x] Real-time updates via WebSocket
- [x] Network statistics

### Advanced Analytics
- [x] Gas price analytics
- [x] Transaction flow visualization
- [x] Address clustering and labeling
- [x] Token analytics
- [x] Smart contract analysis

## Tech Stack

### Backend Dependencies
- `gin-gonic/gin` - Web framework
- `ethereum/go-ethereum` - Ethereum client
- `lib/pq` - PostgreSQL driver
- `gorilla/websocket` - WebSocket support

### Frontend Dependencies
- `next` - React framework
- `ethers` - Ethereum JavaScript library
- `@tanstack/react-query` - Data fetching
- `recharts` - Charts and visualization
- `tailwindcss` - Styling

## Development

See individual README files in `backend/` and `frontend/` directories for detailed setup instructions.
