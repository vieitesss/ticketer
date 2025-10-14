# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ticketer is an AI-powered receipt processing application that extracts structured data from receipt images using Google Gemini 2.5 Flash. The system uses a multi-stage AI approach: first identifying the store, then applying store-specific prompts for accurate data extraction.

**Tech Stack:**
- Backend: Go 1.25, Fiber v3, PostgreSQL (pgx/v5), Google Gemini AI
- Frontend: Next.js 15, TypeScript, HeroUI, Tailwind CSS v4
- Infrastructure: Docker, Docker Compose, Wolfi Linux containers

## Development Commands

### Using Just (Task Runner)

```bash
# Start all services (uses docker-compose.yml)
just up

# Start production services (uses docker-compose.pro.yml)
just up pro

# Stop services
just down

# View logs (optionally filter by service)
just logs [service]

# Build images
just build-backend    # or: just bb
just build-frontend   # or: just bf
just build-all        # or: just ba

# Rebuild and restart
just restart

# Clean up Docker resources
just clean

# Run commands in containers
just run-backend [cmd]
just run-frontend [cmd]
```

### Backend Development

```bash
# Build backend manually
cd backend
go build -o main ./cmd/app

# Build Docker image
docker build -t ticketer-backend:TAG -f backend/build/Dockerfile backend/

# Run with hot reload (Air is configured)
cd backend
air  # Uses .air.toml configuration

# Run tests
cd backend
go test ./...
```

### Frontend Development

```bash
cd frontend

# Development server (with Turbopack)
yarn dev

# Build for production
yarn build

# Start production server
yarn start

# Lint
yarn lint
```

### Database

PostgreSQL is managed via Docker Compose. Schema is in `backend/internal/database/schema.sql`.

Key tables:
- `stores` - Store information
- `products` - Products linked to stores
- `receipts` - Receipt metadata with store and date
- `items` - Line items linking receipts and products with quantity/price

## Architecture

### Backend Structure

The backend follows a clean architecture pattern:

```
backend/
├── cmd/app/           - Application entry point
├── internal/
│   ├── app/           - App initialization and lifecycle
│   ├── config/        - Configuration loading (env vars)
│   ├── database/      - Database layer (PostgreSQL via pgx)
│   │   ├── schema.sql - Database schema
│   │   ├── repository.go - Repository interface
│   │   ├── receipt.go - Receipt CRUD operations
│   │   ├── product.go - Product CRUD operations
│   │   └── store.go - Store CRUD operations
│   ├── models/        - Domain models (Receipt, Product, Store, Item)
│   ├── services/      - Business logic
│   │   ├── receipt.go - Receipt processing service
│   │   └── ai/        - AI integration
│   │       └── gemini.go - Gemini API client with store-specific prompts
│   └── transport/     - HTTP layer
│       ├── http/      - Fiber server setup
│       │   ├── handlers/ - HTTP handlers
│       │   └── routers/  - Route definitions
│       └── dto/       - Data Transfer Objects (API contracts)
└── pkg/               - Shared packages (logger utilities)
```

**Key Flow:**
1. HTTP handler receives receipt upload
2. Service layer calls `GeminiService.ProcessReceipt()`
3. Gemini identifies store type from image
4. Store-specific prompt is selected (ALDI, Carrefour, or generic)
5. Gemini extracts structured data (items, quantities, prices, discounts)
6. Data is saved to PostgreSQL (stores, products, receipts, items)
7. Response DTO is built with calculated totals and returned

**Important Implementation Details:**
- Database is optional - app runs without DATABASE_URL but won't persist data
- Store identification happens first to select the appropriate extraction prompt
- Each store (ALDI, Carrefour) has custom prompts that understand receipt format
- Receipt hash prevents duplicate processing
- Products are normalized to UPPERCASE and linked to stores
- Discounts are tracked at receipt level

### Frontend Structure

```
frontend/
├── app/               - Next.js App Router pages
│   ├── page.tsx       - Home page with upload
│   └── receipts/      - Receipt detail pages
├── components/        - React components
│   ├── ReceiptTable.tsx - Receipt list table
│   ├── ReceiptDetails.tsx - Receipt detail view
│   ├── ReceiptListSidebar.tsx - Sidebar navigation
│   ├── ReceiptDetailView.tsx - Detail layout
│   └── UploadModal.tsx - Upload modal
├── hooks/             - Custom React hooks
├── lib/               - Utility functions
└── types/             - TypeScript type definitions
```

## Environment Variables

Required `.env` file:

```bash
# Backend
GEMINI_API_KEY=your_api_key_here
LOG_LEVEL=info
DATABASE_URL=postgres://user:pass@postgres:5432/dbname

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080

# Database (for docker-compose)
POSTGRES_USER=ticketer
POSTGRES_PASSWORD=ticketer
POSTGRES_DB=ticketer
```

## Docker & Deployment

### Building Images

Backend uses multi-stage build with Chainguard base images for minimal size:
- Builder: `cgr.dev/chainguard/go`
- Runtime: `cgr.dev/chainguard/static`

Build context must be `backend/` directory:
```bash
docker build -f backend/build/Dockerfile backend/
```

### Container Registry

Images are pushed to GitHub Container Registry (ghcr.io):
```bash
# Tag format
ghcr.io/vieitesss/ticketer-backend:TAG
ghcr.io/vieitesss/ticketer-frontend:TAG

# Example: tag with commit hash
git rev-parse --short=8 HEAD  # Get commit hash
docker tag IMAGE ghcr.io/vieitesss/ticketer-backend:HASH
docker push ghcr.io/vieitesss/ticketer-backend:HASH
```

### Docker Compose Configurations

- `docker-compose.yml` - Base configuration (dev mode)
- `docker-compose.override.yml` - Local development overrides
- `docker-compose.dev.yml` - Development-specific settings
- `docker-compose.pro.yml` - Production configuration

Use `just up pro` to run with production config.

## API Endpoints

```
POST /api/receipts/upload    - Upload and process receipt image
GET  /api/receipts           - List all receipts
GET  /api/receipts/:id       - Get receipt details
DELETE /api/receipts/:id     - Delete receipt
PATCH /api/receipts/items/:id - Update item quantity/price
GET  /api/health             - Health check
```

## Store-Specific AI Processing

The Gemini service uses a two-stage approach:

1. **Store Identification**: Analyzes receipt image to detect store name (ALDI, Carrefour, etc.)
2. **Data Extraction**: Applies store-specific prompt templates that understand each store's receipt format

Store prompts are in `backend/internal/services/ai/gemini.go`:
- ALDI: Handles quantity lines before product lines (`X x PRICE €`)
- Carrefour: Handles quantity lines after product lines, discount extraction
- Generic: Column-based or line-based extraction for unknown stores

When adding support for new stores, add a new case in `getStorePrompt()` with format-specific instructions.

## Current Branch

Working on `feat/update-database` branch - database schema and related functionality updates.
