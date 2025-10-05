# Ticketer

AI-powered receipt processing application that extracts structured data from receipt images using Google Gemini.

## Features

- 📸 Upload receipt images (JPG, JPEG, PNG)
- 🤖 AI-powered text extraction with Gemini 2.5 Flash
- 🏪 Store-specific parsing (ALDI, Carrefour, and generic)
- 📊 Detailed breakdown of items, quantities, and prices
- 💰 Automatic discount calculation
- 🌙 Dark theme UI

## Tech Stack

**Backend:**
- Go 1.25
- Google Gemini AI
- Wolfi Linux containers

**Frontend:**
- Next.js 15 (App Router)
- TypeScript
- HeroUI components
- Tailwind CSS v4

## Project Structure

```
├── backend/
│   ├── cmd/server/        # Server entry point
│   ├── internal/
│   │   ├── handlers/      # HTTP handlers
│   │   ├── models/        # Data models
│   │   ├── services/ai/   # Gemini AI service
│   │   └── logger/        # Logging configuration
│   └── Dockerfile
├── frontend/
│   ├── app/               # Next.js pages
│   ├── components/        # React components
│   ├── hooks/             # Custom hooks
│   ├── types/             # TypeScript types
│   └── Dockerfile
├── docker-compose.yml
└── justfile               # Task runner commands
```

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.25+ (for local development)
- Node.js 18+ & Yarn (for local development)
- Google Gemini API key

### Setup

1. Clone the repository and create `.env` file:

```bash
cp example.env .env
```

2. Add your Gemini API key to `.env`:

```
GEMINI_API_KEY=your_api_key_here
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Running with Docker

```bash
# Start all services
just up

# View logs
just logs

# Stop services
just down
```

Visit http://localhost:3000

## Available Commands

```bash
just build-backend    # Build backend Docker image
just build-frontend   # Build frontend Docker image
just up              # Start all services
just down            # Stop all services
just logs [service]  # View logs
just restart         # Rebuild and restart
just clean           # Clean up Docker resources
```

## API Endpoints

- `POST /api/receipts/upload` - Upload and process receipt image
- `GET /api/health` - Health check

## License

MIT
